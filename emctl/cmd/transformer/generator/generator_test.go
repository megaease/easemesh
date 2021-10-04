package generator

import (
	"path/filepath"
	"testing"

	"github.com/dave/jennifer/jen"
)

func TestGenerator(t *testing.T) {
	// for test
	spec := &InterfaceFileSpec{}
	// Get the package of the file with go:generate comment
	goPackage := "github.com/megaease/easemeshctl/cmd/client/command/meshclient"
	spec.Buf = jen.NewFile(goPackage)
	spec.SourceFile = "../../client/command/meshclient/canary_interface.go"
	spec.PkgName = goPackage
	ext := filepath.Ext(spec.SourceFile)
	baseFilename := spec.SourceFile[0 : len(spec.SourceFile)-len(ext)]
	spec.GenerateFileName = baseFilename + "_gen.go"
	spec.ResourceName = "canary"
	err := New(spec).Accept(NewVisitor(spec.ResourceName))
	if err != nil {
		t.Fatalf("generate code error, %s", err)
	}
}
