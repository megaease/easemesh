/*
 * Copyright (c) 2017, MegaEase
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package apply

import (
	"reflect"
	"testing"
	"time"

	"github.com/megaease/easemeshctl/cmd/client/command/meshclient"
	"github.com/megaease/easemeshctl/cmd/client/command/meshclient/fake"
	"github.com/megaease/easemeshctl/cmd/client/resource"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
	"github.com/pkg/errors"
)

type resourceKind struct {
	rType reflect.Type
	kind  string
}

func getResourceKinds() []resourceKind {

	return []resourceKind{
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
		{rType: reflect.TypeOf(resource.Resilience{}), kind: resource.KindResilience},
	}
}

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

func TestApplierCreateSuccessful(t *testing.T) {
	reactorType := "__reactor"
	fake.NewResourceReactorBuilder(reactorType).
		AddReactor("*", "*", "*", func(fake.Action) (bool, []meta.MeshObject, error) {
			return true, nil, nil
		}).
		Added()

	types := getResourceKinds()
	client := meshclient.NewFakeClient(reactorType)
	for _, tp := range types {
		resource := createMeshObjectFromType(tp.rType, tp.kind, "new")
		err := WrapApplierByMeshObject(resource, client, time.Second*1).Apply()
		if err != nil {
			t.Fatalf("apply %+v, error:%s", resource, err)
		}
	}

}

func TestApplierLoopOver(t *testing.T) {
	status := map[string]error{}
	reactorType := "__reactor"
	fake.NewResourceReactorBuilder(reactorType).
		AddReactor("*", "*", "*", func(action fake.Action) (bool, []meta.MeshObject, error) {
			err1, ok := status[action.GetVersionKind().Kind]
			if !ok {
				status[action.GetVersionKind().Kind] = meshclient.ConflictError
				return true, nil, meshclient.ConflictError
			}
			switch {
			case meshclient.IsConflictError(err1):
				status[action.GetVersionKind().Kind] = meshclient.NotFoundError
				return true, nil, meshclient.NotFoundError
			case meshclient.IsNotFoundError(err1):
				return true, nil, nil
			}
			return true, nil, nil
		}).
		Added()
	types := getResourceKinds()
	client := meshclient.NewFakeClient(reactorType)
	for _, tp := range types {
		resource := createMeshObjectFromType(tp.rType, tp.kind, "new")
		err := WrapApplierByMeshObject(resource, client, time.Second*1).Apply()
		if err != nil {
			t.Fatalf("apply %+v, error:%s", resource, err)
		}
	}
}

func TestApplierFastFail(t *testing.T) {
	reactorType := "__reactor"
	fake.NewResourceReactorBuilder(reactorType).
		AddReactor("*", "*", "*", func(action fake.Action) (bool, []meta.MeshObject, error) {
			return true, nil, errors.Errorf("unknown error")
		}).
		Added()
	types := getResourceKinds()
	client := meshclient.NewFakeClient(reactorType)
	for _, tp := range types {
		resource := createMeshObjectFromType(tp.rType, tp.kind, "new")
		err := WrapApplierByMeshObject(resource, client, time.Second*1).Apply()
		if err == nil {
			t.Fatalf("apply %+v, error:%s", resource, err)
		}
	}

	serviceInstance := createMeshObjectFromType(reflect.TypeOf(resource.ServiceInstance{}), resource.KindServiceInstance, "new")
	err := WrapApplierByMeshObject(serviceInstance, client, time.Second*1).Apply()
	if err == nil {
		t.Fatal("serviceinstance applier should failure")
	}
}

func TestApplierCreateFail(t *testing.T) {
	status := map[string]error{}
	reactorType := "__reactor"
	fake.NewResourceReactorBuilder(reactorType).
		AddReactor("*", "*", "*", func(action fake.Action) (bool, []meta.MeshObject, error) {
			err1, ok := status[action.GetVersionKind().Kind]
			if !ok {
				status[action.GetVersionKind().Kind] = meshclient.ConflictError
				return true, nil, meshclient.ConflictError
			}
			switch {
			case meshclient.IsConflictError(err1):
				status[action.GetVersionKind().Kind] = meshclient.NotFoundError
				return true, nil, meshclient.NotFoundError
			case meshclient.IsNotFoundError(err1):
				return true, nil, err1
			}
			return true, nil, nil
		}).
		Added()
	types := getResourceKinds()
	client := meshclient.NewFakeClient(reactorType)
	for _, tp := range types {
		resource := createMeshObjectFromType(tp.rType, tp.kind, "new")
		err := WrapApplierByMeshObject(resource, client, time.Second*1).Apply()
		if err == nil {
			t.Fatalf("apply %+v, should raise an error", resource)
		}
	}
}

func TestApplierPatchFail(t *testing.T) {
	status := map[string]error{}
	reactorType := "__reactor"
	fake.NewResourceReactorBuilder(reactorType).
		AddReactor("*", "*", "*", func(action fake.Action) (bool, []meta.MeshObject, error) {
			err1, ok := status[action.GetVersionKind().Kind]
			if !ok {
				status[action.GetVersionKind().Kind] = meshclient.ConflictError
				return true, nil, meshclient.ConflictError
			}
			switch {
			case meshclient.IsConflictError(err1):
				return true, nil, err1
			}
			return true, nil, nil
		}).
		Added()
	types := getResourceKinds()
	client := meshclient.NewFakeClient(reactorType)
	for _, tp := range types {
		resource := createMeshObjectFromType(tp.rType, tp.kind, "new")
		err := WrapApplierByMeshObject(resource, client, time.Second*1).Apply()
		if err == nil {
			t.Fatalf("apply %+v, should raise an error", resource)
		}
	}
}
