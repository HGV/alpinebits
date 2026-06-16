package schema

import (
	"strings"

	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/xsd"
)

// ValidationError holds structured validation errors from XSD validation.
type ValidationError struct {
	Errors []error
}

func (e *ValidationError) Error() string {
	msgs := make([]string, len(e.Errors))
	for i, err := range e.Errors {
		msgs[i] = err.Error()
	}
	return strings.Join(msgs, "\n")
}

// Schema wraps an XSD schema for XML validation.
type Schema struct {
	xsd *xsd.Schema
}

// Parse parses XSD content and returns a Schema.
func Parse(buf []byte) (*Schema, error) {
	s, err := xsd.Parse(buf)
	if err != nil {
		return nil, err
	}
	return &Schema{xsd: s}, nil
}

// Free releases the underlying XSD schema resources.
func (s *Schema) Free() {
	if s.xsd != nil {
		s.xsd.Free()
	}
}

// Validate validates an XML string against the schema.
// Returns a *ValidationError if validation fails, allowing access to individual errors.
func (s *Schema) Validate(xml string) error {
	doc, err := libxml2.ParseString(xml)
	if err != nil {
		return err
	}
	defer doc.Free()

	if err := s.xsd.Validate(doc); err != nil {
		schemaErr, ok := err.(xsd.SchemaValidationError)
		if !ok {
			return err
		}
		return &ValidationError{Errors: schemaErr.Errors()}
	}
	return nil
}
