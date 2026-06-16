package alpinebits

// InfoCode identifies the type of multimedia content (OTA INF list subset).
type InfoCode int

const (
	InfoCodeDescription InfoCode = 1  // Description
	InfoCodePictures    InfoCode = 23 // Pictures
	InfoCodeLongName    InfoCode = 25 // Long name (title)
)

// MultimediaDescriptions is a collection of multimedia descriptions.
type MultimediaDescriptions []MultimediaDescription

// ByInfoCode returns the MultimediaDescription with the given InfoCode, or nil if not found.
func (mds MultimediaDescriptions) ByInfoCode(code InfoCode) *MultimediaDescription {
	for i, md := range mds {
		if md.InfoCode == code {
			return &mds[i]
		}
	}
	return nil
}

type MultimediaDescription struct {
	InfoCode   InfoCode    `xml:"InfoCode,attr"`
	TextItems  *TextItem   `xml:"TextItems>TextItem"`
	ImageItems []ImageItem `xml:"ImageItems>ImageItem"`
}

type TextItem struct {
	Descriptions []Description `xml:"Description"`
}

// TextFormat specifies the format of description text.
type TextFormat string

const (
	TextFormatPlainText TextFormat = "PlainText"
	TextFormatHTML      TextFormat = "HTML"
)

type Description struct {
	TextFormat TextFormat `xml:"TextFormat,attr"`
	Language   string     `xml:"Language,attr"`
	Value      string     `xml:",chardata"`
}

type ImageItem struct {
	Category     int           `xml:"Category,attr"`
	ImageFormat  ImageFormat   `xml:"ImageFormat"`
	Descriptions []Description `xml:"Description"`
}

type ImageFormat struct {
	URL             string `xml:"URL"`
	CopyrightNotice string `xml:"CopyrightNotice,attr,omitempty"`
}
