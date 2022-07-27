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
	"context"
	"time"

	"github.com/megaease/easemeshctl/cmd/client/command/meshclient"
	"github.com/megaease/easemeshctl/cmd/client/resource"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"

	"github.com/pkg/errors"
)

// WrapDeleterByMeshObject returns a new Deleter from a MeshObject
func WrapDeleterByMeshObject(object meta.MeshObject,
	client meshclient.MeshClient, timeout time.Duration,
) Deleter {
	switch object.Kind() {
	case resource.KindMeshController:
		return &meshControllerDeleter{object: object.(*resource.MeshController), baseDeleter: baseDeleter{client: client, timeout: timeout}}
	case resource.KindService:
		return &serviceDeleter{object: object.(*resource.Service), baseDeleter: baseDeleter{client: client, timeout: timeout}}
	case resource.KindServiceInstance:
		return &serviceInstanceDeleter{object: object.(*resource.ServiceInstance), baseDeleter: baseDeleter{client: client, timeout: timeout}}
	case resource.KindLoadBalance:
		return &loadBalanceDeleter{object: object.(*resource.LoadBalance), baseDeleter: baseDeleter{client: client, timeout: timeout}}
	case resource.KindTenant:
		return &tenantDeleter{object: object.(*resource.Tenant), baseDeleter: baseDeleter{client: client, timeout: timeout}}
	case resource.KindResilience:
		return &resilienceDeleter{object: object.(*resource.Resilience), baseDeleter: baseDeleter{client: client, timeout: timeout}}
	case resource.KindMock:
		return &mockDeleter{object: object.(*resource.Mock), baseDeleter: baseDeleter{client: client, timeout: timeout}}
	case resource.KindObservabilityMetrics:
		return &observabilityMetricsDeleter{object: object.(*resource.ObservabilityMetrics), baseDeleter: baseDeleter{client: client, timeout: timeout}}
	case resource.KindObservabilityOutputServer:
		return &observabilityOutputServerDeleter{object: object.(*resource.ObservabilityOutputServer), baseDeleter: baseDeleter{client: client, timeout: timeout}}
	case resource.KindObservabilityTracings:
		return &observabilityTracingsDeleter{object: object.(*resource.ObservabilityTracings), baseDeleter: baseDeleter{client: client, timeout: timeout}}
	case resource.KindIngress:
		return &ingressDeleter{object: object.(*resource.Ingress), baseDeleter: baseDeleter{client: client, timeout: timeout}}
	case resource.KindHTTPRouteGroup:
		return &httpRouteGroupDeleter{object: object.(*resource.HTTPRouteGroup), baseDeleter: baseDeleter{client: client, timeout: timeout}}
	case resource.KindTrafficTarget:
		return &trafficTargetDeleter{object: object.(*resource.TrafficTarget), baseDeleter: baseDeleter{client: client, timeout: timeout}}
	case resource.KindServiceCanary:
		return &serviceCanaryDeleter{object: object.(*resource.ServiceCanary), baseDeleter: baseDeleter{client: client, timeout: timeout}}
	case resource.KindCustomResourceKind:
		return &customResourceKindDeleter{object: object.(*resource.CustomResourceKind), baseDeleter: baseDeleter{client: client, timeout: timeout}}
	default:
		return &customResourceDeleter{object: object.(*resource.CustomResource), baseDeleter: baseDeleter{client: client, timeout: timeout}}
	}
}

// Deleter deletes configuration from the control plane service of the EaseMesh
type Deleter interface {
	Delete() error
}

type baseDeleter struct {
	client  meshclient.MeshClient
	timeout time.Duration
}

type meshControllerDeleter struct {
	baseDeleter
	object *resource.MeshController
}

func (s *meshControllerDeleter) Delete() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), s.timeout)
	defer cancelFunc()
	return s.client.V2Alpha1().MeshController().Delete(ctx, s.object.Name())
}

type serviceDeleter struct {
	baseDeleter
	object *resource.Service
}

func (s *serviceDeleter) Delete() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), s.timeout)
	defer cancelFunc()
	return s.client.V2Alpha1().Service().Delete(ctx, s.object.Name())
}

type serviceInstanceDeleter struct {
	baseDeleter
	object *resource.ServiceInstance
}

func (s *serviceInstanceDeleter) Delete() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), s.timeout)
	defer cancelFunc()

	serviceName, instanceID, err := s.object.ParseName()
	if err != nil {
		return err
	}

	return s.client.V2Alpha1().ServiceInstance().Delete(ctx,
		serviceName, instanceID)
}

type observabilityTracingsDeleter struct {
	baseDeleter
	object *resource.ObservabilityTracings
}

func (o *observabilityTracingsDeleter) Delete() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), o.timeout)
	defer cancelFunc()

	err := o.client.V2Alpha1().ObservabilityTracings().Delete(ctx, o.object.Name())
	if meshclient.IsNotFoundError(err) {
		return errors.Wrapf(err, "delete observabilityTracings %s", o.object.Name())
	}

	return err
}

type observabilityMetricsDeleter struct {
	baseDeleter
	object *resource.ObservabilityMetrics
}

func (o *observabilityMetricsDeleter) Delete() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), o.timeout)
	defer cancelFunc()

	err := o.client.V2Alpha1().ObservabilityMetrics().Delete(ctx, o.object.Name())
	if meshclient.IsNotFoundError(err) {
		return errors.Wrapf(err, "delete observabilityMetrics %s", o.object.Name())
	}

	return err
}

type observabilityOutputServerDeleter struct {
	baseDeleter
	object *resource.ObservabilityOutputServer
}

func (o *observabilityOutputServerDeleter) Delete() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), o.timeout)
	defer cancelFunc()

	err := o.client.V2Alpha1().ObservabilityOutputServer().Delete(ctx, o.object.Name())
	if meshclient.IsNotFoundError(err) {
		return errors.Wrapf(err, "delete observabilityOutputServer %s", o.object.Name())
	}

	return err
}

type loadBalanceDeleter struct {
	baseDeleter
	object *resource.LoadBalance
}

func (l *loadBalanceDeleter) Delete() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), l.timeout)
	defer cancelFunc()

	err := l.client.V2Alpha1().LoadBalance().Delete(ctx, l.object.Name())
	if meshclient.IsNotFoundError(err) {
		return errors.Wrapf(err, "delete loadBalance %s", l.object.Name())
	}

	return err
}

type tenantDeleter struct {
	baseDeleter
	object *resource.Tenant
}

func (t *tenantDeleter) Delete() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), t.timeout)
	defer cancelFunc()

	err := t.client.V2Alpha1().Tenant().Delete(ctx, t.object.Name())
	if meshclient.IsNotFoundError(err) {
		return errors.Wrapf(err, "delete tenant %s", t.object.Name())
	}

	return err
}

type resilienceDeleter struct {
	baseDeleter
	object *resource.Resilience
}

func (r *resilienceDeleter) Delete() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), r.timeout)
	defer cancelFunc()

	err := r.client.V2Alpha1().Resilience().Delete(ctx, r.object.Name())
	if meshclient.IsNotFoundError(err) {
		return errors.Wrapf(err, "delete resilience %s", r.object.Name())
	}

	return err
}

type mockDeleter struct {
	baseDeleter
	object *resource.Mock
}

func (m *mockDeleter) Delete() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), m.timeout)
	defer cancelFunc()

	err := m.client.V2Alpha1().Mock().Delete(ctx, m.object.Name())
	if meshclient.IsNotFoundError(err) {
		return errors.Wrapf(err, "delete mock %s", m.object.Name())
	}

	return err
}

type ingressDeleter struct {
	baseDeleter
	object *resource.Ingress
}

func (i *ingressDeleter) Delete() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), i.timeout)
	defer cancelFunc()

	err := i.client.V2Alpha1().Ingress().Delete(ctx, i.object.Name())
	if meshclient.IsNotFoundError(err) {
		return errors.Wrapf(err, "delete ingress %s", i.object.Name())
	}

	return err
}

type httpRouteGroupDeleter struct {
	baseDeleter
	object *resource.HTTPRouteGroup
}

func (g *httpRouteGroupDeleter) Delete() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), g.timeout)
	defer cancelFunc()

	err := g.client.V2Alpha1().HTTPRouteGroup().Delete(ctx, g.object.Name())
	if meshclient.IsNotFoundError(err) {
		return errors.Wrapf(err, "delete http route group %s", g.object.Name())
	}

	return err
}

type trafficTargetDeleter struct {
	baseDeleter
	object *resource.TrafficTarget
}

func (tt *trafficTargetDeleter) Delete() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), tt.timeout)
	defer cancelFunc()

	err := tt.client.V2Alpha1().HTTPRouteGroup().Delete(ctx, tt.object.Name())
	if meshclient.IsNotFoundError(err) {
		return errors.Wrapf(err, "delete traffic target %s", tt.object.Name())
	}

	return err
}

type serviceCanaryDeleter struct {
	baseDeleter
	object *resource.ServiceCanary
}

func (sc *serviceCanaryDeleter) Delete() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), sc.timeout)
	defer cancelFunc()

	err := sc.client.V2Alpha1().ServiceCanary().Delete(ctx, sc.object.Name())
	if meshclient.IsNotFoundError(err) {
		return errors.Wrapf(err, "delete serviceCanary %s", sc.object.Name())
	}

	return err
}

type customResourceKindDeleter struct {
	baseDeleter
	object *resource.CustomResourceKind
}

func (k *customResourceKindDeleter) Delete() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), k.timeout)
	defer cancelFunc()

	err := k.client.V2Alpha1().CustomResourceKind().Delete(ctx, k.object.Name())
	if meshclient.IsNotFoundError(err) {
		return errors.Wrapf(err, "delete custom resource kind %s", k.object.Name())
	}

	return err
}

type customResourceDeleter struct {
	baseDeleter
	object *resource.CustomResource
}

func (crd *customResourceDeleter) Delete() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), crd.timeout)
	defer cancelFunc()

	err := crd.client.V2Alpha1().CustomResource().Delete(ctx, crd.object.Kind(), crd.object.Name())
	if meshclient.IsNotFoundError(err) {
		return errors.Wrapf(err, "delete custom resource %s", crd.object.Name())
	}

	return err
}
