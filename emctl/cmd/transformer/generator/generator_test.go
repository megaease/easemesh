package generator

import (
	"bytes"
	"go/ast"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/dave/jennifer/jen"
	utiltesting "k8s.io/client-go/util/testing"
)

func TestGenerator(t *testing.T) {
	// for test
	spec := &InterfaceFileSpec{}
	// Get the package of the file with go:generate comment
	goPackage := "github.com/megaease/easemeshctl/cmd/client/command/meshclient"
	spec.Buf = jen.NewFile(goPackage)
	spec.SourceFile = "../../client/command/meshclient/canary.go"
	spec.PkgName = goPackage
	ext := filepath.Ext(spec.SourceFile)
	baseFilename := spec.SourceFile[0 : len(spec.SourceFile)-len(ext)]
	spec.GenerateFileName = baseFilename + "_gen.go"
	spec.ResourceType = "canary"
	err := New(spec).Accept(NewVisitor(spec.ResourceType))
	if err != nil {
		t.Fatalf("generate code error, %s", err)
	}
}

func TestQualRender(t *testing.T) {
	q := jen.Qual("resource", "Canary")
	buf := &bytes.Buffer{}
	err := q.Render(buf)
	if err != nil {
		t.Fatalf("render buf error %s", err)
	}
	t.Logf("render result:%s", string(buf.Bytes()))
}

func TestExtractResourceName(t *testing.T) {
	tmpDir, err := utiltesting.MkTmpdir("pkgdir")
	if err != nil {
		t.Fatalf("mk test dir failed:%s", err)
	}

	sourceFile := filepath.Join(tmpDir, "observability.go")
	err = ioutil.WriteFile(sourceFile, []byte(source), 0600)
	if err != nil {
		t.Fatalf("write sourcefile error:%s", err)
	}

	interfaceMethod := func(concreateStruct string, verb Verb, method *ast.Field, imports []*ast.ImportSpec, spec *InterfaceFileSpec) error {
		var arguments, results []jen.Code
		var funcName string
		err = covertFuncType(method, imports).
			extractArguments(&arguments).
			extractResults(&results).
			extractFuncName(&funcName).
			error()

		switch Verb(funcName) {
		case Get:
			fallthrough
		case List:
			resource, err := readMethodFetcher(arguments, results)
			if err != nil {
				t.Errorf("readMethodFetch error %s", err)
			}
			if resource != "ObservabilityOutputServer" {
				t.Errorf("resource should be ObservabilityOutputServer but is %s", resource)
			}
		case Create:
			fallthrough
		case Patch:
			resource, err := writeMethodFetcher(arguments, results)
			if err != nil {
				t.Errorf("writeMethodFetch error %s", err)
			}
			if resource != "ObservabilityOutputServer" {
				t.Errorf("method %s resource should be ObservabilityOutputServer but is %s", funcName, resource)
			}
		case Delete:
			resource, err := deleteMethodFetcher(concreateStruct)(arguments, results)
			if err != nil {
				t.Errorf("deleteMethodFetch error %s", err)
			}
			if resource != "ObservabilityOutputServer" {
				t.Errorf("method %s resource should be resource but is %s", funcName, resource)
			}
		default:
			t.Fatalf("unknown funcName %s", funcName)

		}

		return nil
	}
	spec := &InterfaceFileSpec{SourceFile: sourceFile, Buf: jen.NewFile("hhh")}
	g := New(spec)
	g.Accept(&testVisitorAdaptor{interfaceMethod: interfaceMethod})

}

type (
	begin                func(imports []*ast.ImportSpec, spec *InterfaceFileSpec) error
	getterConcreate      func(name string, spec *InterfaceFileSpec) error
	interfaceConcreate   func(name string, spec *InterfaceFileSpec) error
	resourceGetterMethod func(name string, method *ast.Field, imports []*ast.ImportSpec, spec *InterfaceFileSpec) error
	interfaceMethod      func(concreateStruct string, verb Verb, method *ast.Field, imports []*ast.ImportSpec, spec *InterfaceFileSpec) error
	end                  func(buf *jen.File) error
	onErr                func(e error)
)

func (b begin) do(imports []*ast.ImportSpec, spec *InterfaceFileSpec) error {
	if b != nil {
		return b(imports, spec)
	}
	return nil
}
func (i getterConcreate) do(name string, spec *InterfaceFileSpec) error {
	if i != nil {
		return i(name, spec)
	}
	return nil
}
func (i interfaceConcreate) do(name string, spec *InterfaceFileSpec) error {
	if i != nil {
		return i(name, spec)
	}
	return nil
}
func (i resourceGetterMethod) do(name string, method *ast.Field, imports []*ast.ImportSpec, spec *InterfaceFileSpec) error {
	if i != nil {
		return i(name, method, imports, spec)
	}
	return nil
}
func (i interfaceMethod) do(concreateStruct string, verb Verb, method *ast.Field, imports []*ast.ImportSpec, spec *InterfaceFileSpec) error {
	if i != nil {
		i(concreateStruct, verb, method, imports, spec)
	}
	return nil
}
func (i end) do(buf *jen.File) error {
	if i != nil {
		i(buf)
	}
	return nil
}
func (i onErr) do(e error) {
	if i != nil {
		i(e)
	}
}

type testVisitorAdaptor struct {
	begin                begin
	getterConcreate      getterConcreate
	interfaceConcreate   interfaceConcreate
	resourceGetterMethod resourceGetterMethod
	interfaceMethod      interfaceMethod
	end                  end
	onErr                onErr
}

func (t *testVisitorAdaptor) visitorBegin(imports []*ast.ImportSpec, spec *InterfaceFileSpec) error {
	return t.begin.do(imports, spec)
}
func (t *testVisitorAdaptor) visitorResourceGetterConcreatStruct(name string, spec *InterfaceFileSpec) error {
	return t.getterConcreate.do(name, spec)
}
func (t *testVisitorAdaptor) visitorInterfaceConcreatStruct(name string, spec *InterfaceFileSpec) error {

	return t.interfaceConcreate.do(name, spec)
}
func (t *testVisitorAdaptor) visitorResourceGetterMethod(name string, method *ast.Field, imports []*ast.ImportSpec, spec *InterfaceFileSpec) error {
	return t.resourceGetterMethod.do(name, method, imports, spec)
}
func (t *testVisitorAdaptor) visitorIntrefaceMethod(concreateStruct string, verb Verb, method *ast.Field, imports []*ast.ImportSpec, spec *InterfaceFileSpec) error {
	return t.interfaceMethod.do(concreateStruct, verb, method, imports, spec)
}
func (t *testVisitorAdaptor) visitorEnd(spec *InterfaceFileSpec) error {
	return t.end.do(spec.Buf)
}
func (t *testVisitorAdaptor) onError(e error) {
	t.onErr.do(e)
}

var source = `
package meshclient

import (
	"context"

	"github.com/megaease/easemeshctl/cmd/client/resource"
)

// ObservabilityGetter represents an Observability resource accessor
type ObservabilityGetter interface {
	ObservabilityTracings() ObservabilityTracingInterface
	ObservabilityMetrics() ObservabilityMetricInterface
	ObservabilityOutputServer() ObservabilityOutputServerInterface
}

// ObservabilityOutputServerInterface captures the set of operations for interacting with the EaseMesh REST apis of the observability output server resource.
type ObservabilityOutputServerInterface interface {
	Get(context.Context, string) (*resource.ObservabilityOutputServer, error)
	Patch(context.Context, *resource.ObservabilityOutputServer) error
	Create(context.Context, *resource.ObservabilityOutputServer) error
	Delete(context.Context, string) error
	List(context.Context) ([]*resource.ObservabilityOutputServer, error)
}
`
