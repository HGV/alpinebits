package v_2022_10

import (
	_ "embed"

	"github.com/HGV/alpinebits/internal/schema"
	"github.com/HGV/alpinebits/version"
)

//go:embed alpinebits.xsd
var schemaFile []byte

type Version struct {
	schema *schema.Schema
}

var _ version.Version[Action] = new(Version)

func NewVersion() (*Version, error) {
	s, err := schema.Parse(schemaFile)
	if err != nil {
		return nil, err
	}
	return &Version{schema: s}, nil
}

func (v *Version) ValidateXML(xml string) error {
	return v.schema.Validate(xml)
}

func (v *Version) String() string {
	return "2022-10"
}
