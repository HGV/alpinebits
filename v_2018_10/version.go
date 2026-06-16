package v201810

import (
	_ "embed"

	"github.com/HGV/alpinebits"
)

//go:embed alpinebits.xsd
var schemaFile []byte

var Version = alpinebits.MustNewVersion("2018-10", schemaFile)
