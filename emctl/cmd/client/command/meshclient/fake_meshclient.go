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

package meshclient

import (
	"context"

	"github.com/megaease/easemeshctl/cmd/client/command/meshclient/fake"
	"github.com/megaease/easemeshctl/cmd/client/resource"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
	"github.com/pkg/errors"
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
	fakeServiceGetter struct {
		baseGetter
	}

	fakeObservabilityOutputServerGetter struct {
		baseGetter
	}

	fakeObservabilityMetricGetter struct {
		baseGetter
	}

	fakeObservabilityTracingGetter struct {
		baseGetter
	}

	fakeResilienceGetter struct {
		baseGetter
	}

	fakeCanaryGetter struct {
		baseGetter
	}

	fakeLoadbalanceGetter struct {
		baseGetter
	}

	fakeServiceInstanceGetter struct {
		baseGetter
	}

	fakeIngressGetter struct {
		baseGetter
	}

	fakeServiceCanaryGetter struct {
		baseGetter
	}

	fakeCustomResourceKindGetter struct {
		baseGetter
	}

	fakeCustomResourceGetter struct {
		baseGetter
	}
	fakeV1alpha1 struct {
		resourceReactor fake.ResourceReactor
	}
	fakeMeshClient struct {
		reactorType string
	}
)

var _ MeshClient = &fakeMeshClient{}

func (b *baseGetter) doModifyRequest(kind, name string, obj meta.MeshObject) error {

	_, err := b.resourceReactor.DoRequest("get", kind, name, obj)
	if err != nil {
		return err
	}
	return nil
}

func (f *fakeMeshClient) V1Alpha1() V1Alpha1Interface {
	return &fakeV1alpha1{fake.ResourceReactorForType(f.reactorType)}
}
func (f *fakeV1alpha1) Tenant() TenantInterface {
	return &fakeTenantGetter{baseGetter: baseGetter{resourceReactor: f.resourceReactor,
		kind: resource.KindTenant}}
}
func (f *fakeV1alpha1) Service() ServiceInterface {
	return &fakeServiceGetter{baseGetter: baseGetter{resourceReactor: f.resourceReactor,
		kind: resource.KindTenant}}
}
func (f *fakeV1alpha1) ServiceInstance() ServiceInstanceInterface {
	return &fakeServiceInstanceGetter{baseGetter: baseGetter{resourceReactor: f.resourceReactor,
		kind: resource.KindTenant}}
}

func (f *fakeV1alpha1) LoadBalance() LoadBalanceInterface {
	return &fakeLoadbalanceGetter{baseGetter: baseGetter{resourceReactor: f.resourceReactor,
		kind: resource.KindTenant}}
}
func (f *fakeV1alpha1) Canary() CanaryInterface {
	return &fakeCanaryGetter{baseGetter: baseGetter{resourceReactor: f.resourceReactor,
		kind: resource.KindTenant}}
}
func (f *fakeV1alpha1) ObservabilityTracings() ObservabilityTracingsInterface {
	return &fakeObservabilityTracingGetter{baseGetter: baseGetter{resourceReactor: f.resourceReactor,
		kind: resource.KindTenant}}
}

func (f *fakeV1alpha1) ObservabilityMetrics() ObservabilityMetricsInterface {
	return &fakeObservabilityMetricGetter{baseGetter: baseGetter{resourceReactor: f.resourceReactor,
		kind: resource.KindTenant}}
}
func (f *fakeV1alpha1) ObservabilityOutputServer() ObservabilityOutputServerInterface {
	return &fakeObservabilityOutputServerGetter{baseGetter: baseGetter{resourceReactor: f.resourceReactor,
		kind: resource.KindTenant}}
}

func (f *fakeV1alpha1) Resilience() ResilienceInterface {
	return &fakeResilienceGetter{baseGetter: baseGetter{resourceReactor: f.resourceReactor,
		kind: resource.KindTenant}}
}

func (f *fakeV1alpha1) Ingress() IngressInterface {
	return &fakeIngressGetter{baseGetter: baseGetter{resourceReactor: f.resourceReactor,
		kind: resource.KindTenant}}
}

func (f *fakeV1alpha1) ServiceCanary() ServiceCanaryInterface {
	return &fakeServiceCanaryGetter{
		baseGetter: baseGetter{
			resourceReactor: f.resourceReactor,
			kind:            resource.KindServiceCanary,
		},
	}
}

func (f *fakeV1alpha1) CustomResourceKind() CustomResourceKindInterface {
	return &fakeCustomResourceKindGetter{baseGetter: baseGetter{resourceReactor: f.resourceReactor,
		kind: resource.KindTenant}}
}

func (f *fakeV1alpha1) CustomResource() CustomResourceInterface {
	return &fakeCustomResourceGetter{baseGetter: baseGetter{resourceReactor: f.resourceReactor,
		kind: resource.KindTenant}}
}

func (f *fakeV1alpha1) MeshController() MeshControllerInterface {
	return &fakeMeshControllerGetter{baseGetter: baseGetter{resourceReactor: f.resourceReactor,
		kind: resource.KindMeshController}}
}

func (f *fakeMeshControllerGetter) Get(ctx context.Context, name string) (*resource.MeshController, error) {
	o, err := f.resourceReactor.DoRequest("get", resource.KindMeshController, name, nil)
	if err != nil {
		return nil, err
	}
	if len(o) == 0 {
		return nil, NotFoundError
	}
	meshController, ok := o[0].(*resource.MeshController)
	if !ok {
		return nil, errors.Errorf("get an unknown MeshObject %+v", o)
	}
	return meshController, nil
}

func (f *fakeMeshControllerGetter) Patch(ctx context.Context, r *resource.MeshController) error {
	return f.doModifyRequest(resource.KindMeshController, r.Name(), r)
}
func (f *fakeMeshControllerGetter) Create(ctx context.Context, r *resource.MeshController) error {
	return f.doModifyRequest(resource.KindMeshController, r.Name(), r)
}
func (f *fakeMeshControllerGetter) Delete(ctx context.Context, name string) error {
	return f.doModifyRequest(resource.KindMeshController, name, nil)
}
func (f *fakeMeshControllerGetter) List(context.Context) ([]*resource.MeshController, error) {
	o, err := f.resourceReactor.DoRequest("list", resource.KindMeshController, "", nil)
	if err != nil {
		return nil, err
	}
	if len(o) == 0 {
		return nil, NotFoundError
	}
	result := []*resource.MeshController{}
	for _, m := range o {
		c := m.(*resource.MeshController)
		if c != nil {
			result = append(result, c)
		}
	}
	return result, nil
}

func (f *fakeTenantGetter) Get(ctx context.Context, name string) (*resource.Tenant, error) {
	o, err := f.resourceReactor.DoRequest("get", resource.KindTenant, name, nil)
	if err != nil {
		return nil, err
	}
	if len(o) == 0 {
		return nil, NotFoundError
	}
	tenant, ok := o[0].(*resource.Tenant)
	if !ok {
		return nil, errors.Errorf("get an unknown MeshObject %+v", o)
	}
	return tenant, nil
}

func (f *fakeTenantGetter) Patch(ctx context.Context, t *resource.Tenant) error {
	return f.doModifyRequest(resource.KindTenant, t.Name(), t)
}
func (f *fakeTenantGetter) Create(ctx context.Context, t *resource.Tenant) error {
	return f.doModifyRequest(resource.KindTenant, t.Name(), t)
}
func (f *fakeTenantGetter) Delete(ctx context.Context, name string) error {
	return f.doModifyRequest(resource.KindTenant, name, nil)
}
func (f *fakeTenantGetter) List(ctx context.Context) ([]*resource.Tenant, error) {
	o, err := f.resourceReactor.DoRequest("list", resource.KindTenant, "", nil)
	if err != nil {
		return nil, err
	}
	if len(o) == 0 {
		return nil, NotFoundError
	}
	result := []*resource.Tenant{}
	for _, m := range o {
		c := m.(*resource.Tenant)
		if c != nil {
			result = append(result, c)
		}
	}
	return result, nil
}

// fakeServiceGetter implementation

func (f *fakeServiceGetter) Get(ctx context.Context, name string) (*resource.Service, error) {
	o, err := f.resourceReactor.DoRequest("get", resource.KindService, name, nil)
	if err != nil {
		return nil, err
	}
	if len(o) == 0 {
		return nil, NotFoundError
	}
	service, ok := o[0].(*resource.Service)
	if !ok {
		return nil, errors.Errorf("get an unknown MeshObject %+v", o)
	}
	return service, nil
}

func (f *fakeServiceGetter) Patch(ctx context.Context, t *resource.Service) error {
	return f.doModifyRequest(resource.KindService, t.Name(), t)
}

func (f *fakeServiceGetter) Create(ctx context.Context, t *resource.Service) error {
	return f.doModifyRequest(resource.KindService, t.Name(), t)
}

func (f *fakeServiceGetter) Delete(ctx context.Context, name string) error {
	return f.doModifyRequest(resource.KindService, name, nil)
}

func (f *fakeServiceGetter) List(ctx context.Context) ([]*resource.Service, error) {
	o, err := f.resourceReactor.DoRequest("list", resource.KindService, "", nil)
	if err != nil {
		return nil, err
	}
	if len(o) == 0 {
		return nil, NotFoundError
	}
	result := []*resource.Service{}
	for _, m := range o {
		c := m.(*resource.Service)
		if c != nil {
			result = append(result, c)
		}
	}
	return result, nil
}

// fakeServiceInstanceGetter implementation

func (f *fakeServiceInstanceGetter) Get(ctx context.Context, name, instanceID string) (*resource.ServiceInstance, error) {
	o, err := f.resourceReactor.DoRequest("get", resource.KindServiceInstance, name+"/"+instanceID, nil)
	if err != nil {
		return nil, err
	}
	if len(o) == 0 {
		return nil, NotFoundError
	}
	service, ok := o[0].(*resource.ServiceInstance)
	if !ok {
		return nil, errors.Errorf("get an unknown MeshObject %+v", o)
	}
	return service, nil
}

func (f *fakeServiceInstanceGetter) Delete(ctx context.Context, name, instanceID string) error {
	return f.doModifyRequest(resource.KindServiceInstance, name+"/"+instanceID, nil)
}

func (f *fakeServiceInstanceGetter) List(ctx context.Context) ([]*resource.ServiceInstance, error) {
	o, err := f.resourceReactor.DoRequest("list", resource.KindServiceInstance, "", nil)
	if err != nil {
		return nil, err
	}
	if len(o) == 0 {
		return nil, NotFoundError
	}
	result := []*resource.ServiceInstance{}
	for _, m := range o {
		c := m.(*resource.ServiceInstance)
		if c != nil {
			result = append(result, c)
		}
	}
	return result, nil
}

// fakeLoadbalanceGetter implementation

func (f *fakeLoadbalanceGetter) Get(ctx context.Context, name string) (*resource.LoadBalance, error) {
	o, err := f.resourceReactor.DoRequest("get", resource.KindLoadBalance, name, nil)
	if err != nil {
		return nil, err
	}
	if len(o) == 0 {
		return nil, NotFoundError
	}
	result, ok := o[0].(*resource.LoadBalance)
	if !ok {
		return nil, errors.Errorf("get an unknown MeshObject %+v", o)
	}
	return result, nil
}

func (f *fakeLoadbalanceGetter) Patch(ctx context.Context, t *resource.LoadBalance) error {
	return f.doModifyRequest(resource.KindLoadBalance, t.Name(), t)
}

func (f *fakeLoadbalanceGetter) Create(ctx context.Context, t *resource.LoadBalance) error {
	return f.doModifyRequest(resource.KindLoadBalance, t.Name(), t)
}

func (f *fakeLoadbalanceGetter) Delete(ctx context.Context, name string) error {
	return f.doModifyRequest(resource.KindLoadBalance, name, nil)
}

func (f *fakeLoadbalanceGetter) List(ctx context.Context) ([]*resource.LoadBalance, error) {
	o, err := f.resourceReactor.DoRequest("list", resource.KindLoadBalance, "", nil)
	if err != nil {
		return nil, err
	}
	if len(o) == 0 {
		return nil, NotFoundError
	}
	result := []*resource.LoadBalance{}
	for _, m := range o {
		c := m.(*resource.LoadBalance)
		if c != nil {
			result = append(result, c)
		}
	}
	return result, nil
}

// fakeCanaryGetter implementation

func (f *fakeCanaryGetter) Get(ctx context.Context, name string) (*resource.Canary, error) {
	o, err := f.resourceReactor.DoRequest("get", resource.KindCanary, name, nil)
	if err != nil {
		return nil, err
	}
	if len(o) == 0 {
		return nil, NotFoundError
	}
	result, ok := o[0].(*resource.Canary)
	if !ok {
		return nil, errors.Errorf("get an unknown MeshObject %+v", o)
	}
	return result, nil
}

func (f *fakeCanaryGetter) Patch(ctx context.Context, t *resource.Canary) error {
	return f.doModifyRequest(resource.KindCanary, t.Name(), t)
}

func (f *fakeCanaryGetter) Create(ctx context.Context, t *resource.Canary) error {
	return f.doModifyRequest(resource.KindCanary, t.Name(), t)
}

func (f *fakeCanaryGetter) Delete(ctx context.Context, name string) error {
	return f.doModifyRequest(resource.KindCanary, name, nil)
}

func (f *fakeCanaryGetter) List(ctx context.Context) ([]*resource.Canary, error) {
	o, err := f.resourceReactor.DoRequest("list", resource.KindCanary, "", nil)
	if err != nil {
		return nil, err
	}
	if len(o) == 0 {
		return nil, NotFoundError
	}
	result := []*resource.Canary{}
	for _, m := range o {
		c := m.(*resource.Canary)
		if c != nil {
			result = append(result, c)
		}
	}
	return result, nil
}

// fakeObservabilityTracingGetter implementation

func (f *fakeObservabilityTracingGetter) Get(ctx context.Context, name string) (*resource.ObservabilityTracings, error) {
	o, err := f.resourceReactor.DoRequest("get", resource.KindObservabilityTracings, name, nil)
	if err != nil {
		return nil, err
	}
	if len(o) == 0 {
		return nil, NotFoundError
	}
	result, ok := o[0].(*resource.ObservabilityTracings)
	if !ok {
		return nil, errors.Errorf("get an unknown MeshObject %+v", o)
	}
	return result, nil
}

func (f *fakeObservabilityTracingGetter) Patch(ctx context.Context, t *resource.ObservabilityTracings) error {
	return f.doModifyRequest(resource.KindObservabilityTracings, t.Name(), t)
}

func (f *fakeObservabilityTracingGetter) Create(ctx context.Context, t *resource.ObservabilityTracings) error {
	return f.doModifyRequest(resource.KindObservabilityTracings, t.Name(), t)
}

func (f *fakeObservabilityTracingGetter) Delete(ctx context.Context, name string) error {
	return f.doModifyRequest(resource.KindObservabilityTracings, name, nil)
}

func (f *fakeObservabilityTracingGetter) List(ctx context.Context) ([]*resource.ObservabilityTracings, error) {
	o, err := f.resourceReactor.DoRequest("list", resource.KindObservabilityTracings, "", nil)
	if err != nil {
		return nil, err
	}
	if len(o) == 0 {
		return nil, NotFoundError
	}
	result := []*resource.ObservabilityTracings{}
	for _, m := range o {
		c := m.(*resource.ObservabilityTracings)
		if c != nil {
			result = append(result, c)
		}
	}
	return result, nil
}

// fakeObservabilityMetricGetter implementation

func (f *fakeObservabilityMetricGetter) Get(ctx context.Context, name string) (*resource.ObservabilityMetrics, error) {
	o, err := f.resourceReactor.DoRequest("get", resource.KindObservabilityMetrics, name, nil)
	if err != nil {
		return nil, err
	}
	if len(o) == 0 {
		return nil, NotFoundError
	}
	result, ok := o[0].(*resource.ObservabilityMetrics)
	if !ok {
		return nil, errors.Errorf("get an unknown MeshObject %+v", o)
	}
	return result, nil
}

func (f *fakeObservabilityMetricGetter) Patch(ctx context.Context, t *resource.ObservabilityMetrics) error {
	return f.doModifyRequest(resource.KindObservabilityMetrics, t.Name(), t)
}

func (f *fakeObservabilityMetricGetter) Create(ctx context.Context, t *resource.ObservabilityMetrics) error {
	return f.doModifyRequest(resource.KindObservabilityMetrics, t.Name(), t)
}

func (f *fakeObservabilityMetricGetter) Delete(ctx context.Context, name string) error {
	return f.doModifyRequest(resource.KindObservabilityMetrics, name, nil)
}

func (f *fakeObservabilityMetricGetter) List(ctx context.Context) ([]*resource.ObservabilityMetrics, error) {
	o, err := f.resourceReactor.DoRequest("list", resource.KindObservabilityMetrics, "", nil)
	if err != nil {
		return nil, err
	}
	if len(o) == 0 {
		return nil, NotFoundError
	}
	result := []*resource.ObservabilityMetrics{}
	for _, m := range o {
		c := m.(*resource.ObservabilityMetrics)
		if c != nil {
			result = append(result, c)
		}
	}
	return result, nil
}

// fakeObservabilityOutputServerGetter implementation

func (f *fakeObservabilityOutputServerGetter) Get(ctx context.Context, name string) (*resource.ObservabilityOutputServer, error) {
	o, err := f.resourceReactor.DoRequest("get", resource.KindObservabilityOutputServer, name, nil)
	if err != nil {
		return nil, err
	}
	if len(o) == 0 {
		return nil, NotFoundError
	}
	result, ok := o[0].(*resource.ObservabilityOutputServer)
	if !ok {
		return nil, errors.Errorf("get an unknown MeshObject %+v", o)
	}
	return result, nil
}

func (f *fakeObservabilityOutputServerGetter) Patch(ctx context.Context, t *resource.ObservabilityOutputServer) error {
	return f.doModifyRequest(resource.KindObservabilityOutputServer, t.Name(), t)
}

func (f *fakeObservabilityOutputServerGetter) Create(ctx context.Context, t *resource.ObservabilityOutputServer) error {
	return f.doModifyRequest(resource.KindObservabilityOutputServer, t.Name(), t)
}

func (f *fakeObservabilityOutputServerGetter) Delete(ctx context.Context, name string) error {
	return f.doModifyRequest(resource.KindObservabilityOutputServer, name, nil)
}

func (f *fakeObservabilityOutputServerGetter) List(ctx context.Context) ([]*resource.ObservabilityOutputServer, error) {
	o, err := f.resourceReactor.DoRequest("list", resource.KindObservabilityOutputServer, "", nil)
	if err != nil {
		return nil, err
	}
	if len(o) == 0 {
		return nil, NotFoundError
	}
	result := []*resource.ObservabilityOutputServer{}
	for _, m := range o {
		c := m.(*resource.ObservabilityOutputServer)
		if c != nil {
			result = append(result, c)
		}
	}
	return result, nil
}

// fakeResilienceGetter implementation

func (f *fakeResilienceGetter) Get(ctx context.Context, name string) (*resource.Resilience, error) {
	o, err := f.resourceReactor.DoRequest("get", resource.KindResilience, name, nil)
	if err != nil {
		return nil, err
	}
	if len(o) == 0 {
		return nil, NotFoundError
	}
	result, ok := o[0].(*resource.Resilience)
	if !ok {
		return nil, errors.Errorf("get an unknown MeshObject %+v", o)
	}
	return result, nil
}

func (f *fakeResilienceGetter) Patch(ctx context.Context, t *resource.Resilience) error {
	return f.doModifyRequest(resource.KindResilience, t.Name(), t)
}

func (f *fakeResilienceGetter) Create(ctx context.Context, t *resource.Resilience) error {
	return f.doModifyRequest(resource.KindResilience, t.Name(), t)
}

func (f *fakeResilienceGetter) Delete(ctx context.Context, name string) error {
	return f.doModifyRequest(resource.KindResilience, name, nil)
}

func (f *fakeResilienceGetter) List(ctx context.Context) ([]*resource.Resilience, error) {
	o, err := f.resourceReactor.DoRequest("list", resource.KindResilience, "", nil)
	if err != nil {
		return nil, err
	}
	if len(o) == 0 {
		return nil, NotFoundError
	}
	result := []*resource.Resilience{}
	for _, m := range o {
		c := m.(*resource.Resilience)
		if c != nil {
			result = append(result, c)
		}
	}
	return result, nil
}

// fakeIngressGetter implementation

func (f *fakeIngressGetter) Get(ctx context.Context, name string) (*resource.Ingress, error) {
	o, err := f.resourceReactor.DoRequest("get", resource.KindIngress, name, nil)
	if err != nil {
		return nil, err
	}
	if len(o) == 0 {
		return nil, NotFoundError
	}
	result, ok := o[0].(*resource.Ingress)
	if !ok {
		return nil, errors.Errorf("get an unknown MeshObject %+v", o)
	}
	return result, nil
}

func (f *fakeIngressGetter) Patch(ctx context.Context, t *resource.Ingress) error {
	return f.doModifyRequest(resource.KindIngress, t.Name(), t)
}

func (f *fakeIngressGetter) Create(ctx context.Context, t *resource.Ingress) error {
	return f.doModifyRequest(resource.KindIngress, t.Name(), t)
}

func (f *fakeIngressGetter) Delete(ctx context.Context, name string) error {
	return f.doModifyRequest(resource.KindIngress, name, nil)
}

func (f *fakeIngressGetter) List(ctx context.Context) ([]*resource.Ingress, error) {
	o, err := f.resourceReactor.DoRequest("list", resource.KindIngress, "", nil)
	if err != nil {
		return nil, err
	}
	if len(o) == 0 {
		return nil, NotFoundError
	}
	result := []*resource.Ingress{}
	for _, m := range o {
		c := m.(*resource.Ingress)
		if c != nil {
			result = append(result, c)
		}
	}
	return result, nil
}

func (f *fakeServiceCanaryGetter) Get(ctx context.Context, name string) (*resource.ServiceCanary, error) {
	o, err := f.resourceReactor.DoRequest("get", resource.KindServiceCanary, name, nil)
	if err != nil {
		return nil, err
	}
	if len(o) == 0 {
		return nil, NotFoundError
	}
	result, ok := o[0].(*resource.ServiceCanary)
	if !ok {
		return nil, errors.Errorf("get an unknown MeshObject %+v", o)
	}
	return result, nil
}

func (f *fakeServiceCanaryGetter) Patch(ctx context.Context, t *resource.ServiceCanary) error {
	return f.doModifyRequest(resource.KindServiceCanary, t.Name(), t)
}

func (f *fakeServiceCanaryGetter) Create(ctx context.Context, t *resource.ServiceCanary) error {
	return f.doModifyRequest(resource.KindServiceCanary, t.Name(), t)
}

func (f *fakeServiceCanaryGetter) Delete(ctx context.Context, name string) error {
	return f.doModifyRequest(resource.KindServiceCanary, name, nil)
}

func (f *fakeServiceCanaryGetter) List(ctx context.Context) ([]*resource.ServiceCanary, error) {
	o, err := f.resourceReactor.DoRequest("list", resource.KindServiceCanary, "", nil)
	if err != nil {
		return nil, err
	}
	if len(o) == 0 {
		return nil, NotFoundError
	}
	result := []*resource.ServiceCanary{}
	for _, m := range o {
		c := m.(*resource.ServiceCanary)
		if c != nil {
			result = append(result, c)
		}
	}
	return result, nil
}

// fakeCustomResourceKindGetter implementation

func (f *fakeCustomResourceKindGetter) Get(ctx context.Context, name string) (*resource.CustomResourceKind, error) {
	o, err := f.resourceReactor.DoRequest("get", resource.KindCustomResourceKind, name, nil)
	if err != nil {
		return nil, err
	}
	if len(o) == 0 {
		return nil, NotFoundError
	}
	result, ok := o[0].(*resource.CustomResourceKind)
	if !ok {
		return nil, errors.Errorf("get an unknown MeshObject %+v", o)
	}
	return result, nil
}

func (f *fakeCustomResourceKindGetter) Patch(ctx context.Context, t *resource.CustomResourceKind) error {
	return f.doModifyRequest(resource.KindCustomResourceKind, t.Name(), t)
}

func (f *fakeCustomResourceKindGetter) Create(ctx context.Context, t *resource.CustomResourceKind) error {
	return f.doModifyRequest(resource.KindCustomResourceKind, t.Name(), t)
}

func (f *fakeCustomResourceKindGetter) Delete(ctx context.Context, name string) error {
	return f.doModifyRequest(resource.KindCustomResourceKind, name, nil)
}

func (f *fakeCustomResourceKindGetter) List(ctx context.Context) ([]*resource.CustomResourceKind, error) {
	o, err := f.resourceReactor.DoRequest("list", resource.KindCustomResourceKind, "", nil)
	if err != nil {
		return nil, err
	}
	if len(o) == 0 {
		return nil, NotFoundError
	}
	result := []*resource.CustomResourceKind{}
	for _, m := range o {
		c := m.(*resource.CustomResourceKind)
		if c != nil {
			result = append(result, c)
		}
	}
	return result, nil
}

// fakeCustomResourceGetter implementation

func (f *fakeCustomResourceGetter) Get(ctx context.Context, kind, name string) (*resource.CustomResource, error) {
	// TODO: how to handle customresource ?
	o, err := f.resourceReactor.DoRequest("get", kind, name, nil)
	if err != nil {
		return nil, err
	}
	if len(o) == 0 {
		return nil, NotFoundError
	}
	result, ok := o[0].(*resource.CustomResource)
	if !ok {
		return nil, errors.Errorf("get an unknown MeshObject %+v", o)
	}
	return result, nil
}

func (f *fakeCustomResourceGetter) Patch(ctx context.Context, t *resource.CustomResource) error {
	return f.doModifyRequest(t.Kind(), t.Name(), t)
}

func (f *fakeCustomResourceGetter) Create(ctx context.Context, t *resource.CustomResource) error {
	return f.doModifyRequest(t.Kind(), t.Name(), t)
}

func (f *fakeCustomResourceGetter) Delete(ctx context.Context, kind, name string) error {
	return f.doModifyRequest(kind, name, nil)
}

func (f *fakeCustomResourceGetter) List(ctx context.Context, kind string) ([]*resource.CustomResource, error) {
	o, err := f.resourceReactor.DoRequest("list", kind, "", nil)
	if err != nil {
		return nil, err
	}
	if len(o) == 0 {
		return nil, NotFoundError
	}
	result := []*resource.CustomResource{}
	for _, m := range o {
		c := m.(*resource.CustomResource)
		if c != nil {
			result = append(result, c)
		}
	}
	return result, nil
}

// NewFakeClient return a fake meshclient
func NewFakeClient(t string) MeshClient {
	return &fakeMeshClient{reactorType: t}
}
