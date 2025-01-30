package handshake

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
