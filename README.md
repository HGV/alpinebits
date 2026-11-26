# alpinebits

Go server and client for [AlpineBits HotelData](https://www.alpinebits.org/), versions 2018-10 to 2026-10.

```
go get github.com/HGV/alpinebits
```

Needs Go 1.27 and cgo with libxml2 for XSD validation (`apt install libxml2-dev pkg-config`).

## Versions and actions

Each version is a package (`v201810` … `v202610`) holding its `Version` and actions. The message types and rules of each action are in a subpackage.

| Action | Variable | Subpackage |
|---|---|---|
| FreeRooms (push) | `ActionHotelAvailNotif` (2018-10), `ActionHotelInvCountNotif` (2020-10 on) | `freerooms` |
| GuestRequests (pull) | `ActionRead`, `ActionNotifReport` | `guestrequests` |
| Inventory/Basic (push) | `ActionHotelDescriptiveContentNotif` | `inventory` |
| RatePlans (push) | `ActionHotelRatePlanNotif` | `rateplans` |

## Server

```go
r := alpinebits.NewRouter(alpinebits.WithClientStore(alpinebits.NewInMemoryClientStore())).
	Handle(v202410.ActionHotelInvCountNotif, pushFreeRooms,
		alpinebits.WithCapabilities(freerooms.CapabilityAcceptRooms, freerooms.CapabilityAcceptDeltas)).
	Handle(v202410.ActionRead, pullGuestRequests).
	Handle(v202410.ActionNotifReport, ackGuestRequests)

http.Handle("/alpinebits", basicAuth(r))
```

The router answers the handshake, then checks every request against the XSD, the spec's rules and the capabilities the client negotiated. Only valid requests reach the handler:

```go
func pushFreeRooms(ctx context.Context, r alpinebits.Request[freerooms.HotelInvCountNotifRQ]) (freerooms.HotelInvCountNotifRS, error) {
	// Checks against your data are up to you.
	hotel := freerooms.HotelInventory{Categories: map[string]int{"DOUBLE": 10}}
	if errs := freerooms.CheckInventory(r.Data, hotel); len(errs) > 0 {
		return freerooms.HotelInvCountNotifRS{}, alpinebits.Reject(errs...)
	}
	// store r.Data for r.User
	return freerooms.HotelInvCountNotifRS{Response: alpinebits.OK()}, nil
}
```

- `alpinebits.Reject(errs...)` answers with an AlpineBits error response. Any other error answers 500 and is logged.
- The router checks responses only against the XSD. When you build guest requests, skip the ones that break the rules:

  ```go
  for _, res := range reservations {
  	if errs := guestrequests.ValidateHotelReservation(res); len(errs) == 0 {
  		rs.HotelReservations = append(rs.HotelReservations, res)
  	}
  }
  ```
- The router doesn't authenticate. Wrap it, and answer failures with `alpinebits.Unauthorized(w, "invalid credentials")`.
- Without a `ClientStore`, every request gets all the capabilities you list in `Handle`. With one, clients are held to their last handshake. Implement `ClientStore` to keep handshakes in your database.
- While a message's types are aliases across versions, one handler can serve several of them, such as `v202210.ActionHotelInvCountNotif` and `v202410.ActionHotelInvCountNotif`.

## Client

```go
c := alpinebits.NewClient("https://example.com/alpinebits", "user", "secret", v202410.Version)

negotiated, err := c.Handshake(ctx, alpinebits.ActionCapabilities{
	v202410.ActionHotelInvCountNotif.HandshakeName(): {freerooms.CapabilityAcceptRooms, freerooms.CapabilityAcceptDeltas},
})

// Send checks only the XSD, so check the rules first.
caps := negotiated[v202410.ActionHotelInvCountNotif.HandshakeName()]
if errs := freerooms.Validate(rq, caps); len(errs) > 0 {
	// fix rq; errs lists what is wrong
}

rs, err := c.Send(ctx, v202410.ActionHotelInvCountNotif, rq)
if se, ok := errors.AsType[*alpinebits.ServerError](err); ok {
	// No AlpineBits response: HTTP error or "ERROR:" text. Retry on status 500.
}
if rs != nil && rs.Errors != nil {
	// The server rejected the request.
}
```

## Development

```
go test ./...
```

- A version that doesn't change a message reuses the previous version's types as aliases and delegates to its rules.
- A version that changes a message has its own types and a copy of the previous rules, with the changes marked `// 20xx-10:`. `TestRulesMatch<previous>` keeps the copied checks in step.
- `testdata` holds the official AlpineBits samples that changed in that version. Tests read the others from the version where they first appeared.
