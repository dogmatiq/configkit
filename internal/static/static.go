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
			// Any packages.Package that can not be built results in a nil
			// ssa.Package. We ignore any such packages so that we can still
			// obtain information about applications from any valid packages.
			continue
		}

		for _, m := range pkg.Members {
			// The sequence of the if-blocks below is important as a type
			// implements an interface only if the methods in the interface's
			// method set have non-pointer receivers. Hence the implementation
			// check for the "raw" (non-pointer) type is made first.
			//
			// A pointer to the type, on the other hand, implements the
			// interface regardless of whether pointer receivers are used or
			// not.
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
