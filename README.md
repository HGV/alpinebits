# AlpineBits

![github.com/HGV/alpinebits](https://github.com/HGV/alpinebits/workflows/test/badge.svg)

A robust Go library that supports all modern versions of the AlpineBits standard, providing seamless integration and interoperability.

## Installation

```sh
go get github.com/HGV/alpinebits
```

## Usage

###  Routing

```go
package main

import (
    "net/http"

	"github.com/HGV/alpinebits/v_2018_10"
    "github.com/HGV/alpinebits/v_2020_10"
    "github.com/go-chi/chi/v5"
)

func main() {
    r := chi.NewRouter()
    r.Mount("/", alpinebitsRouter())
    http.ListenAndServe(":8080", r)
}

func alpinebitsRouter() http.Handler {
    v201810, _ := v_2018_10.NewVersion()
    v202010, _ := v_2020_10.NewVersion()

    r := NewRouter()
    r.Version(v201810, func(s *Subrouter) {
		s.Action(v_2018_10.ActionHotelAvailNotif, pushHotelAvailNotif)
	})
    r.Version(v202010, func(s *Subrouter) {
		s.Action(v_2020_10.ActionHotelInvCountNotif, pushHotelInvCountNotif, alpinebits.WithCapabilities(
			v_2020_10.CapabilityHotelInvCountNotifAcceptRooms,
			v_2020_10.CapabilityHotelInvCountNotifAcceptDeltas,
			v_2020_10.CapabilityHotelInvCountNotifAcceptOutOfOrder,
			v_2020_10.CapabilityHotelInvCountNotifAcceptOutOfMarket,
			v_2020_10.CapabilityHotelInvCountNotifAcceptClosingSeasons,
		))
	})
}

func pushHotelAvailNotif(r Request) (any, error) {
	return nil, nil
}

func pushHotelInvCountNotif(r Request) (any, error) {
	return nil, nil
}
```