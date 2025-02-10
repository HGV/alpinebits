package common

type TextFormat string

const (
	TextFormatPlainText = "PlainText"
	TextFormatHTML      = "HTML"
)

type Description struct {
	TextFormat TextFormat `xml:"TextFormat,attr"`
	Language   string     `xml:"Language,attr"`
	Value      string     `xml:",innerxml"`
}

type URL struct {
	Value string `xml:",innerxml"`
}
