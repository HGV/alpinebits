// Package xsd validates XML against an XSD schema using libxml2. It is the
// only package that needs cgo; only the version packages import it.
package xsd

import (
	"strings"

	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/parser"
	libxsd "github.com/lestrrat-go/libxml2/xsd"
)

// parseOptions keep libxml2 from writing to stderr and from using the network.
// XMLParseNoEnt must stay off, or external entities in requests get resolved.
const parseOptions = parser.XMLParseNoError | parser.XMLParseNoWarning | parser.XMLParseNoNet

// ValidationError holds the individual schema violations of a document.
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

// Schema is a parsed XSD schema. It is safe for concurrent use.
type Schema struct {
	xsd *libxsd.Schema
}

// Parse parses an XSD schema.
func Parse(b []byte) (*Schema, error) {
	s, err := libxsd.Parse(b)
	if err != nil {
		return nil, err
	}
	return &Schema{xsd: s}, nil
}

// MustParse is like Parse but panics on error.
func MustParse(b []byte) *Schema {
	s, err := Parse(b)
	if err != nil {
		panic("xsd: parse schema: " + err.Error())
	}
	return s
}

// Free releases the schema's C memory.
func (s *Schema) Free() {
	if s.xsd != nil {
		s.xsd.Free()
	}
}

// Validate checks doc against the schema. Schema violations are returned as
// *ValidationError, parse errors as is.
func (s *Schema) Validate(doc []byte) error {
	d, err := libxml2.Parse(doc, parseOptions)
	if err != nil {
		return err
	}
	defer d.Free()

	if err := s.xsd.Validate(d); err != nil {
		schemaErr, ok := err.(libxsd.SchemaValidationError)
		if !ok {
			return err
		}
		return &ValidationError{Errors: schemaErr.Errors()}
	}
	return nil
}
