package v202010

import (
	"context"
	"testing"

	"github.com/HGV/alpinebits"
	v201810 "github.com/HGV/alpinebits/v201810"
	"github.com/HGV/alpinebits/v202010/guestrequests"
)

// Stops compiling when 2020-10 replaces an alias with its own type.
func TestHandlersServeEarlierVersions(t *testing.T) {
	ack := func(context.Context, alpinebits.Request[guestrequests.NotifReportRQ]) (guestrequests.NotifReportRS, error) {
		return guestrequests.NotifReportRS{}, nil
	}
	_ = alpinebits.NewRouter().
		Handle(v201810.ActionNotifReport, ack).
		Handle(ActionNotifReport, ack)
}
