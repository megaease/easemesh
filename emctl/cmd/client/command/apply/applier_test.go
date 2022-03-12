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
	"github.com/stretchr/testify/assert"
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
		meshResource := meshtesting.CreateMeshObjectFromType(tp.Type, tp.Kind, "new")
		err := WrapApplierByMeshObject(meshResource, client, time.Second*1).Apply()
		assert.ErrorIs(t, err, nil)
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
		meshResource := meshtesting.CreateMeshObjectFromType(tp.Type, tp.Kind, "new")
		err := WrapApplierByMeshObject(meshResource, client, time.Second*1).Apply()
		assert.NotErrorIs(t, err, meshclient.ConflictError)
		assert.NotErrorIs(t, err, meshclient.NotFoundError)
		assert.NotErrorIs(t, err, nil)
	}

	serviceInstance := meshtesting.CreateMeshObjectFromType(reflect.TypeOf(resource.ServiceInstance{}), resource.KindServiceInstance, "new")
	err := WrapApplierByMeshObject(serviceInstance, client, time.Second*1).Apply()
	assert.NotErrorIs(t, err, meshclient.ConflictError)
	assert.NotErrorIs(t, err, meshclient.NotFoundError)
	assert.NotErrorIs(t, err, nil)
}

func TestApplierCreateFailButPatchSuccess(t *testing.T) {
	status := map[string]error{}
	reactorType := "__reactor"
	fake.NewResourceReactorBuilder(reactorType).
		AddReactor("*", "*", "*", func(action fake.Action) (bool, []meta.MeshObject, error) {
			err1, ok := status[action.GetVersionKind().Kind]
			if !ok {
				status[action.GetVersionKind().Kind] = meshclient.ConflictError
				return true, nil, meshclient.ConflictError
			}
			assert.ErrorIs(t, err1, meshclient.ConflictError)
			return true, nil, nil
		}).
		Added()
	types := meshtesting.GetAllResourceKinds()
	client := meshclient.NewFakeClient(reactorType)
	for _, tp := range types {
		meshResource := meshtesting.CreateMeshObjectFromType(tp.Type, tp.Kind, "new")
		err := WrapApplierByMeshObject(meshResource, client, time.Second*1).Apply()
		assert.ErrorIs(t, err, nil)
	}
}

func TestApplierCreateFailAndPatchFail(t *testing.T) {
	status := map[string]error{}
	reactorType := "__reactor"
	fake.NewResourceReactorBuilder(reactorType).
		AddReactor("*", "*", "*", func(action fake.Action) (bool, []meta.MeshObject, error) {
			err1, ok := status[action.GetVersionKind().Kind]
			if !ok {
				status[action.GetVersionKind().Kind] = meshclient.ConflictError
				return true, nil, meshclient.ConflictError
			}
			if meshclient.IsConflictError(err1) {
				return true, nil, meshclient.NotFoundError
			}
			return true, nil, nil
		}).
		Added()
	types := meshtesting.GetAllResourceKinds()
	client := meshclient.NewFakeClient(reactorType)
	for _, tp := range types {
		meshResource := meshtesting.CreateMeshObjectFromType(tp.Type, tp.Kind, "new")
		err := WrapApplierByMeshObject(meshResource, client, time.Second*1).Apply()
		assert.ErrorIs(t, err, meshclient.NotFoundError)
	}
}
