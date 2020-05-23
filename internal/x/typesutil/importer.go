package typesutil

import (
	"go/importer"
	"go/types"
)

type pkgimporter struct {
	def types.Importer // default
	src types.Importer // fallback
}

// NewImporter initializes and returns a new instance of types.Importer.
func NewImporter() types.Importer {
	return &pkgimporter{
		def: importer.Default(),
		src: importer.For("source", nil),
	}
}

// Import returns the imported package for the given import path. Import first
// attempts to import the package using the default importer, if that fails
// another attempt is made to import the package using a "source" importer,
// but if that fails as well an error will be returned.
func (i pkgimporter) Import(path string) (*types.Package, error) {
	//pkg, err := i.def.Import(path)
	//if err != nil {
	//	return i.src.Import(path)
	//}
	//return pkg, nil

	return i.src.Import(path)
}
