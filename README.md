# AlpineBits

![github.com/HGV/alpinebits-http](./docs/alpine_bits.png)

![github.com/HGV/alpinebits](https://github.com/HGV/alpinebits/workflows/test/badge.svg)

A Go library that supports all modern versions of the AlpineBits standard, providing seamless integration and interoperability.

## Installation

```sh
go get github.com/HGV/alpinebits
```

## Usage

### Routing

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
    r := alpinebits.NewRouter()

    v201810, _ := v_2018_10.NewVersion()
    r.Version(v201810, func(s *Subrouter) {
        s.Action(v_2018_10.ActionHotelAvailNotif, pushHotelAvailNotif)
    })

    v202010, _ := v_2020_10.NewVersion()
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

### Validation

```go
validator := v_2018_10.NewHotelAvailNotifValidator(
    WithRooms(true, &map[string]map[string]struct{}{
        "DZ": {"101": {}, "102": {}},
    }),
    WithDeltas(true),
    WithBookingThreshold(true),
)
err := validator.Validate(hotelAvailNotifRQ)
```

## Testing

> [!IMPORTANT]
> Ensure `libxml2` and `libxml2-dev` are installed before running the tests. You can install them using:

```sh
sudo apt install libxml2 libxml2-dev
```

Run all tests:

```sh
go test ./...
```
