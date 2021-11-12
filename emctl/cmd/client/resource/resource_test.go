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

	"github.com/megaease/easemesh-api/v1alpha1"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
)

func TestObjectCreator(t *testing.T) {

	kinds := []string{KindCanary, KindCustomResourceKind, KindIngress, KindLoadBalance,
		KindMeshController, KindObservabilityMetrics, KindObservabilityOutputServer, KindObservabilityTracings,
		KindResilience, KindService, KindServiceInstance, KindTenant, "CustomResource"}

	NewObjectCreator().NewFromResource(meta.MeshResource{
		VersionKind: meta.VersionKind{
			Kind:       KindCanary,
			APIVersion: DefaultAPIVersion,
		}})

	for _, kind := range kinds {
		resource, err := NewObjectCreator().NewFromKind(meta.VersionKind{Kind: kind})
		if err != nil {
			t.Fatalf("resource should be create from kind %+v but got an error: %s", kind, err)
		}
		switch r := resource.(type) {
		case *LoadBalance:
			l := r.ToV1Alpha1()
			ToLoadBalance("new", l).Columns()
			r.Spec = &v1alpha1.LoadBalance{}
			l = r.ToV1Alpha1()
			ToLoadBalance("new", l).Columns()
		case *MeshController:
			ToMeshController(r.ToV1Alpha1()).Columns()
		case *Ingress:
			ToIngress(r.ToV1Alpha1())
			r.Spec = &IngressSpec{}
			ToIngress(r.ToV1Alpha1())
		case *Canary:
			r.Spec = &v1alpha1.Canary{}
			ToCanary("new", r.ToV1Alpha1())
		case *CustomResourceKind:
			ToCustomResourceKind(r.ToV1Alpha1())
		case *ObservabilityMetrics:
			ToObservabilityMetrics("new", r.ToV1Alpha1())
		case *ObservabilityOutputServer:
			ToObservabilityOutputServer("new", r.ToV1Alpha1())
		case *ObservabilityTracings:
			ToObservabilityTracings("new", r.ToV1Alpha1())
		case *Resilience:
			r.Spec = &v1alpha1.Resilience{}
			ToResilience("new", r.ToV1Alpha1())
		case *Service:
			ToService(r.ToV1Alpha1()).Columns()
			r.Spec = &ServiceSpec{}
			s := ToService(r.ToV1Alpha1())
			s.Spec = nil
			s.Columns()
		case *ServiceInstance:
			r.Spec = &v1alpha1.ServiceInstance{
				ServiceName: "aaa",
				InstanceID:  "bbb",
			}
			ToServiceInstance(r.ToV1Alpha1()).Columns()
			s := ToServiceInstance(r.ToV1Alpha1())
			s.Spec = nil
			s.Columns()
			r.ParseName()
			r.MetaData.Name = "aaa/bbb"
			r.ParseName()
		case *Tenant:
			r.Spec = &TenantSpec{}
			ToTenant(r.ToV1Alpha1()).Columns()
			t := ToTenant(r.ToV1Alpha1())
			t.Spec = nil
			t.Columns()
		case *CustomResource:
			ToCustomResource(map[string]interface{}{
				"name": "name",
				"kind": "kind1",
			}).ToV1Alpha1()
		}

	}
}
