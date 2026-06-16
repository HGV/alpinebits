package v202210

import (
	_ "embed"

	"github.com/HGV/alpinebits"
)

//go:embed alpinebits.xsd
var schemaFile []byte

var Version = alpinebits.MustNewVersion("2022-10", schemaFile)
