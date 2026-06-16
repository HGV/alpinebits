package v202010

import (
	_ "embed"

	"github.com/HGV/alpinebits"
)

//go:embed alpinebits.xsd
var schemaFile []byte

var Version = alpinebits.MustNewVersion("2020-10", schemaFile)
