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
package resource

import (
	"testing"

	"github.com/megaease/easemesh-api/v2alpha1"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
	"gopkg.in/yaml.v2"
)

func TestObjectCreator(t *testing.T) {
	kinds := []string{
		KindCustomResourceKind, KindIngress, KindLoadBalance,
		KindMeshController, KindObservabilityMetrics, KindObservabilityOutputServer, KindObservabilityTracings,
		KindResilience, KindService, KindServiceInstance, KindTenant, "CustomResource",
	}

	NewObjectCreator().NewFromResource(meta.MeshResource{
		VersionKind: meta.VersionKind{
			Kind:       KindService,
			APIVersion: DefaultAPIVersion,
		},
	})

	for _, kind := range kinds {
		resource, err := NewObjectCreator().NewFromKind(meta.VersionKind{Kind: kind})
		if err != nil {
			t.Fatalf("resource should be create from kind %+v but got an error: %s", kind, err)
		}
		switch r := resource.(type) {
		case *LoadBalance:
			l := r.ToV2Alpha1()
			ToLoadBalance("new", l).Columns()
			r.Spec = &v2alpha1.LoadBalance{}
			l = r.ToV2Alpha1()
			ToLoadBalance("new", l).Columns()
		case *MeshController:
			ToMeshController(r.ToV2Alpha1()).Columns()
		case *Ingress:
			r.Spec = &IngressSpec{}
			ToIngress(r.ToV2Alpha1())
		case *HTTPRouteGroup:
			r.Spec = &HTTPRouteGroupSpec{}
			ToHTTPRouteGroup(r.ToV2Alpha1())
		case *TrafficTarget:
			r.Spec = &TrafficTargetSpec{}
			ToTrafficTarget(r.ToV2Alpha1())
		case *CustomResourceKind:
			ToCustomResourceKind(r.ToV2Alpha1())
		case *ObservabilityMetrics:
			ToObservabilityMetrics("new", r.ToV2Alpha1())
		case *ObservabilityOutputServer:
			ToObservabilityOutputServer("new", r.ToV2Alpha1())
		case *ObservabilityTracings:
			ToObservabilityTracings("new", r.ToV2Alpha1())
		case *Resilience:
			r.Spec = &v2alpha1.Resilience{}
			ToResilience("new", r.ToV2Alpha1())
		case *Mock:
			r.Spec = &v2alpha1.Mock{}
			ToMock("new", r.ToV2Alpha1())
		case *Service:
			ToService(r.ToV2Alpha1()).Columns()
			r.Spec = &ServiceSpec{}
			s := ToService(r.ToV2Alpha1())
			s.Spec = nil
			s.Columns()
		case *ServiceInstance:
			r.Spec = &v2alpha1.ServiceInstance{
				ServiceName: "aaa",
				InstanceID:  "bbb",
			}
			ToServiceInstance(r.ToV2Alpha1()).Columns()
			s := ToServiceInstance(r.ToV2Alpha1())
			s.Spec = nil
			s.Columns()
			r.ParseName()
			r.MetaData.Name = "aaa/bbb"
			r.ParseName()
		case *Tenant:
			r.Spec = &TenantSpec{}
			ToTenant(r.ToV2Alpha1()).Columns()
			t := ToTenant(r.ToV2Alpha1())
			t.Spec = nil
			t.Columns()
		case *CustomResource:
			ToCustomResource(map[string]interface{}{
				"name": "name",
				"kind": "kind1",
			}).ToV2Alpha1()
		}

	}
}

func TestDynamicObject(t *testing.T) {
	r := DynamicObject{}
	r["field1"] = map[string]interface{}{
		"sub1": 1,
		"sub2": "value2",
	}
	r["field2"] = []interface{}{
		"sub1", "sub2",
	}

	data, err := yaml.Marshal(r)
	if err != nil {
		t.Errorf("yaml.Marshal should succeed: %v", err.Error())
	}

	err = yaml.Unmarshal(data, &r)
	if err != nil {
		t.Errorf("yaml.Marshal should succeed: %v", err.Error())
	}

	if _, ok := r["field1"].(map[string]interface{}); !ok {
		t.Errorf("the type of 'field1' should be 'map[string]interface{}'")
	}
}
