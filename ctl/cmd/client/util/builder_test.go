package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/megaease/easemeshctl/cmd/client/resource"

	"github.com/davecgh/go-spew/spew"
	utiltesting "k8s.io/client-go/util/testing"
)

func createTestDir(t *testing.T, path string) {
	if err := os.MkdirAll(path, 0750); err != nil {
		t.Fatalf("error creating test dir: %v", err)
	}
}

func writeTestFile(t *testing.T, path string, contents string) {
	if err := ioutil.WriteFile(path, []byte(contents), 0644); err != nil {
		t.Fatalf("error creating test file %#v", err)
	}
}

const (
	aTenant = `kind: Tenant
apiVersion: mesh.megaease.com/v1alpha1
metadata:
  name: tenant_{id}
spec:
  service: []
`

	aService = `kind: Service
apiVersion: mesh.megaease.com/v1alpha1
metadata:
  name: service_{id}
spec:
   registerTenant: tenant_{id}
`
)

func TestBuilderVisitor(t *testing.T) {
	// create test dirs
	tmpDir, err := utiltesting.MkTmpdir("spec_test")
	if err != nil {
		t.Fatalf("error creating temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	createTestDir(t, fmt.Sprintf("%s/%s", tmpDir, "recursive/tenant/tenant1"))
	createTestDir(t, fmt.Sprintf("%s/%s", tmpDir, "recursive/service/service1"))

	writeTestFile(t, fmt.Sprintf("%s/recursive/tenant/tenant0.yaml", tmpDir), strings.Replace(aTenant, "{id}", "0", -1))
	writeTestFile(t, fmt.Sprintf("%s/recursive/tenant/tenant1/tenant1.yaml", tmpDir), strings.Replace(aTenant, "{id}", "1", -1))
	writeTestFile(t, fmt.Sprintf("%s/recursive/service/service0.yaml", tmpDir), strings.Replace(aService, "{id}", "0", -1))
	writeTestFile(t, fmt.Sprintf("%s/recursive/service/service1/service1.yaml", tmpDir), strings.Replace(aService, "{id}", "1", -1))

	tests := []struct {
		name          string
		meshObject    resource.MeshObject
		recursive     bool
		directory     string
		expectedNames []string
	}{
		{"recursive-service", &resource.Service{}, true, fmt.Sprintf("%s/recursive/service", tmpDir), []string{"service_0", "service_1"}},
		{"recursive-tenant", &resource.Tenant{}, true, fmt.Sprintf("%s/recursive/tenant", tmpDir), []string{"tenant_0", "tenant_1"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vs, err := NewVisitorBuilder().
				FilenameParam(&FilenameOptions{Recursive: tt.recursive, Filenames: []string{tt.directory}}).
				Do()
			if err != nil {
				t.Fatalf("build visitor error: %s", err)
			}

			if len(vs) < 1 {
				t.Fatal("number of visitors built should greater than 1 ")
			}

			var results []resource.MeshObject
			for _, v := range vs {
				v.Visit(func(mo resource.MeshObject, e error) error {
					if e != nil {
						t.Errorf("visitor error meshobject: %s", e)
					}
					results = append(results, mo)
					return nil
				})
			}

			if len(results) < 1 {
				t.Error("number of results should greater than 1")
			}

			for i, r := range results {
				switch tt.meshObject.(type) {
				case *resource.Tenant:
					if tenant, ok := r.(*resource.Tenant); !ok || tenant.Name() != tt.expectedNames[i] {
						t.Errorf("expect tenant: %s but unexpected info: %v", tt.expectedNames[i], spew.Sdump(r))
					}
				case *resource.Service:
					if service, ok := r.(*resource.Service); !ok || service.Name() != tt.expectedNames[i] {
						t.Errorf("expect service: %s but unexpected info: %v", tt.expectedNames[i], spew.Sdump(r))
					}
				}
			}

		})
	}
}
