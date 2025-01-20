package v_2018_10

type ErrorWarningType int

const (
	ErrorWarningTypeAdvisory         ErrorWarningType = 11
	ErrorWarningTypeApplicationError ErrorWarningType = 13
)

type Status string

const (
	StatusSendInventory Status = "ALPINEBITS_SEND_INVENTORY"
	StatusSendFreeRooms Status = "ALPINEBITS_SEND_FREEROOMS"
	StatusSendRatePlans Status = "ALPINEBITS_SEND_RATEPLANS"
)

type Success struct{}

type Warning struct {
	Type   ErrorWarningType `xml:"Type,attr"`
	Code   int              `xml:"Code,attr,omitempty"`
	Status Status           `xml:"Status,attr,omitempty"`
	Value  string           `xml:",innerxml"`
}

type Error struct {
	Type   ErrorWarningType `xml:"Type,attr"`
	Code   int              `xml:"Code,attr,omitempty"`
	Status Status           `xml:"Status,attr,omitempty"`
	Value  string           `xml:",innerxml"`
}

func (err Error) Error() string {
	return err.Value
}

type response struct {
	Success  *Success   `xml:"Success"`
	Warnings *[]Warning `xml:"Warnings>Warning"`
	Errors   *[]Error   `xml:"Errors>Error"`
}
