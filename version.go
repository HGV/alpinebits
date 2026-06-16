package alpinebits

import "github.com/HGV/alpinebits/internal/schema"

type Version struct {
	name   string
	schema *schema.Schema
}

// NewVersion creates a new version with the given name and XSD schema.
func NewVersion(name string, schemaFile []byte) (*Version, error) {
	s, err := schema.Parse(schemaFile)
	if err != nil {
		return nil, err
	}
	return &Version{name: name, schema: s}, nil
}

// MustNewVersion is like NewVersion but panics on error.
// Use this for package-level version initialization.
func MustNewVersion(name string, schemaFile []byte) *Version {
	v, err := NewVersion(name, schemaFile)
	if err != nil {
		panic("alpinebits: failed to create version " + name + ": " + err.Error())
	}
	return v
}

// Name returns the version identifier (e.g., "2018-10").
func (v *Version) Name() string {
	return v.name
}

// Validate validates an XML string against the version's schema.
func (v *Version) Validate(xml string) error {
	return v.schema.Validate(xml)
}

// Free releases the underlying schema resources.
// For package-level versions, this is typically not needed as they live
// for the program's lifetime.
func (v *Version) Free() {
	v.schema.Free()
}
