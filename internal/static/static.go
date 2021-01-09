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

	if dogmaPkg := prog.ImportedPackage(dogmaPkgPath); dogmaPkg != nil {
		return fromSSAPackages(dogmaPkg, packages)
	}

	// If Dogma package is not found as an import, none of the packages in the
	// built application implement dogma.Application interface that requires the
	// import of the github.com/dogmatiq/dogma package.
	return nil
}

// fromSSAPackages loads the list of configkit.Application from the given list
// of ssa.Package items.
func fromSSAPackages(dogmaPkg *ssa.Package, pkgs []*ssa.Package) []configkit.Application {
	var result []configkit.Application

	for _, pkg := range pkgs {
		// The package can be returned as nil by ssautil.AllPackages() if it
		// previously contained errors in types.Package.Errors.
		if pkg == nil {
			continue
		}

		a := dogmaPkg.Pkg.Scope().Lookup("Application")
		iface := a.Type().Underlying().(*types.Interface)

		for _, m := range pkg.Members {
			// NOTE: the sequence of the if-blocks below is important as the
			// value of a type implements the interface only if the
			// interface-implementing methods have value receiver. Hence is
			// the implementation check for the type value first.
			//
			// However, a pointer to the type implements the interface in
			// both cases: when interface-implementing methods have value or
			// pointer receiver.
			if types.Implements(m.Type(), iface) {
				result = append(
					result,
					parse(m, m.Type()),
				)
				continue
			}

			if p := types.NewPointer(m.Type()); types.Implements(p, iface) {
				result = append(
					result,
					parse(m, p),
				)
			}
		}
	}

	return result
}
