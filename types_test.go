package alpinebits

import (
	"encoding/xml"
	"testing"
)

func TestDecimal(t *testing.T) {
	type amounts struct {
		XMLName xml.Name `xml:"A"`
		Big     Decimal  `xml:"Big,attr"`
		Small   *Decimal `xml:"Small,attr,omitempty"`
		Zero    Decimal  `xml:"Zero,attr,omitempty"`
	}
	out, err := xml.Marshal(amounts{Big: 1500000, Small: new(Decimal(0.00001))})
	if want := `<A Big="1500000" Small="0.00001"></A>`; err != nil || string(out) != want {
		t.Errorf("marshal: %s, %v; want %s", out, err, want)
	}
	var back amounts
	if err := xml.Unmarshal([]byte(`<A Big=" 76.80 "/>`), &back); err != nil || back.Big != 76.8 || back.Small != nil {
		t.Errorf("unmarshal: %+v, %v", back, err)
	}
}
