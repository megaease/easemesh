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

package delete

import (
	"reflect"
	"testing"
	"time"

	"github.com/megaease/easemeshctl/cmd/client/command/meshclient"
	"github.com/megaease/easemeshctl/cmd/client/command/meshclient/fake"
	"github.com/megaease/easemeshctl/cmd/client/resource"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
	meshtesting "github.com/megaease/easemeshctl/cmd/client/testing"
)

func TestDeleter(t *testing.T) {
	reactorType := "__test_reactor_type"
	fake.NewResourceReactorBuilder(reactorType).
		AddReactor("*", "*", "*", func(action fake.Action) (handled bool, rets []meta.MeshObject, err error) {
			return true, nil, nil
		}).Added()

	rks := meshtesting.GetAllResourceKinds()
	rks = append(rks, meshtesting.ResourceTypeKind{Type: reflect.TypeOf(resource.ServiceInstance{}),
		Kind: resource.KindServiceInstance})
	client := meshclient.NewFakeClient(reactorType)
	for _, rk := range rks {

		obj := meshtesting.CreateMeshObjectFromType(rk.Type, rk.Kind, "new/bbb")
		err := WrapDeleterByMeshObject(obj, client, time.Second*1).Delete()
		if err != nil {
			t.Fatalf("delete resource should be successful, but %s", err)
		}
	}
}

func TestDeleterFail(t *testing.T) {
	reactorType := "__test_reactor_type"
	fake.NewResourceReactorBuilder(reactorType).
		AddReactor("*", "*", "*", func(action fake.Action) (handled bool, rets []meta.MeshObject, err error) {
			return true, nil, meshclient.NotFoundError
		}).Added()

	rks := meshtesting.GetAllResourceKinds()
	rks = append(rks, meshtesting.ResourceTypeKind{Type: reflect.TypeOf(resource.ServiceInstance{}),
		Kind: resource.KindServiceInstance})
	client := meshclient.NewFakeClient(reactorType)
	for _, rk := range rks {

		obj := meshtesting.CreateMeshObjectFromType(rk.Type, rk.Kind, "newbbb")
		err := WrapDeleterByMeshObject(obj, client, time.Second*1).Delete()
		if err == nil {
			t.Fatalf("delete resource should be successful, but %s", err)
		}
	}
}
