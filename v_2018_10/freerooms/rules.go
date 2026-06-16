package freerooms

import (
	"slices"
	"strings"

	"github.com/HGV/alpinebits"
)

// Capabilities for FreeRooms action.
const (
	CapAcceptRooms            alpinebits.Capability = "action_OTA_HotelAvailNotif_accept_rooms"
	CapAcceptDeltas           alpinebits.Capability = "action_OTA_HotelAvailNotif_accept_deltas"
	CapAcceptBookingThreshold alpinebits.Capability = "action_OTA_HotelAvailNotif_accept_BookingThreshold"
)

// ValidateOptions holds parameters for validation.
type ValidateOptions struct {
	Caps       []alpinebits.Capability
	Rooms      map[string]map[string]struct{}
	Categories map[string]struct{}
}

// Validate runs all validation rules for HotelAvailNotifRQ.
func Validate(rq HotelAvailNotifRQ, opts ValidateOptions) alpinebits.RuleResult {
	return alpinebits.Validate(rq,
		alpinebits.RequiredHotelCode,
		requiredUniqueID(opts.Caps),
		alpinebits.When(notReset, requiredInvTypeCode),
		alpinebits.When(notReset, validBookingLimits(opts.Caps)),
		alpinebits.When(notReset, validInventory(opts.Caps, opts.Rooms, opts.Categories)),
		alpinebits.When(notReset, validDateRanges),
		alpinebits.When(notReset, validOverlaps(opts.Caps)),
	)
}

func notReset(rq HotelAvailNotifRQ) bool {
	return !rq.AvailStatusMessages.IsReset()
}

func requiredUniqueID(caps []alpinebits.Capability) func(HotelAvailNotifRQ) alpinebits.RuleResult {
	return func(rq HotelAvailNotifRQ) alpinebits.RuleResult {
		if rq.UniqueID != nil {
			return alpinebits.RuleResult{}
		}
		if slices.Contains(caps, CapAcceptDeltas) {
			return alpinebits.RuleResult{}
		}
		return alpinebits.RuleResult{
			Errors: []alpinebits.Error{
				alpinebits.ApplicationError(alpinebits.ErrCodeRequiredField, "UniqueID required when deltas not supported"),
			},
		}
	}
}

func requiredInvTypeCode(rq HotelAvailNotifRQ) alpinebits.RuleResult {
	var result alpinebits.RuleResult
	for i, msg := range rq.AvailStatusMessages.AvailStatusMessages {
		if strings.TrimSpace(msg.StatusApplicationControl.InvTypeCode) == "" {
			result.Errors = append(result.Errors, alpinebits.ApplicationErrorf(
				alpinebits.ErrCodeRequiredField,
				"AvailStatusMessage[%d]: missing InvTypeCode", i,
			))
		}
	}
	return result
}

func validInventory(
	caps []alpinebits.Capability,
	rooms map[string]map[string]struct{},
	categories map[string]struct{},
) func(HotelAvailNotifRQ) alpinebits.RuleResult {
	return func(rq HotelAvailNotifRQ) alpinebits.RuleResult {
		var result alpinebits.RuleResult
		supportsRooms := slices.Contains(caps, CapAcceptRooms)

		for i, msg := range rq.AvailStatusMessages.AvailStatusMessages {
			s := msg.StatusApplicationControl
			if supportsRooms {
				if strings.TrimSpace(s.InvCode) == "" {
					result.Errors = append(result.Errors, alpinebits.ApplicationErrorf(
						alpinebits.ErrCodeRequiredField,
						"AvailStatusMessage[%d]: missing InvCode", i,
					))
				} else if rooms != nil {
					if _, ok := rooms[s.InvTypeCode][s.InvCode]; !ok {
						result.Errors = append(result.Errors, alpinebits.ApplicationErrorf(
							alpinebits.ErrCodeInvalidValue,
							"AvailStatusMessage[%d]: unknown room %s/%s", i, s.InvTypeCode, s.InvCode,
						))
					}
				}
			} else if categories != nil {
				if _, ok := categories[s.InvTypeCode]; !ok {
					result.Errors = append(result.Errors, alpinebits.ApplicationErrorf(
						alpinebits.ErrCodeInvalidValue,
						"AvailStatusMessage[%d]: unknown category %s", i, s.InvTypeCode,
					))
				}
			}
		}
		return result
	}
}

func validBookingLimits(caps []alpinebits.Capability) func(HotelAvailNotifRQ) alpinebits.RuleResult {
	return func(rq HotelAvailNotifRQ) alpinebits.RuleResult {
		var result alpinebits.RuleResult
		supportsRooms := slices.Contains(caps, CapAcceptRooms)
		supportsThreshold := slices.Contains(caps, CapAcceptBookingThreshold)

		for i, msg := range rq.AvailStatusMessages.AvailStatusMessages {
			if supportsRooms && msg.BookingLimit > 1 {
				result.Errors = append(result.Errors, alpinebits.ApplicationErrorf(
					alpinebits.ErrCodeInvalidValue,
					"AvailStatusMessage[%d]: BookingLimit must be 0 or 1 when rooms supported", i,
				))
			}

			if !supportsThreshold && msg.BookingThreshold > 0 {
				result.Errors = append(result.Errors, alpinebits.ApplicationErrorf(
					alpinebits.ErrCodeInvalidValue,
					"AvailStatusMessage[%d]: BookingThreshold not supported", i,
				))
			}

			if msg.BookingThreshold > msg.BookingLimit {
				result.Errors = append(result.Errors, alpinebits.ApplicationErrorf(
					alpinebits.ErrCodeInvalidValue,
					"AvailStatusMessage[%d]: BookingThreshold exceeds BookingLimit", i,
				))
			}
		}
		return result
	}
}

func validDateRanges(rq HotelAvailNotifRQ) alpinebits.RuleResult {
	var result alpinebits.RuleResult
	for i, msg := range rq.AvailStatusMessages.AvailStatusMessages {
		if msg.StatusApplicationControl.Start.After(msg.StatusApplicationControl.End) {
			result.Errors = append(result.Errors, alpinebits.ApplicationErrorf(
				alpinebits.ErrCodeInvalidDateRange,
				"AvailStatusMessage[%d]: end before start", i,
			))
		}
	}
	return result
}

func validOverlaps(caps []alpinebits.Capability) func(HotelAvailNotifRQ) alpinebits.RuleResult {
	return func(rq HotelAvailNotifRQ) alpinebits.RuleResult {
		supportsRooms := slices.Contains(caps, CapAcceptRooms)

		// Group messages by InvCode (rooms) or InvTypeCode (categories)
		groups := make(map[string][]AvailStatusMessage)
		for _, msg := range rq.AvailStatusMessages.AvailStatusMessages {
			var key string
			if supportsRooms {
				key = msg.StatusApplicationControl.InvCode
			} else {
				key = msg.StatusApplicationControl.InvTypeCode
			}
			groups[key] = append(groups[key], msg)
		}

		// Check overlaps within each group
		var result alpinebits.RuleResult
		for key, msgs := range groups {
			if err := alpinebits.CheckOverlaps(msgs); err != nil {
				result.Errors = append(result.Errors, alpinebits.ApplicationErrorf(
					alpinebits.ErrCodeInvalidDateRange,
					"%s: %v", key, err,
				))
			}
		}
		return result
	}
}
