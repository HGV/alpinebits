package guestrequests

import (
	"fmt"
	"slices"
	"strings"

	"github.com/HGV/alpinebits"
)

// ValidMealPlanCodes are the allowed values for MealsIncluded.MealPlanCodes.
var ValidMealPlanCodes = []string{
	"1",  // All inclusive
	"3",  // Bed and breakfast
	"10", // Full board
	"12", // Half board
	"14", // Room only
}

// resCtx provides context for HotelReservation validation errors.
type resCtx struct {
	index int
}

func (c resCtx) err(code alpinebits.ErrCode, msg string) alpinebits.Error {
	return alpinebits.ApplicationErrorf(code, "HotelReservation[%d]: %s", c.index, msg)
}

func (c resCtx) errf(code alpinebits.ErrCode, format string, args ...any) alpinebits.Error {
	return alpinebits.ApplicationErrorf(code, "HotelReservation[%d]: "+format, append([]any{c.index}, args...)...)
}

// stayCtx provides context for RoomStay validation errors.
type stayCtx struct {
	resIndex  int
	stayIndex int
}

func (c stayCtx) err(code alpinebits.ErrCode, msg string) alpinebits.Error {
	return alpinebits.ApplicationErrorf(code, "HotelReservation[%d].RoomStay[%d]: %s", c.resIndex, c.stayIndex, msg)
}

func (c stayCtx) errf(code alpinebits.ErrCode, format string, args ...any) alpinebits.Error {
	msg := fmt.Sprintf(format, args...)
	return alpinebits.ApplicationErrorf(code, "HotelReservation[%d].RoomStay[%d]: %s", c.resIndex, c.stayIndex, msg)
}

func (c stayCtx) result(code alpinebits.ErrCode, msg string) alpinebits.RuleResult {
	return alpinebits.RuleResult{Errors: []alpinebits.Error{c.err(code, msg)}}
}

func (c stayCtx) resultf(code alpinebits.ErrCode, format string, args ...any) alpinebits.RuleResult {
	return alpinebits.RuleResult{Errors: []alpinebits.Error{c.errf(code, format, args...)}}
}

// ValidateReadRQ runs validation rules for OTA_ReadRQ.
func ValidateReadRQ(rq ReadRQ) alpinebits.RuleResult {
	return alpinebits.Validate(rq,
		alpinebits.RequiredHotelCode,
	)
}

// ValidateResRetrieveRS runs validation rules for OTA_ResRetrieveRS.
func ValidateResRetrieveRS(rs ResRetrieveRS) alpinebits.RuleResult {
	if rs.HotelReservations == nil {
		return alpinebits.RuleResult{}
	}

	var result alpinebits.RuleResult
	for i, res := range *rs.HotelReservations {
		ctx := resCtx{index: i}
		var r alpinebits.RuleResult

		switch res.ResStatus {
		case ResStatusCancelled:
			r = validateCancellation(res, ctx)
		case ResStatusRequested:
			r = validateRequest(res, ctx)
		case ResStatusReserved, ResStatusModify:
			r = validateReservation(res, ctx)
		}

		result.Merge(r)
	}
	return result
}

// validateCancellation validates Cancelled status.
func validateCancellation(res HotelReservation, ctx resCtx) alpinebits.RuleResult {
	return alpinebits.Validate(res,
		uniqueIDType15(ctx),
	)
}

// validateRequest validates Requested status (quote requests).
func validateRequest(res HotelReservation, ctx resCtx) alpinebits.RuleResult {
	return alpinebits.Validate(res,
		uniqueIDType14(ctx),
		requireResGuest(ctx),
		requireCustomerSurname(ctx),
		requireResGlobalInfo(ctx),
		validateRequestRoomStays(ctx),
	)
}

// validateReservation validates Reserved and Modify status.
func validateReservation(res HotelReservation, ctx resCtx) alpinebits.RuleResult {
	return alpinebits.Validate(res,
		uniqueIDType14(ctx),
		requireResGuest(ctx),
		requireCustomerSurname(ctx),
		requireResGlobalInfo(ctx),
		requireRoomStays(ctx),
		validateReservationRoomStays(ctx),
	)
}

// --- HotelReservation rules ---

func uniqueIDType14(ctx resCtx) func(HotelReservation) alpinebits.RuleResult {
	return func(res HotelReservation) alpinebits.RuleResult {
		if res.UniqueID.Type != UniqueIDTypeReservation {
			return alpinebits.RuleResult{
				Errors: []alpinebits.Error{ctx.errf(alpinebits.ErrCodeInvalidValue,
					"UniqueID Type must be 14 for ResStatus %q", res.ResStatus)},
			}
		}
		return alpinebits.RuleResult{}
	}
}

func uniqueIDType15(ctx resCtx) func(HotelReservation) alpinebits.RuleResult {
	return func(res HotelReservation) alpinebits.RuleResult {
		if res.UniqueID.Type != UniqueIDTypeCancellation {
			return alpinebits.RuleResult{
				Errors: []alpinebits.Error{ctx.errf(alpinebits.ErrCodeInvalidValue,
					"UniqueID Type must be 15 for ResStatus %q", res.ResStatus)},
			}
		}
		return alpinebits.RuleResult{}
	}
}

func requireResGuest(ctx resCtx) func(HotelReservation) alpinebits.RuleResult {
	return func(res HotelReservation) alpinebits.RuleResult {
		if res.ResGuest == nil {
			return alpinebits.RuleResult{
				Errors: []alpinebits.Error{ctx.err(alpinebits.ErrCodeRequiredField, "ResGuests required")},
			}
		}
		if res.ResGuest.Profile.Customer == nil {
			return alpinebits.RuleResult{
				Errors: []alpinebits.Error{ctx.err(alpinebits.ErrCodeRequiredField, "Customer required")},
			}
		}
		return alpinebits.RuleResult{}
	}
}

func requireCustomerSurname(ctx resCtx) func(HotelReservation) alpinebits.RuleResult {
	return func(res HotelReservation) alpinebits.RuleResult {
		if res.ResGuest == nil || res.ResGuest.Profile.Customer == nil {
			return alpinebits.RuleResult{}
		}
		if strings.TrimSpace(res.ResGuest.Profile.Customer.PersonName.Surname) == "" {
			return alpinebits.RuleResult{
				Errors: []alpinebits.Error{ctx.err(alpinebits.ErrCodeRequiredField, "Customer Surname required")},
			}
		}
		return alpinebits.RuleResult{}
	}
}

func requireResGlobalInfo(ctx resCtx) func(HotelReservation) alpinebits.RuleResult {
	return func(res HotelReservation) alpinebits.RuleResult {
		if res.ResGlobalInfo == nil {
			return alpinebits.RuleResult{
				Errors: []alpinebits.Error{ctx.err(alpinebits.ErrCodeRequiredField, "ResGlobalInfo required")},
			}
		}
		return alpinebits.RuleResult{}
	}
}

func requireRoomStays(ctx resCtx) func(HotelReservation) alpinebits.RuleResult {
	return func(res HotelReservation) alpinebits.RuleResult {
		if res.RoomStays == nil || len(*res.RoomStays) == 0 {
			return alpinebits.RuleResult{
				Errors: []alpinebits.Error{ctx.err(alpinebits.ErrCodeRequiredField, "RoomStays required")},
			}
		}
		return alpinebits.RuleResult{}
	}
}

func validateReservationRoomStays(ctx resCtx) func(HotelReservation) alpinebits.RuleResult {
	return func(res HotelReservation) alpinebits.RuleResult {
		if res.RoomStays == nil {
			return alpinebits.RuleResult{}
		}

		var result alpinebits.RuleResult
		for i, stay := range *res.RoomStays {
			sCtx := stayCtx{resIndex: ctx.index, stayIndex: i}
			result.Merge(validateReservationRoomStay(stay, sCtx))
		}
		return result
	}
}

func validateRequestRoomStays(ctx resCtx) func(HotelReservation) alpinebits.RuleResult {
	return func(res HotelReservation) alpinebits.RuleResult {
		if res.RoomStays == nil {
			return alpinebits.RuleResult{}
		}

		var result alpinebits.RuleResult
		for i, stay := range *res.RoomStays {
			sCtx := stayCtx{resIndex: ctx.index, stayIndex: i}
			result.Merge(validateRequestRoomStay(stay, sCtx))
		}
		return result
	}
}

// --- RoomStay validation ---

func validateReservationRoomStay(stay RoomStay, ctx stayCtx) alpinebits.RuleResult {
	// RoomType required with RoomTypeCode
	if stay.RoomType == nil {
		return ctx.result(alpinebits.ErrCodeRequiredField, "RoomType required")
	}
	if strings.TrimSpace(stay.RoomType.RoomTypeCode) == "" {
		return ctx.result(alpinebits.ErrCodeRequiredField, "RoomTypeCode required")
	}

	// RatePlan required with RatePlanCode and MealsIncluded
	if stay.RatePlan == nil {
		return ctx.result(alpinebits.ErrCodeRequiredField, "RatePlan required")
	}
	if strings.TrimSpace(stay.RatePlan.RatePlanCode) == "" {
		return ctx.result(alpinebits.ErrCodeRequiredField, "RatePlanCode required")
	}
	if stay.RatePlan.MealsIncluded == nil {
		return ctx.result(alpinebits.ErrCodeRequiredField, "MealsIncluded required")
	}
	if r := validateMealsIncluded(stay.RatePlan.MealsIncluded, ctx); !r.Ok() {
		return r
	}

	// GuestCounts required with total > 0
	if stay.GuestCounts == nil || len(*stay.GuestCounts) == 0 {
		return ctx.result(alpinebits.ErrCodeRequiredField, "GuestCounts required")
	}
	total := 0
	for _, gc := range *stay.GuestCounts {
		total += gc.Count
	}
	if total == 0 {
		return ctx.result(alpinebits.ErrCodeInvalidValue, "GuestCounts total must be > 0")
	}

	// TimeSpan must have Start and End
	if stay.TimeSpan.Start.IsZero() {
		return ctx.result(alpinebits.ErrCodeRequiredField, "TimeSpan Start required")
	}
	if stay.TimeSpan.End.IsZero() {
		return ctx.result(alpinebits.ErrCodeRequiredField, "TimeSpan End required")
	}
	if !stay.TimeSpan.End.After(stay.TimeSpan.Start) {
		return ctx.result(alpinebits.ErrCodeInvalidDateRange, "TimeSpan End must be after Start")
	}

	// Total required
	if stay.Total == nil {
		return ctx.result(alpinebits.ErrCodeRequiredField, "Total required")
	}
	if strings.TrimSpace(stay.Total.CurrencyCode) == "" {
		return ctx.result(alpinebits.ErrCodeRequiredField, "Total CurrencyCode required")
	}

	return alpinebits.RuleResult{}
}

func validateMealsIncluded(m *MealsIncluded, ctx stayCtx) alpinebits.RuleResult {
	if !m.MealPlanIndicator {
		return ctx.result(alpinebits.ErrCodeInvalidValue, "MealPlanIndicator must be true")
	}
	if strings.TrimSpace(m.MealPlanCodes) == "" {
		return ctx.result(alpinebits.ErrCodeRequiredField, "MealPlanCodes required")
	}
	if !slices.Contains(ValidMealPlanCodes, m.MealPlanCodes) {
		return ctx.resultf(alpinebits.ErrCodeInvalidValue, "invalid MealPlanCodes %q", m.MealPlanCodes)
	}
	return alpinebits.RuleResult{}
}

func validateRequestRoomStay(stay RoomStay, ctx stayCtx) alpinebits.RuleResult {
	hasStartEnd := !stay.TimeSpan.Start.IsZero() && !stay.TimeSpan.End.IsZero()
	hasWindow := stay.TimeSpan.StartDateWindow != nil
	hasDuration := stay.TimeSpan.Duration > 0

	// Must have either Start/End or StartDateWindow
	if stay.TimeSpan.Start.IsZero() && stay.TimeSpan.StartDateWindow == nil {
		return ctx.result(alpinebits.ErrCodeRequiredField, "TimeSpan requires Start or StartDateWindow")
	}

	// If using StartDateWindow, Duration or End should be present
	if hasWindow && !hasDuration && stay.TimeSpan.End.IsZero() {
		return ctx.result(alpinebits.ErrCodeRequiredField, "TimeSpan Duration or End required with StartDateWindow")
	}

	// Validate date ranges
	if hasStartEnd && !stay.TimeSpan.End.After(stay.TimeSpan.Start) {
		return ctx.result(alpinebits.ErrCodeInvalidDateRange, "TimeSpan End must be after Start")
	}

	if hasWindow {
		w := stay.TimeSpan.StartDateWindow
		if !w.LatestDate.IsZero() && !w.EarliestDate.IsZero() && w.LatestDate.Before(w.EarliestDate) {
			return ctx.result(alpinebits.ErrCodeInvalidDateRange, "StartDateWindow LatestDate must be >= EarliestDate")
		}
	}

	return alpinebits.RuleResult{}
}

// ValidateNotifReportRQ runs validation rules for OTA_NotifReportRQ.
func ValidateNotifReportRQ(rq NotifReportRQ) alpinebits.RuleResult {
	var result alpinebits.RuleResult
	for i, res := range rq.HotelReservations {
		if strings.TrimSpace(res.UniqueID.ID) == "" {
			result.Errors = append(result.Errors, alpinebits.ApplicationErrorf(
				alpinebits.ErrCodeRequiredField,
				"HotelReservation[%d]: UniqueID ID required", i,
			))
		}
	}
	return result
}
