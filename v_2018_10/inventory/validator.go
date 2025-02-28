package inventory

import (
	"strings"

	"github.com/HGV/alpinebits/v_2018_10/common"
	"github.com/HGV/x/slicesx"
)

type HotelDescriptiveContentNotifValidator struct {
	supportsRooms             bool
	supportsOccupancyChildren bool
}

var _ common.Validatable[HotelDescriptiveContentNotifRQ] = (*HotelDescriptiveContentNotifValidator)(nil)

type HotelDescriptiveContentNotifValidatorFunc func(*HotelDescriptiveContentNotifValidator)

func NewHotelDescriptiveContentNotifValidator(opts ...HotelDescriptiveContentNotifValidatorFunc) HotelDescriptiveContentNotifValidator {
	var v HotelDescriptiveContentNotifValidator
	for _, opt := range opts {
		opt(&v)
	}
	return v
}

func WithRooms() HotelDescriptiveContentNotifValidatorFunc {
	return func(v *HotelDescriptiveContentNotifValidator) {
		v.supportsRooms = true
	}
}

func WithOccupancyChildren() HotelDescriptiveContentNotifValidatorFunc {
	return func(v *HotelDescriptiveContentNotifValidator) {
		v.supportsOccupancyChildren = true
	}
}

func (v HotelDescriptiveContentNotifValidator) Validate(r HotelDescriptiveContentNotifRQ) error {
	if err := common.ValidateHotelCode(r.HotelDescriptiveContent.HotelCode); err != nil {
		return err
	}

	if err := v.validateGuestRooms(r.HotelDescriptiveContent.GuestRooms); err != nil {
		return err
	}

	return nil
}

func (v HotelDescriptiveContentNotifValidator) validateGuestRooms(guestRooms []GuestRoom) error {
	guestRoomsByCode := slicesx.GroupByFunc(guestRooms, func(g GuestRoom) string {
		return g.Code
	})
	for _, guestRooms := range guestRoomsByCode {
		if err := v.validateGuestRoom(guestRooms); err != nil {
			return err
		}
	}
	return nil
}

func (v HotelDescriptiveContentNotifValidator) validateGuestRoom(guestRooms []GuestRoom) error {
	headGuestRoom := guestRooms[0]
	if strings.TrimSpace(headGuestRoom.Code) == "" {
		return common.ErrMissingCode
	}

	if err := v.validateOccupancies(headGuestRoom); err != nil {
		return err
	}

	if err := v.validateTypeRoom(headGuestRoom.TypeRoom); err != nil {
		return err
	}

	if err := v.validateAmenities(headGuestRoom.Amenities); err != nil {
		return err
	}

	if err := v.validateMultimediaDescriptions(headGuestRoom.MultimediaDescriptions); err != nil {
		return err
	}

	tailGuestRooms := guestRooms[1:]
	if err := v.validateRooms(tailGuestRooms); err != nil {
		return err
	}

	return nil
}

func (v HotelDescriptiveContentNotifValidator) validateOccupancies(guestRoom GuestRoom) error {
	min := guestRoom.MinOccupancy
	std := guestRoom.TypeRoom.StandardOccupancy
	max := guestRoom.MaxOccupancy
	mco := guestRoom.MaxChildOccupancy

	if !v.supportsOccupancyChildren && mco > 0 {
		return common.ErrChildOccupancyNotSupported
	}

	if mco > max {
		return common.ErrMaxChildOccGreaterThanMaxOcc
	}

	if std < min {
		return common.ErrStdOccLowerThanMinOcc
	}

	if max < std {
		return common.ErrMaxOccLowerThanStdOcc
	}

	return nil
}

func (v *HotelDescriptiveContentNotifValidator) validateTypeRoom(typeRoom TypeRoom) error {
	if typeRoom.RoomClassificationCode < 1 || typeRoom.RoomClassificationCode > 83 {
		return common.ErrInvalidRoomClassificationCode(typeRoom.RoomClassificationCode)
	}

	if typeRoom.RoomType > 0 {
		allowed := map[int]int{
			1: 42, // Room
			2: 13, // Apartments
			3: 13, // Mobile Homes
			4: 13, // Bungalows
			5: 13, // Holiday Homes
			6: 5,  // Camping Grounds
			7: 5,  // Pitches
			8: 5,  // Camping Grounds/Pitches
			9: 42, // Resting places
		}
		rcc, ok := allowed[typeRoom.RoomType]
		if !ok {
			return common.ErrInvalidRoomType(typeRoom.RoomType)
		}
		if typeRoom.RoomClassificationCode != rcc {
			return common.ErrInvalidRoomClassificationCode(typeRoom.RoomClassificationCode)
		}
	}

	return nil
}

func (v *HotelDescriptiveContentNotifValidator) validateAmenities(amenities *[]Amenity) error {
	if amenities == nil {
		return nil
	}

	for _, amenity := range *amenities {
		if code := amenity.RoomAmenityCode; code < 1 || code > 293 {
			return common.ErrInvalidRoomAmenityType(code)
		}
	}

	return nil
}

func (v *HotelDescriptiveContentNotifValidator) validateMultimediaDescriptions(mds MultimediaDescriptions) error {
	longNames := mds.LongNames()
	if len(longNames) == 0 {
		return common.ErrMissingLongName
	}

	for _, md := range mds {
		switch md.InfoCode {
		case InformationTypeLongName:
			if err := common.ValidateLanguageUniqueness(*md.TextItems); err != nil {
				return err
			}
		case InformationTypeDescription:
			if err := common.ValidateLanguageUniqueness(*md.TextItems); err != nil {
				return err
			}
		case InformationTypePictures:
			if err := v.validateImages(*md.ImageItems); err != nil {
				return err
			}
		}
	}

	return nil
}

func (v *HotelDescriptiveContentNotifValidator) validateImages(images []ImageItem) error {
	for _, image := range images {
		if category := image.Category; category < 1 || category > 23 {
			return common.ErrInvalidPictureCategoryCode(category)
		}
		if err := common.ValidateLanguageUniqueness(image.Descriptions); err != nil {
			return err
		}
	}
	return nil
}

func (v *HotelDescriptiveContentNotifValidator) validateRooms(rooms []GuestRoom) error {
	if !v.supportsRooms && len(rooms) > 0 {
		return common.ErrRoomsNotSupported
	}

	for _, room := range rooms {
		if strings.TrimSpace(room.TypeRoom.RoomID) == "" {
			return common.ErrMissingRoomID
		}
	}

	return nil
}
