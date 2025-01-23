package schema

import (
	"errors"
	"strings"

	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/xsd"
)

type Schema struct {
	xsd *xsd.Schema
}

func (s *Schema) Validate(xml string) error {
	d, err := libxml2.ParseString(xml)
	if err != nil {
		return err
	}
	defer d.Free()

	if err := s.xsd.Validate(d); err != nil {
		var errs []string
		for _, err := range err.(xsd.SchemaValidationError).Errors() {
			errs = append(errs, err.Error())
		}
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

func Parse(buf []byte) (*Schema, error) {
	s, err := xsd.Parse(buf)
	if err != nil {
		return nil, err
	}
	return &Schema{xsd: s}, nil
}
