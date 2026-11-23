// Package v202010 holds the Version and actions of AlpineBits 2020-10. The
// messages of each action are in a subpackage, such as freerooms.
package v202010

import (
	_ "embed"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/internal/xsd"
)

//go:embed alpinebits.xsd
var schemaFile []byte

// Version is AlpineBits 2020-10. Its XSD validation requires cgo.
var Version = alpinebits.NewVersion("2020-10", xsd.MustParse(schemaFile))
