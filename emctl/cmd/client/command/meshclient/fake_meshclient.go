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
package meshclient

import (
	"context"

	"github.com/megaease/easemeshctl/cmd/client/command/meshclient/fake"
	"github.com/megaease/easemeshctl/cmd/client/resource"
)

type (
	baseGetter struct {
		resourceReactor fake.ResourceReactor
		kind            string
	}

	fakeMeshControllerGetter struct {
		baseGetter
	}

	fakeTenantGetter struct {
		baseGetter
	}

	fakeMeshClient struct {
		reactorType string
	}

	fakeV1alpha1 struct {
		resourceReactor fake.ResourceReactor
	}
)

var _ MeshClient = &fakeMeshClient{}

func (f *fakeMeshClient) V1Alpha1() V1Alpha1Interface {
	return &fakeV1alpha1{fake.ResourceReactorForType(f.reactorType)}
}
func (f *fakeV1alpha1) Tenant() TenantInterface {
	return &fakeTenantGetter{baseGetter: baseGetter{resourceReactor: f.resourceReactor,
		kind: resource.KindTenant}}
}
func (f *fakeV1alpha1) Service() ServiceInterface                            { return nil }
func (f *fakeV1alpha1) ServiceInstance() ServiceInstanceInterface            { return nil }
func (f *fakeV1alpha1) LoadBalance() LoadBalanceInterface                    { return nil }
func (f *fakeV1alpha1) Canary() CanaryInterface                              { return nil }
func (f *fakeV1alpha1) ObservabilityTracings() ObservabilityTracingInterface { return nil }
func (f *fakeV1alpha1) ObservabilityMetrics() ObservabilityMetricInterface   { return nil }
func (f *fakeV1alpha1) ObservabilityOutputServer() ObservabilityOutputServerInterface {
	return nil
}
func (f *fakeV1alpha1) Resilience() ResilienceInterface                 { return nil }
func (f *fakeV1alpha1) Ingress() IngressInterface                       { return nil }
func (f *fakeV1alpha1) CustomResourceKind() CustomResourceKindInterface { return nil }
func (f *fakeV1alpha1) CustomResource() CustomResourceInterface         { return nil }

func (f *fakeV1alpha1) MeshController() MeshControllerInterface {
	return &fakeMeshControllerGetter{baseGetter: baseGetter{resourceReactor: f.resourceReactor,
		kind: resource.KindMeshController}}
}

func (f *fakeMeshControllerGetter) Get(context.Context, string) (*resource.MeshController, error) {
	return nil, nil
}
func (f *fakeMeshControllerGetter) Patch(context.Context, *resource.MeshController) error { return nil }
func (f *fakeMeshControllerGetter) Create(context.Context, *resource.MeshController) error {
	return nil
}
func (f *fakeMeshControllerGetter) Delete(context.Context, string) error { return nil }
func (f *fakeMeshControllerGetter) List(context.Context) ([]*resource.MeshController, error) {
	return nil, nil
}

func (f *fakeTenantGetter) Get(context.Context, string) (*resource.Tenant, error) { return nil, nil }
func (f *fakeTenantGetter) Patch(context.Context, *resource.Tenant) error         { return nil }
func (f *fakeTenantGetter) Create(context.Context, *resource.Tenant) error        { return nil }
func (f *fakeTenantGetter) Delete(context.Context, string) error                  { return nil }
func (f *fakeTenantGetter) List(context.Context) ([]*resource.Tenant, error)      { return nil, nil }
