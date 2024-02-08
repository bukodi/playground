package gotools

import (
	"golang.org/x/tools/go/packages"
	"os"
	"testing"
)

func TestListPackages(t *testing.T) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedDeps |
			packages.NeedExportFile | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedTypesSizes |
			packages.NeedModule | packages.NeedEmbedFiles | packages.NeedEmbedPatterns,
		Dir:   "",
		Tests: true,
	}
	pkgs, err := packages.Load(cfg)
	if err != nil {
		t.Fatalf("load: %v\n", err)
		os.Exit(1)
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	// Print the names of the source files
	// for each package listed on the command line.
	for _, pkg := range pkgs {
		t.Logf("%s, %v", pkg.ID, pkg.GoFiles)
	}
}
