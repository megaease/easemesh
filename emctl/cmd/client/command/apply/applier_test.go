package apply

import (
	"reflect"
	"testing"
	"time"

	"github.com/megaease/easemeshctl/cmd/client/command/meshclient"
	"github.com/megaease/easemeshctl/cmd/client/command/meshclient/fake"
	"github.com/megaease/easemeshctl/cmd/client/resource"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
)

func createMeshObjectFromType(t reflect.Type, kind, nm string) meta.MeshObject {
	meshObject := reflect.New(t).
		Elem() // reflect.Value

	versionKind := meshObject.FieldByName("VersionKind")
	version := versionKind.FieldByName("APIVersion")
	knd := versionKind.FieldByName("Kind")

	knd.SetString(kind)
	version.SetString("v1alpha1")

	metaData := meshObject.FieldByName("MetaData")
	name := metaData.FieldByName("Name")
	name.SetString(nm)
	return meshObject.Addr().Interface().(meta.MeshObject)
}

func TestApplierCreate(t *testing.T) {
	reactorType := "__reactor"
	fake.NewResourceReactorBuilder(reactorType).
		AddReactor("*", "*", "*", func(fake.Action) (bool, []meta.MeshObject, error) {
			return true, nil, nil
		}).
		Added()

	types := []struct {
		rType reflect.Type
		kind  string
	}{
		{rType: reflect.TypeOf(resource.Tenant{}), kind: resource.KindTenant},
		{rType: reflect.TypeOf(resource.MeshController{}), kind: resource.KindMeshController},
		{rType: reflect.TypeOf(resource.Ingress{}), kind: resource.KindIngress},
		{rType: reflect.TypeOf(resource.CustomResourceKind{}), kind: resource.KindCustomResourceKind},
		{rType: reflect.TypeOf(resource.CustomResource{}), kind: "-"},
		{rType: reflect.TypeOf(resource.LoadBalance{}), kind: resource.KindLoadBalance},
		{rType: reflect.TypeOf(resource.ObservabilityMetrics{}), kind: resource.KindObservabilityMetrics},
		{rType: reflect.TypeOf(resource.ObservabilityOutputServer{}), kind: resource.KindObservabilityOutputServer},
		{rType: reflect.TypeOf(resource.ObservabilityTracings{}), kind: resource.KindObservabilityTracings},
		{rType: reflect.TypeOf(resource.Canary{}), kind: resource.KindCanary},
		{rType: reflect.TypeOf(resource.Service{}), kind: resource.KindService},
		{rType: reflect.TypeOf(resource.ServiceInstance{}), kind: resource.KindServiceInstance},
		{rType: reflect.TypeOf(resource.Resilience{}), kind: resource.KindResilience},
	}

	client := meshclient.NewFakeClient(reactorType)
	for _, t := range types {
		resource := createMeshObjectFromType(t.rType, t.kind, "new")
		WrapApplierByMeshObject(resource, client, time.Second*1).Apply()
	}

}
