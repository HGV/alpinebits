package main

import (
	"encoding/json"
	"testing"

	"github.com/HGV/alpinebits/v_2018_10"
	"github.com/HGV/alpinebits/v_2020_10"
	"github.com/stretchr/testify/assert"
)

func TestNewHandshakeDataFromRouter(t *testing.T) {
	r := NewRouter()

	v202010, _ := v_2020_10.NewVersion()
	r.Version(v202010, func(s *Subrouter) {
		s.Action(v_2020_10.ActionPing, nil)
		s.Action(v_2020_10.ActionHotelInvCountNotif, nil, WithCapabilities(
			v_2020_10.CapabilityHotelInvCountNotifAcceptRooms,
			v_2020_10.CapabilityHotelInvCountNotifAcceptDeltas,
			v_2020_10.CapabilityHotelInvCountNotifAcceptOutOfOrder,
			v_2020_10.CapabilityHotelInvCountNotifAcceptOutOfMarket,
			v_2020_10.CapabilityHotelInvCountNotifAcceptClosingSeasons,
		))
	})

	v201810, _ := v_2018_10.NewVersion()
	r.Version(v201810, func(s *Subrouter) {
		s.Action(v_2018_10.ActionHotelAvailNotif, nil)
		s.Action(v_2018_10.ActionReadGuestRequests, nil, WithExcludeFromHandshake())
	})

	handshakeData := NewHandshakeDataFromRouter(*r)

	expected := HandshakeData{
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
			"action_OTA_HotelAvailNotif": nil,
		},
	}

	assert.Equal(t, expected, handshakeData)
}

func TestMarshalJSON(t *testing.T) {
	handshakeData := HandshakeData{
		"2022-10": map[string][]string{},
		"2020-10": map[string][]string{
			"action_OTA_Ping": nil,
			"action_OTA_Read": {},
			"action_OTA_HotelAvailNotif": {
				"OTA_HotelAvailNotif_accept_rooms",
				"OTA_HotelAvailNotif_accept_deltas",
				"OTA_HotelAvailNotif_accept_BookingThreshold",
			},
		},
		"2018-10": map[string][]string{
			"action_OTA_Ping": nil,
		},
	}

	got, err := json.Marshal(handshakeData)
	assert.Nil(t, err)

	expected := `{
		"versions": [
			{
				"version": "2022-10"
			},
			{
				"version": "2020-10",
				"actions": [
					{
						"action": "action_OTA_HotelAvailNotif",
						"supports": [
							"OTA_HotelAvailNotif_accept_rooms",
							"OTA_HotelAvailNotif_accept_deltas",
							"OTA_HotelAvailNotif_accept_BookingThreshold"
						]
					},
					{
						"action": "action_OTA_Ping"
					},
					{
						"action": "action_OTA_Read"
					}
				]
			},
			{
				"version": "2018-10",
				"actions": [
					{
					"action": "action_OTA_Ping"
					}
				]
			}
		]
	}`

	assert.JSONEq(t, expected, string(got))
}

func TestUnmarshalJSON(t *testing.T) {
	got := `{
		"versions": [
			{
				"version": "2022-10"
			},
			{
				"version": "2020-10",
				"actions": [
					{
						"action": "action_OTA_HotelAvailNotif",
						"supports": [
							"OTA_HotelAvailNotif_accept_rooms",
							"OTA_HotelAvailNotif_accept_deltas",
							"OTA_HotelAvailNotif_accept_BookingThreshold"
						]
					},
					{
						"action": "action_OTA_Ping"
					},
					{
						"action": "action_OTA_Read"
					}
				]
			},
			{
				"version": "2018-10",
				"actions": [
					{
					"action": "action_OTA_Ping"
					}
				]
			}
		]
	}`

	var handshakeData HandshakeData
	err := json.Unmarshal([]byte(got), &handshakeData)
	assert.Nil(t, err)

	expected := HandshakeData{
		"2022-10": map[string][]string{},
		"2020-10": map[string][]string{
			"action_OTA_Ping": nil,
			"action_OTA_Read": nil,
			"action_OTA_HotelAvailNotif": {
				"OTA_HotelAvailNotif_accept_rooms",
				"OTA_HotelAvailNotif_accept_deltas",
				"OTA_HotelAvailNotif_accept_BookingThreshold",
			},
		},
		"2018-10": map[string][]string{
			"action_OTA_Ping": nil,
		},
	}

	assert.Equal(t, expected, handshakeData)
}

func TestIntersect(t *testing.T) {
	serverHandshakeData := HandshakeData{
		"2022-10": map[string][]string{},
		"2020-10": map[string][]string{
			"action_OTA_Ping": nil,
			"action_OTA_Read": nil,
			"action_OTA_HotelAvailNotif": {
				"OTA_HotelAvailNotif_accept_rooms",
				"OTA_HotelAvailNotif_accept_deltas",
				"OTA_HotelAvailNotif_accept_BookingThreshold",
			},
		},
		"2018-10": map[string][]string{
			"action_OTA_Ping": nil,
		},
	}
	clientHandshakeData := HandshakeData{
		"2020-10": map[string][]string{
			"action_OTA_Ping": nil,
			"action_OTA_HotelAvailNotif": {
				"OTA_HotelAvailNotif_accept_categories",
				"OTA_HotelAvailNotif_accept_deltas",
			},
		},
		"2018-10": map[string][]string{
			"action_OTA_Ping": nil,
		},
	}

	expected := HandshakeData{
		"2020-10": map[string][]string{
			"action_OTA_Ping": nil,
			"action_OTA_HotelAvailNotif": {
				"OTA_HotelAvailNotif_accept_deltas",
			},
		},
		"2018-10": map[string][]string{
			"action_OTA_Ping": nil,
		},
	}

	assert.Equal(t, expected, serverHandshakeData.Intersect(clientHandshakeData))
}
