# ![github.com/HGV/alpinebits-http](./docs/alpinebits.png)

[![Go Reference](https://pkg.go.dev/badge/github.com/HGV/alpinebits.svg)](https://pkg.go.dev/github.com/HGV/alpinebits)
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
import "github.com/HGV/alpinebits/v_2018_10/freerooms"

validator := freerooms.NewHotelAvailNotifValidator(
    freerooms.WithRooms(true, &map[string]map[string]struct{}{
        "DZ": {"101": {}, "102": {}},
    }),
    freerooms.WithDeltas(true),
    freerooms.WithBookingThreshold(true),
)
err := validator.Validate(hotelAvailNotifRQ)
```

### Handshake & Client Request

```go
handshakeConfig := alpinebits.HandshakeClientConfig{
    // client supported versions, actions and capabilities
    HandshakeData: HandshakeData{
        "2020-10": map[string][]string{
            "action_OTA_Ping": nil,
            "action_OTA_HotelInvCountNotif": {
                "OTA_HotelInvCountNotif_accept_rooms",
                "OTA_HotelInvCountNotif_accept_deltas",
                "OTA_HotelInvCountNotif_accept_out_of_order",
                "OTA_HotelInvCountNotif_accept_out_of_market",
                "OTA_HotelInvCountNotif_accept_closing_seasons",
            },
        },
        "2018-10": map[string][]string{
            "action_OTA_Ping": nil,
        },
    },
}
handshakeClient, _ := alpinebits.NewHandshakeClient(handshakeConfig)
handshakeData, _, _ := handshakeClient.Ping(context.TODO())

// Use one of the versions specified in `handshakeData` as needed.
// `NegotiatedVersion()` selects the highest version supported by both
// the client and server.
switch version, actions := handshakeData.NegotiatedVersion(); version {
case "2020-10":
    client, _ := v_2020_10.NewClient(v_2020_10.ClientConfig{
        NegotiatedVersion: actions,
    })
case "2018-10":
    client, _ := v_2018_10.NewClient(v_2018_10.ClientConfig{
        NegotiatedVersion: actions,
    })
}
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
