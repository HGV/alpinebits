// Package v202610 holds the Version and actions of AlpineBits 2026-10. The
// messages of each action are in a subpackage, such as freerooms.
package v202610

import (
	_ "embed"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/internal/xsd"
)

//go:embed alpinebits.xsd
var schemaFile []byte

// Version is AlpineBits 2026-10. Its XSD validation requires cgo.
var Version = alpinebits.NewVersion("2026-10", xsd.MustParse(schemaFile))
