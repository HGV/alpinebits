package inventory

import (
	"slices"
	"strings"

	"github.com/HGV/alpinebits"
)

// Capabilities for Inventory action.
const (
	CapUseRooms          alpinebits.Capability = "OTA_HotelDescriptiveContentNotif_Inventory_use_rooms"
	CapOccupancyChildren alpinebits.Capability = "OTA_HotelDescriptiveContentNotif_Inventory_occupancy_children"
)

// ValidateOptions holds parameters for validation.
type ValidateOptions struct {
	Caps []alpinebits.Capability
}

// Validate runs all validation rules for HotelDescriptiveContentNotifRQ.
// Note: Basic field requirements (StandardOccupancy > 0, RoomType values, etc.)
// are validated by the XSD schema before these rules run.
func Validate(rq HotelDescriptiveContentNotifRQ, opts ValidateOptions) alpinebits.RuleResult {
	roomTypes, individualRooms := Partition(rq.HotelDescriptiveContent.GuestRooms)

	return alpinebits.Validate(rq,
		alpinebits.RequiredHotelCode,
		requiredGuestRoomCode,
		validTypeRoom(roomTypes),
		validOccupancy(roomTypes, opts.Caps),
		validRooms(individualRooms, opts.Caps),
		validMultimediaDescriptions(roomTypes),
	)
}

// Partition separates GuestRooms into room types and individual rooms.
// Room types are the first occurrence of each Code.
// Individual rooms are subsequent occurrences with the same Code.
func Partition(rooms []GuestRoom) (roomTypes, individualRooms []GuestRoom) {
	seen := make(map[string]bool)
	for _, room := range rooms {
		if seen[room.Code] {
			individualRooms = append(individualRooms, room)
		} else {
			roomTypes = append(roomTypes, room)
			seen[room.Code] = true
		}
	}
	return
}

func requiredGuestRoomCode(rq HotelDescriptiveContentNotifRQ) alpinebits.RuleResult {
	var result alpinebits.RuleResult
	for i, room := range rq.HotelDescriptiveContent.GuestRooms {
		if strings.TrimSpace(room.Code) == "" {
			result.Errors = append(result.Errors, alpinebits.ApplicationErrorf(
				alpinebits.ErrCodeRequiredField,
				"GuestRoom[%d]: missing Code", i,
			))
		}
	}
	return result
}

// validRoomTypeForClassification maps RoomClassificationCode to valid RoomType values.
var validRoomTypeForClassification = map[int][]int{
	5:  {6, 7, 8},    // camping: camping grounds, pitches, camping grounds/pitches
	13: {2, 3, 4, 5}, // apartment/holiday home: apartments, mobile homes, bungalows, holiday homes
	42: {1, 9},       // room: rooms, resting places
}

// validTypeRoom validates RoomType against RoomClassificationCode.
// Only certain RoomType values are valid for each RoomClassificationCode.
func validTypeRoom(roomTypes []GuestRoom) func(HotelDescriptiveContentNotifRQ) alpinebits.RuleResult {
	return func(_ HotelDescriptiveContentNotifRQ) alpinebits.RuleResult {
		var result alpinebits.RuleResult
		for _, room := range roomTypes {
			tr := room.TypeRoom
			if tr.RoomType == 0 {
				continue
			}

			validTypes, ok := validRoomTypeForClassification[tr.RoomClassificationCode]
			if !ok {
				result.Errors = append(result.Errors, alpinebits.ApplicationErrorf(
					alpinebits.ErrCodeInvalidValue,
					"GuestRoom %q: RoomType not allowed for RoomClassificationCode %d", room.Code, tr.RoomClassificationCode,
				))
			} else if !slices.Contains(validTypes, tr.RoomType) {
				result.Errors = append(result.Errors, alpinebits.ApplicationErrorf(
					alpinebits.ErrCodeInvalidValue,
					"GuestRoom %q: invalid RoomType %d for RoomClassificationCode %d", room.Code, tr.RoomType, tr.RoomClassificationCode,
				))
			}
		}
		return result
	}
}

func validOccupancy(roomTypes []GuestRoom, caps []alpinebits.Capability) func(HotelDescriptiveContentNotifRQ) alpinebits.RuleResult {
	return func(_ HotelDescriptiveContentNotifRQ) alpinebits.RuleResult {
		supportsOccupancyChildren := slices.Contains(caps, CapOccupancyChildren)

		var result alpinebits.RuleResult
		for _, room := range roomTypes {
			if !supportsOccupancyChildren && room.MaxChildOccupancy > 0 {
				result.Errors = append(result.Errors, alpinebits.ApplicationErrorf(
					alpinebits.ErrCodeInvalidValue,
					"GuestRoom %q: MaxChildOccupancy not supported", room.Code,
				))
			}
			if room.MaxOccupancy > 0 && room.MinOccupancy > room.MaxOccupancy {
				result.Errors = append(result.Errors, alpinebits.ApplicationErrorf(
					alpinebits.ErrCodeInvalidValue,
					"GuestRoom %q: MinOccupancy exceeds MaxOccupancy", room.Code,
				))
			}
			if room.MaxOccupancy > 0 && room.MaxChildOccupancy > room.MaxOccupancy {
				result.Errors = append(result.Errors, alpinebits.ApplicationErrorf(
					alpinebits.ErrCodeInvalidValue,
					"GuestRoom %q: MaxChildOccupancy exceeds MaxOccupancy", room.Code,
				))
			}
			if room.MaxOccupancy > 0 && room.TypeRoom.StandardOccupancy > room.MaxOccupancy {
				result.Errors = append(result.Errors, alpinebits.ApplicationErrorf(
					alpinebits.ErrCodeInvalidValue,
					"GuestRoom %q: StandardOccupancy exceeds MaxOccupancy", room.Code,
				))
			}
		}
		return result
	}
}

func validRooms(individualRooms []GuestRoom, caps []alpinebits.Capability) func(HotelDescriptiveContentNotifRQ) alpinebits.RuleResult {
	return func(_ HotelDescriptiveContentNotifRQ) alpinebits.RuleResult {
		if len(individualRooms) > 0 && !slices.Contains(caps, CapUseRooms) {
			return alpinebits.RuleResult{
				Errors: []alpinebits.Error{
					alpinebits.ApplicationError(alpinebits.ErrCodeInvalidValue, "individual rooms not supported"),
				},
			}
		}
		return alpinebits.RuleResult{}
	}
}

func validMultimediaDescriptions(roomTypes []GuestRoom) func(HotelDescriptiveContentNotifRQ) alpinebits.RuleResult {
	return func(_ HotelDescriptiveContentNotifRQ) alpinebits.RuleResult {
		var result alpinebits.RuleResult
		for _, room := range roomTypes {
			var longNameCount, descriptionCount, picturesCount int
			for _, md := range room.MultimediaDescriptions {
				switch md.InfoCode {
				case alpinebits.InfoCodeLongName:
					longNameCount++
				case alpinebits.InfoCodeDescription:
					descriptionCount++
				case alpinebits.InfoCodePictures:
					picturesCount++
				}
			}

			if longNameCount == 0 {
				result.Errors = append(result.Errors, alpinebits.ApplicationErrorf(
					alpinebits.ErrCodeRequiredField,
					"GuestRoom %q: missing MultimediaDescription with InfoCode 25 (Long name)", room.Code,
				))
			} else if longNameCount > 1 {
				result.Errors = append(result.Errors, alpinebits.ApplicationErrorf(
					alpinebits.ErrCodeInvalidValue,
					"GuestRoom %q: multiple MultimediaDescription with InfoCode 25 (Long name)", room.Code,
				))
			}
			if descriptionCount > 1 {
				result.Errors = append(result.Errors, alpinebits.ApplicationErrorf(
					alpinebits.ErrCodeInvalidValue,
					"GuestRoom %q: multiple MultimediaDescription with InfoCode 1 (Description)", room.Code,
				))
			}
			if picturesCount > 1 {
				result.Errors = append(result.Errors, alpinebits.ApplicationErrorf(
					alpinebits.ErrCodeInvalidValue,
					"GuestRoom %q: multiple MultimediaDescription with InfoCode 23 (Pictures)", room.Code,
				))
			}
		}
		return result
	}
}
