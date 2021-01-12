package static

import (
	"go/types"

	"github.com/dogmatiq/configkit"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

const (
	// dogmaPkgPath is the full path of dogma package.
	dogmaPkgPath = "github.com/dogmatiq/dogma"
)

// FromPackages loads the list of configkit.Application from the given list
// of packages.Package items.
func FromPackages(pkgs []*packages.Package) []configkit.Application {
	prog, packages := ssautil.AllPackages(pkgs, ssa.SanityCheckFunctions)
	prog.Build()

	dogmaPkg := prog.ImportedPackage(dogmaPkgPath)
	if dogmaPkg == nil {
		// If Dogma package is not found as an import, none of the packages in
		// the built application implement dogma.Application interface.
		return nil
	}

	var apps []configkit.Application
	iface := dogmaPkg.Pkg.Scope().Lookup("Application").Type().Underlying().(*types.Interface)

	for _, pkg := range packages {
		if pkg == nil {
			continue
		}

		for _, m := range pkg.Members {
			// NOTE: the sequence of the if-blocks below is important as the
			// value of a type implements the interface only if the
			// interface-implementing methods have value receiver. Hence is the
			// implementation check for the type value is positioned first.
			//
			// A pointer to the type, on the other hand, implements the
			// interface in both cases: when interface-implementing methods have
			// value or pointer receiver.
			if types.Implements(m.Type(), iface) {
				apps = append(apps, parse(m, m.Type()))
				continue
			}

			if p := types.NewPointer(m.Type()); types.Implements(p, iface) {
				apps = append(apps, parse(m, p))
			}
		}
	}

	return apps
}
