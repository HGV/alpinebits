// Package alpinebits implements the AlpineBits protocol: a Router that serves
// its actions over HTTP and a Client that sends them. The versions and their
// actions are in subpackages, such as v201810.
package alpinebits

// Validator checks an XML document against a version's schema. It must be
// safe for concurrent use.
type Validator interface {
	Validate(doc []byte) error
}

// Version is an AlpineBits protocol version.
type Version struct {
	name      string
	validator Validator
}

// NewVersion creates a version whose documents are checked by v, which must
// not be nil.
func NewVersion(name string, v Validator) *Version {
	if v == nil {
		panic("alpinebits: NewVersion(" + name + "): nil Validator")
	}
	return &Version{name: name, validator: v}
}

// Name returns the version identifier (e.g., "2018-10").
func (v *Version) Name() string {
	return v.name
}

// Validate checks doc against the version's schema.
func (v *Version) Validate(doc []byte) error {
	return v.validator.Validate(doc)
}

// xmlBody reports whether requests of v are an XML document with the action in
// the X-AlpineBits-Action header (AlpineBits 2026-10 §2.1) rather than
// multipart/form-data.
func (v *Version) xmlBody() bool {
	return v.name >= "2026-10"
}
