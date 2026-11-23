// Package v202210 holds the Version and actions of AlpineBits 2022-10. The
// messages of each action are in a subpackage, such as freerooms.
package v202210

import (
	_ "embed"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/internal/xsd"
)

//go:embed alpinebits.xsd
var schemaFile []byte

// Version is AlpineBits 2022-10. Its XSD validation requires cgo.
var Version = alpinebits.NewVersion("2022-10", xsd.MustParse(schemaFile))
