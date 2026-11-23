// Package v201810 holds the Version and actions of AlpineBits 2018-10. The
// messages of each action are in a subpackage, such as freerooms.
package v201810

import (
	_ "embed"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/internal/xsd"
)

//go:embed alpinebits.xsd
var schemaFile []byte

// Version is AlpineBits 2018-10. Its XSD validation requires cgo.
var Version = alpinebits.NewVersion("2018-10", xsd.MustParse(schemaFile))
