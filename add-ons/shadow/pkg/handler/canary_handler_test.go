package handler

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/megaease/easemeshctl/cmd/client/resource"
	k8Yaml "k8s.io/apimachinery/pkg/util/yaml"
)

func fakeServiceCanary() *resource.ServiceCanary {
	decoder := k8Yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(fakeCanaryYaml)), 1000)
	serviceCanary := &resource.ServiceCanary{}
	_ = decoder.Decode(serviceCanary)
	return serviceCanary
}

func Test_createShadowServiceCanary(t *testing.T) {
	shadowService := fakeShadowService()
	fakeCanary := fakeServiceCanary()
	newCanary := createShadowServiceCanary(&shadowService)

	if !reflect.DeepEqual(newCanary, fakeCanary) {
		t.Errorf("createShadowServiceCanary() = %v, want %v", newCanary, fakeCanary)
	}
}

const fakeCanaryYaml = `
apiVersion: mesh.megaease.com/v1alpha1
kind: ServiceCanary
metadata:
  name: shadow-visits-service
spec:
  priority: 5 
  selector:
    matchServices: [visits-service]
    matchInstanceLabels: {version: shadow}
  trafficRules:
    headers:
      X-Mesh-Canary:
        exact: shadow
`
