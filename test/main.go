package main

import (
	"context"
	"crypto/subtle"
	"flag"
	"log"
	"net/http"

	"github.com/HGV/alpinebits"
	v201810 "github.com/HGV/alpinebits/v201810"
	"github.com/HGV/alpinebits/v201810/freerooms"
	"github.com/HGV/alpinebits/v201810/rateplans"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var users = map[string]string{
	"test": "test",
}

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()

	ab := alpinebits.NewRouter(
		alpinebits.WithClientStore(alpinebits.NewInMemoryClientStore()),
	)
	ab.Version(v201810.Version,
		alpinebits.Handle(v201810.ActionHotelAvailNotif, pushFreeRooms,
			alpinebits.WithCapabilities(
				freerooms.CapAcceptRooms,
				freerooms.CapAcceptDeltas,
			),
		),
		alpinebits.Handle(v201810.ActionHotelRatePlanNotif, pushRatePlans,
			alpinebits.WithCapabilities(
				rateplans.CapAcceptRatePlanBookingRule,
				rateplans.CapAcceptSupplements,
			),
		),
	)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.With(basicAuth).Handle("/alpinebits", ab)

	log.Printf("listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, r))
}

// basicAuth answers failures with Unauthorized, as AlpineBits §2 requires.
func basicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		user, pass, ok := req.BasicAuth()
		want, known := users[user]
		if !ok || !known || subtle.ConstantTimeCompare([]byte(pass), []byte(want)) != 1 {
			alpinebits.Unauthorized(w, "invalid credentials")
			return
		}
		next.ServeHTTP(w, req)
	})
}

// The router has already checked the XSD and the spec's rules.
func pushFreeRooms(_ context.Context, r alpinebits.Request[freerooms.HotelAvailNotifRQ]) (freerooms.HotelAvailNotifRS, error) {
	hotel := freerooms.HotelInventory{
		Rooms: map[string]map[string]struct{}{"DZ": {"101": {}, "102": {}}},
	}
	if errs := freerooms.CheckInventory(r.Data, hotel); len(errs) > 0 {
		return freerooms.HotelAvailNotifRS{}, alpinebits.Reject(errs...)
	}
	// ... store availability ...
	return freerooms.HotelAvailNotifRS{Response: alpinebits.OK()}, nil
}

func pushRatePlans(_ context.Context, r alpinebits.Request[rateplans.HotelRatePlanNotifRQ]) (rateplans.HotelRatePlanNotifRS, error) {
	categories := map[string]struct{}{"DZ": {}}
	if errs := rateplans.CheckInventory(r.Data, categories); len(errs) > 0 {
		return rateplans.HotelRatePlanNotifRS{}, alpinebits.Reject(errs...)
	}
	// ... store rate plans ...
	return rateplans.HotelRatePlanNotifRS{Response: alpinebits.OK()}, nil
}
