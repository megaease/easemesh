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

package get

import (
	"context"
	"time"

	"github.com/megaease/easemeshctl/cmd/client/command/meshclient"
	"github.com/megaease/easemeshctl/cmd/client/resource"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
)

// WrapGetterByMeshObject wraps getter for mesh object.
func WrapGetterByMeshObject(object meta.MeshObject,
	client meshclient.MeshClient, timeout time.Duration) Getter {

	base := baseGetter{
		client:  client,
		timeout: timeout,
	}

	switch object.Kind() {
	case resource.KindMeshController:
		return &meshControllerGetter{object: object.(*resource.MeshController), baseGetter: base}
	case resource.KindService:
		return &serviceGetter{object: object.(*resource.Service), baseGetter: base}
	case resource.KindServiceInstance:
		return &serviceInstanceGetter{object: object.(*resource.ServiceInstance), baseGetter: base}
	case resource.KindCanary:
		return &canaryGetter{object: object.(*resource.Canary), baseGetter: base}
	case resource.KindLoadBalance:
		return &loadBalanceGetter{object: object.(*resource.LoadBalance), baseGetter: base}
	case resource.KindTenant:
		return &tenantGetter{object: object.(*resource.Tenant), baseGetter: base}
	case resource.KindResilience:
		return &resilienceGetter{object: object.(*resource.Resilience), baseGetter: base}
	case resource.KindObservabilityMetrics:
		return &observabilityMetricsGetter{object: object.(*resource.ObservabilityMetrics), baseGetter: base}
	case resource.KindObservabilityOutputServer:
		return &observabilityOutputServerGetter{object: object.(*resource.ObservabilityOutputServer), baseGetter: base}
	case resource.KindObservabilityTracings:
		return &observabilityTracingsGetter{object: object.(*resource.ObservabilityTracings), baseGetter: base}
	case resource.KindIngress:
		return &ingressGetter{object: object.(*resource.Ingress), baseGetter: base}
	case resource.KindCustomResourceKind:
		return &customResourceKindGetter{object: object.(*resource.CustomResourceKind), baseGetter: base}
	default:
		return &customResourceGetter{object: object.(*resource.CustomResource), baseGetter: base}
	}
}

type (
	// Getter is the getter interface fo mesh object.
	Getter interface {
		Get() ([]meta.MeshObject, error)
	}

	baseGetter struct {
		client  meshclient.MeshClient
		timeout time.Duration
	}
)

type meshControllerGetter struct {
	baseGetter
	object *resource.MeshController
}

func (s *meshControllerGetter) Get() ([]meta.MeshObject, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), s.timeout)
	defer cancelFunc()

	if s.object.Name() != "" {
		meshController, err := s.client.V1Alpha1().MeshController().Get(ctx, s.object.Name())
		if err != nil {
			return nil, err
		}

		return []meta.MeshObject{meshController}, nil
	}

	meshControllers, err := s.client.V1Alpha1().MeshController().List(ctx)
	if err != nil {
		return nil, err
	}

	objects := make([]meta.MeshObject, len(meshControllers))
	for i := range meshControllers {
		objects[i] = meshControllers[i]
	}

	return objects, nil
}

type serviceGetter struct {
	baseGetter
	object *resource.Service
}

func (s *serviceGetter) Get() ([]meta.MeshObject, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), s.timeout)
	defer cancelFunc()

	if s.object.Name() != "" {
		service, err := s.client.V1Alpha1().Service().Get(ctx, s.object.Name())
		if err != nil {
			return nil, err
		}

		return []meta.MeshObject{service}, nil
	}

	services, err := s.client.V1Alpha1().Service().List(ctx)
	if err != nil {
		return nil, err
	}

	objects := make([]meta.MeshObject, len(services))
	for i := range services {
		objects[i] = services[i]
	}

	return objects, nil
}

type serviceInstanceGetter struct {
	baseGetter
	object *resource.ServiceInstance
}

func (s *serviceInstanceGetter) Get() ([]meta.MeshObject, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), s.timeout)
	defer cancelFunc()

	if s.object.Name() != "" {
		serviceName, instanceID, err := s.object.ParseName()
		if err != nil {
			return nil, err
		}

		serviceInstance, err := s.client.V1Alpha1().ServiceInstance().Get(ctx,
			serviceName, instanceID)
		if err != nil {
			return nil, err
		}

		return []meta.MeshObject{serviceInstance}, nil
	}

	serviceInstances, err := s.client.V1Alpha1().ServiceInstance().List(ctx)
	if err != nil {
		return nil, err
	}

	objects := make([]meta.MeshObject, len(serviceInstances))
	for i := range serviceInstances {
		objects[i] = serviceInstances[i]
	}

	return objects, nil
}

type canaryGetter struct {
	baseGetter
	object *resource.Canary
}

func (c *canaryGetter) Get() ([]meta.MeshObject, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFunc()

	if c.object.Name() != "" {
		canary, err := c.client.V1Alpha1().Canary().Get(ctx, c.object.Name())
		if err != nil {
			return nil, err
		}

		return []meta.MeshObject{canary}, nil
	}

	canaries, err := c.client.V1Alpha1().Canary().List(ctx)
	if err != nil {
		return nil, err
	}

	objects := make([]meta.MeshObject, len(canaries))
	for i := range canaries {
		objects[i] = canaries[i]
	}

	return objects, nil
}

type observabilityTracingsGetter struct {
	baseGetter
	object *resource.ObservabilityTracings
}

func (o *observabilityTracingsGetter) Get() ([]meta.MeshObject, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), o.timeout)
	defer cancelFunc()

	if o.object.Name() != "" {
		tracings, err := o.client.V1Alpha1().ObservabilityTracings().Get(ctx, o.object.Name())
		if err != nil {
			return nil, err
		}

		return []meta.MeshObject{tracings}, nil
	}

	tracings, err := o.client.V1Alpha1().ObservabilityTracings().List(ctx)
	if err != nil {
		return nil, err
	}

	objects := make([]meta.MeshObject, len(tracings))
	for i := range tracings {
		objects[i] = tracings[i]
	}

	return objects, nil
}

type observabilityMetricsGetter struct {
	baseGetter
	object *resource.ObservabilityMetrics
}

func (o *observabilityMetricsGetter) Get() ([]meta.MeshObject, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), o.timeout)
	defer cancelFunc()

	if o.object.Name() != "" {
		metrics, err := o.client.V1Alpha1().ObservabilityMetrics().Get(ctx, o.object.Name())
		if err != nil {
			return nil, err
		}

		return []meta.MeshObject{metrics}, nil
	}

	metrics, err := o.client.V1Alpha1().ObservabilityMetrics().List(ctx)
	if err != nil {
		return nil, err
	}

	objects := make([]meta.MeshObject, len(metrics))
	for i := range metrics {
		objects[i] = metrics[i]
	}

	return objects, nil
}

type observabilityOutputServerGetter struct {
	baseGetter
	object *resource.ObservabilityOutputServer
}

func (o *observabilityOutputServerGetter) Get() ([]meta.MeshObject, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), o.timeout)
	defer cancelFunc()

	if o.object.Name() != "" {
		server, err := o.client.V1Alpha1().ObservabilityOutputServer().Get(ctx, o.object.Name())
		if err != nil {
			return nil, err
		}

		return []meta.MeshObject{server}, nil
	}

	servers, err := o.client.V1Alpha1().ObservabilityOutputServer().List(ctx)
	if err != nil {
		return nil, err
	}

	objects := make([]meta.MeshObject, len(servers))
	for i := range servers {
		objects[i] = servers[i]
	}

	return objects, nil
}

type loadBalanceGetter struct {
	baseGetter
	object *resource.LoadBalance
}

func (l *loadBalanceGetter) Get() ([]meta.MeshObject, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), l.timeout)
	defer cancelFunc()

	if l.object.Name() != "" {
		lb, err := l.client.V1Alpha1().LoadBalance().Get(ctx, l.object.Name())
		if err != nil {
			return nil, err
		}

		return []meta.MeshObject{lb}, nil
	}

	lbs, err := l.client.V1Alpha1().LoadBalance().List(ctx)
	if err != nil {
		return nil, err
	}

	objects := make([]meta.MeshObject, len(lbs))
	for i := range lbs {
		objects[i] = lbs[i]
	}

	return objects, nil
}

type tenantGetter struct {
	baseGetter
	object *resource.Tenant
}

func (t *tenantGetter) Get() ([]meta.MeshObject, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), t.timeout)
	defer cancelFunc()

	if t.object.Name() != "" {
		tenant, err := t.client.V1Alpha1().Tenant().Get(ctx, t.object.Name())
		if err != nil {
			return nil, err
		}

		return []meta.MeshObject{tenant}, nil
	}

	tenants, err := t.client.V1Alpha1().Tenant().List(ctx)
	if err != nil {
		return nil, err
	}

	objects := make([]meta.MeshObject, len(tenants))
	for i := range tenants {
		objects[i] = tenants[i]
	}

	return objects, nil
}

type resilienceGetter struct {
	baseGetter
	object *resource.Resilience
}

func (r *resilienceGetter) Get() ([]meta.MeshObject, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), r.timeout)
	defer cancelFunc()

	if r.object.Name() != "" {
		resilience, err := r.client.V1Alpha1().Resilience().Get(ctx, r.object.Name())
		if err != nil {
			return nil, err
		}

		return []meta.MeshObject{resilience}, nil
	}

	resiliences, err := r.client.V1Alpha1().Resilience().List(ctx)
	if err != nil {
		return nil, err
	}

	objects := make([]meta.MeshObject, len(resiliences))
	for i := range resiliences {
		objects[i] = resiliences[i]
	}

	return objects, nil
}

type ingressGetter struct {
	baseGetter
	object *resource.Ingress
}

func (i *ingressGetter) Get() ([]meta.MeshObject, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), i.timeout)
	defer cancelFunc()

	if i.object.Name() != "" {
		ingress, err := i.client.V1Alpha1().Ingress().Get(ctx, i.object.Name())
		if err != nil {
			return nil, err
		}

		return []meta.MeshObject{ingress}, nil
	}

	ingresses, err := i.client.V1Alpha1().Ingress().List(ctx)
	if err != nil {
		return nil, err
	}

	objects := make([]meta.MeshObject, len(ingresses))
	for i := range ingresses {
		objects[i] = ingresses[i]
	}

	return objects, nil
}

type customResourceKindGetter struct {
	baseGetter
	object *resource.CustomResourceKind
}

func (k *customResourceKindGetter) Get() ([]meta.MeshObject, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), k.timeout)
	defer cancelFunc()

	if k.object.Name() != "" {
		customResourceKind, err := k.client.V1Alpha1().CustomResourceKind().Get(ctx, k.object.Name())
		if err != nil {
			return nil, err
		}

		return []meta.MeshObject{customResourceKind}, nil
	}

	customResourceKinds, err := k.client.V1Alpha1().CustomResourceKind().List(ctx)
	if err != nil {
		return nil, err
	}

	objects := make([]meta.MeshObject, len(customResourceKinds))
	for i := range customResourceKinds {
		objects[i] = customResourceKinds[i]
	}

	return objects, nil
}

type customResourceGetter struct {
	baseGetter
	object *resource.CustomResource
}

func (crg *customResourceGetter) Get() ([]meta.MeshObject, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), crg.timeout)
	defer cancelFunc()

	if crg.object.Name() != "" {
		customResource, err := crg.client.V1Alpha1().CustomResource().Get(ctx, crg.object.Kind(), crg.object.Name())
		if err != nil {
			return nil, err
		}

		return []meta.MeshObject{customResource}, nil
	}

	customResources, err := crg.client.V1Alpha1().CustomResource().List(ctx, crg.object.Kind())
	if err != nil {
		return nil, err
	}

	objects := make([]meta.MeshObject, len(customResources))
	for i := range customResources {
		objects[i] = customResources[i]
	}

	return objects, nil
}
