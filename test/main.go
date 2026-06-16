package test

import (
	"github.com/HGV/alpinebits"
	v201810 "github.com/HGV/alpinebits/v_2018_10"
	"github.com/HGV/alpinebits/v_2018_10/freerooms"
	"github.com/HGV/alpinebits/v_2018_10/rateplans"
)

func main() {
	r := alpinebits.NewRouter()

	r.Version(v201810.Version,
		alpinebits.Handle(v201810.ActionHotelAvailNotif, pushFreeRooms,
			alpinebits.WithCapabilities(
				freerooms.CapAcceptRooms,
				freerooms.CapAcceptDeltas,
			),
		),
		alpinebits.Handle(v201810.ActionHotelRatePlanNotif, pushRatePlans),
	)

	_ = r
}

func pushFreeRooms(r alpinebits.Request[freerooms.HotelAvailNotifRQ]) (freerooms.HotelAvailNotifRS, error) {
	validation := freerooms.Validate(r.Data, freerooms.ValidateOptions{
		Caps: r.Capabilities,
	})
	if !validation.Ok() {
		return freerooms.HotelAvailNotifRS{
			Version:  alpinebits.OTAVersion,
			Warnings: validation.WarningsPtr(),
			Errors:   &validation.Errors,
		}, nil
	}
	return freerooms.HotelAvailNotifRS{
		Version:  alpinebits.OTAVersion,
		Success:  &alpinebits.Success{},
		Warnings: validation.WarningsPtr(),
	}, nil
}

func pushRatePlans(r alpinebits.Request[rateplans.HotelRatePlanNotifRQ]) (rateplans.HotelRatePlanNotifRS, error) {
	return rateplans.HotelRatePlanNotifRS{}, nil
}
