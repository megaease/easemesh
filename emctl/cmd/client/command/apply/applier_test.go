/*
 * Copyright (c) 2021, MegaEase
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
	meshtesting "github.com/megaease/easemeshctl/cmd/client/testing"
	"github.com/pkg/errors"
)

func TestApplierCreateSuccessful(t *testing.T) {
	reactorType := "__reactor"
	fake.NewResourceReactorBuilder(reactorType).
		AddReactor("*", "*", "*", func(fake.Action) (bool, []meta.MeshObject, error) {
			return true, nil, nil
		}).
		Added()

	types := meshtesting.GetAllResourceKinds()
	client := meshclient.NewFakeClient(reactorType)
	for _, tp := range types {
		resource := meshtesting.CreateMeshObjectFromType(tp.Type, tp.Kind, "new")
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
	types := meshtesting.GetAllResourceKinds()
	client := meshclient.NewFakeClient(reactorType)
	for _, tp := range types {
		resource := meshtesting.CreateMeshObjectFromType(tp.Type, tp.Kind, "new")
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
	types := meshtesting.GetAllResourceKinds()
	client := meshclient.NewFakeClient(reactorType)
	for _, tp := range types {
		resource := meshtesting.CreateMeshObjectFromType(tp.Type, tp.Kind, "new")
		err := WrapApplierByMeshObject(resource, client, time.Second*1).Apply()
		if err == nil {
			t.Fatalf("apply %+v, error:%s", resource, err)
		}
	}

	serviceInstance := meshtesting.CreateMeshObjectFromType(reflect.TypeOf(resource.ServiceInstance{}), resource.KindServiceInstance, "new")
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
	types := meshtesting.GetAllResourceKinds()
	client := meshclient.NewFakeClient(reactorType)
	for _, tp := range types {
		resource := meshtesting.CreateMeshObjectFromType(tp.Type, tp.Kind, "new")
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
	types := meshtesting.GetAllResourceKinds()
	client := meshclient.NewFakeClient(reactorType)
	for _, tp := range types {
		resource := meshtesting.CreateMeshObjectFromType(tp.Type, tp.Kind, "new")
		err := WrapApplierByMeshObject(resource, client, time.Second*1).Apply()
		if err == nil {
			t.Fatalf("apply %+v, should raise an error", resource)
		}
	}
}
