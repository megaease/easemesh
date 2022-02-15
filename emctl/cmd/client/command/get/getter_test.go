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

package get

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

func TestGetter(t *testing.T) {
	reactorType := "__test_reactor_type"
	rks := meshtesting.GetAllResourceKinds()
	rks = append(rks, meshtesting.ResourceTypeKind{
		Type: reflect.TypeOf(resource.ServiceInstance{}),
		Kind: resource.KindServiceInstance,
	})
	fake.NewResourceReactorBuilder(reactorType).
		AddReactor("*", "*", "*", func(action fake.Action) (handled bool, rets []meta.MeshObject, err error) {
			for _, rk := range rks {
				if action.GetVersionKind().Kind == rk.Kind {
					return true,
						[]meta.MeshObject{
							meshtesting.CreateMeshObjectFromType(rk.Type, rk.Kind, "new/bbb"),
						},
						nil
				}
			}
			return true, nil, nil
		}).Added()

	client := meshclient.NewFakeClient(reactorType)
	for _, rk := range rks {

		obj := meshtesting.CreateMeshObjectFromType(rk.Type, rk.Kind, "new/bbb")
		_, err := WrapGetterByMeshObject(obj, client, time.Second*1).Get()
		if err != nil {
			t.Fatalf("get resource should be successful, but %s", err)
		}
	}

	for _, rk := range rks {
		obj := meshtesting.CreateMeshObjectFromType(rk.Type, rk.Kind, "")
		_, err := WrapGetterByMeshObject(obj, client, time.Second*1).Get()
		if err != nil {
			t.Fatalf("get resource should be successful, but %s", err)
		}

	}
}

func TestGetterFail(t *testing.T) {
	reactorType := "__test_reactor_type"
	rks := meshtesting.GetAllResourceKinds()
	rks = append(rks, meshtesting.ResourceTypeKind{
		Type: reflect.TypeOf(resource.ServiceInstance{}),
		Kind: resource.KindServiceInstance,
	})
	fake.NewResourceReactorBuilder(reactorType).
		AddReactor("*", "*", "*", func(action fake.Action) (handled bool, rets []meta.MeshObject, err error) {
			return true, nil, nil
		}).Added()

	client := meshclient.NewFakeClient(reactorType)
	for _, rk := range rks {

		obj := meshtesting.CreateMeshObjectFromType(rk.Type, rk.Kind, "new/bbb")
		_, err := WrapGetterByMeshObject(obj, client, time.Second*1).Get()
		if err == nil {
			t.Fatalf("expect got an error with getting operation")
		}
	}

	for _, rk := range rks {
		obj := meshtesting.CreateMeshObjectFromType(rk.Type, rk.Kind, "")
		_, err := WrapGetterByMeshObject(obj, client, time.Second*1).Get()
		if err == nil {
			t.Fatalf("expect got an error with getting operation")
		}
	}

	WrapGetterByMeshObject(meshtesting.CreateMeshObjectFromType(reflect.TypeOf(resource.ServiceInstance{}),
		resource.KindServiceInstance, "aaaa"), client, time.Second*1).Get()
}
