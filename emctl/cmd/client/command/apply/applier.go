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
	"context"
	"time"

	"github.com/megaease/easemeshctl/cmd/client/command/meshclient"
	"github.com/megaease/easemeshctl/cmd/client/resource"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"

	"github.com/pkg/errors"
)

// Applier applies configuration to control plane service of the EaseMesh
type Applier interface {
	Apply() error
}

var _ Applier = &serviceApplier{}

type baseApplier struct {
	client  meshclient.MeshClient
	timeout time.Duration
}

// WrapApplierByMeshObject returns a Applier from a MeshObject
func WrapApplierByMeshObject(object meta.MeshObject,
	client meshclient.MeshClient, timeout time.Duration) Applier {
	switch object.Kind() {
	case resource.KindMeshController:
		return &meshControllerApplier{object: object.(*resource.MeshController), baseApplier: baseApplier{client: client, timeout: timeout}}
	case resource.KindService:
		return &serviceApplier{object: object.(*resource.Service), baseApplier: baseApplier{client: client, timeout: timeout}}
	case resource.KindServiceInstance:
		return &serviceInstanceApplier{object: object.(*resource.ServiceInstance), baseApplier: baseApplier{client: client, timeout: timeout}}
	case resource.KindCanary:
		return &canaryApplier{object: object.(*resource.Canary), baseApplier: baseApplier{client: client, timeout: timeout}}
	case resource.KindLoadBalance:
		return &loadBalanceApplier{object: object.(*resource.LoadBalance), baseApplier: baseApplier{client: client, timeout: timeout}}
	case resource.KindTenant:
		return &tenantApplier{object: object.(*resource.Tenant), baseApplier: baseApplier{client: client, timeout: timeout}}
	case resource.KindResilience:
		return &resilienceApplier{object: object.(*resource.Resilience), baseApplier: baseApplier{client: client, timeout: timeout}}
	case resource.KindMock:
		return &mockApplier{object: object.(*resource.Mock), baseApplier: baseApplier{client: client, timeout: timeout}}
	case resource.KindObservabilityMetrics:
		return &observabilityMetricsApplier{object: object.(*resource.ObservabilityMetrics), baseApplier: baseApplier{client: client, timeout: timeout}}
	case resource.KindObservabilityOutputServer:
		return &observabilityOutputServerApplier{object: object.(*resource.ObservabilityOutputServer), baseApplier: baseApplier{client: client, timeout: timeout}}
	case resource.KindObservabilityTracings:
		return &observabilityTracingsApplier{object: object.(*resource.ObservabilityTracings), baseApplier: baseApplier{client: client, timeout: timeout}}
	case resource.KindIngress:
		return &ingressApplier{object: object.(*resource.Ingress), baseApplier: baseApplier{client: client, timeout: timeout}}
	case resource.KindHTTPRouteGroup:
		return &httpRouteGroupApplier{object: object.(*resource.HTTPRouteGroup), baseApplier: baseApplier{client: client, timeout: timeout}}
	case resource.KindTrafficTarget:
		return &trafficTargetApplier{object: object.(*resource.TrafficTarget), baseApplier: baseApplier{client: client, timeout: timeout}}
	case resource.KindServiceCanary:
		return &serviceCanaryApplier{object: object.(*resource.ServiceCanary), baseApplier: baseApplier{client: client, timeout: timeout}}
	case resource.KindCustomResourceKind:
		return &customResourceKindApplier{object: object.(*resource.CustomResourceKind), baseApplier: baseApplier{client: client, timeout: timeout}}
	default:
		return &customResourceApplier{object: object.(*resource.CustomResource), baseApplier: baseApplier{client: client, timeout: timeout}}
	}
}

type meshControllerApplier struct {
	baseApplier
	object *resource.MeshController
}

func (mc *meshControllerApplier) Apply() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), mc.timeout)
	defer cancelFunc()
	err := mc.client.V1Alpha1().MeshController().Create(ctx, mc.object)
	if err != nil && meshclient.IsConflictError(err) {
		err = mc.client.V1Alpha1().MeshController().Patch(ctx, mc.object)
		if err != nil {
			return errors.Wrapf(err, "update meshController %s", mc.object.Name())
		}
	} else if err != nil {
		return errors.Wrapf(err, "apply meshController %s", mc.object.Name())
	}
	return nil
}

type serviceApplier struct {
	baseApplier
	object *resource.Service
}

func (s *serviceApplier) Apply() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), s.timeout)
	defer cancelFunc()
	err := s.client.V1Alpha1().Service().Create(ctx, s.object)
	if err != nil && meshclient.IsConflictError(err) {
		err = s.client.V1Alpha1().Service().Patch(ctx, s.object)
		if err != nil {
			return errors.Wrapf(err, "update service %s", s.object.Name())
		}
	} else if err != nil {
		return errors.Wrapf(err, "apply service %s", s.object.Name())
	}
	return nil
}

type serviceInstanceApplier struct {
	baseApplier
	object *resource.ServiceInstance
}

func (si *serviceInstanceApplier) Apply() error {
	return errors.New("not support applying service instance")
}

type canaryApplier struct {
	baseApplier
	object *resource.Canary
}

func (c *canaryApplier) Apply() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFunc()
	err := c.client.V1Alpha1().Canary().Create(ctx, c.object)
	if err != nil && meshclient.IsConflictError(err) {
		err = c.client.V1Alpha1().Canary().Patch(ctx, c.object)
		if err != nil {
			return errors.Wrapf(err, "update canary %s", c.object.Name())
		}
	} else if err != nil {
		return errors.Wrapf(err, "apply canary %s", c.object.Name())
	}
	return nil
}

type observabilityTracingsApplier struct {
	baseApplier
	object *resource.ObservabilityTracings
}

func (o *observabilityTracingsApplier) Apply() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), o.timeout)
	defer cancelFunc()
	err := o.client.V1Alpha1().ObservabilityTracings().Create(ctx, o.object)
	if err != nil && meshclient.IsConflictError(err) {
		err = o.client.V1Alpha1().ObservabilityTracings().Patch(ctx, o.object)
		if err != nil {
			return errors.Wrapf(err, "update observabilityTracings %s", o.object.Name())
		}
	} else if err != nil {
		return errors.Wrapf(err, "apply observabilityTracings %s", o.object.Name())
	}
	return nil
}

type observabilityMetricsApplier struct {
	baseApplier
	object *resource.ObservabilityMetrics
}

func (o *observabilityMetricsApplier) Apply() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), o.timeout)
	defer cancelFunc()
	err := o.client.V1Alpha1().ObservabilityMetrics().Create(ctx, o.object)
	if err != nil && meshclient.IsConflictError(err) {
		err = o.client.V1Alpha1().ObservabilityMetrics().Patch(ctx, o.object)
		if err != nil {
			return errors.Wrapf(err, "update observabilityMetrics %s", o.object.Name())
		}
	} else if err != nil {
		return errors.Wrapf(err, "apply observabilityMetrics %s", o.object.Name())
	}
	return nil
}

type observabilityOutputServerApplier struct {
	baseApplier
	object *resource.ObservabilityOutputServer
}

func (o *observabilityOutputServerApplier) Apply() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), o.timeout)
	defer cancelFunc()
	err := o.client.V1Alpha1().ObservabilityOutputServer().Create(ctx, o.object)
	if err != nil && meshclient.IsConflictError(err) {
		err = o.client.V1Alpha1().ObservabilityOutputServer().Patch(ctx, o.object)
		if err != nil {
			return errors.Wrapf(err, "update observabilityOutputServer %s", o.object.Name())
		}
	} else if err != nil {
		return errors.Wrapf(err, "apply observabilityOutputServer %s", o.object.Name())
	}
	return nil
}

type loadBalanceApplier struct {
	baseApplier
	object *resource.LoadBalance
}

func (l *loadBalanceApplier) Apply() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), l.timeout)
	defer cancelFunc()
	err := l.client.V1Alpha1().LoadBalance().Create(ctx, l.object)
	if err != nil && meshclient.IsConflictError(err) {
		err = l.client.V1Alpha1().LoadBalance().Patch(ctx, l.object)
		if err != nil {
			return errors.Wrapf(err, "update loadbalance %s", l.object.Name())
		}
	} else if err != nil {
		return errors.Wrapf(err, "apply loadbalance %s", l.object.Name())
	}
	return nil
}

type tenantApplier struct {
	baseApplier
	object *resource.Tenant
}

func (t *tenantApplier) Apply() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), t.timeout)
	defer cancelFunc()
	err := t.client.V1Alpha1().Tenant().Create(ctx, t.object)
	if err != nil && meshclient.IsConflictError(err) {
		err = t.client.V1Alpha1().Tenant().Patch(ctx, t.object)
		if err != nil {
			return errors.Wrapf(err, "update tenant %s", t.object.Name())
		}
	} else if err != nil {
		return errors.Wrapf(err, "apply tenant %s", t.object.Name())
	}
	return nil
}

type resilienceApplier struct {
	baseApplier
	object *resource.Resilience
}

func (r *resilienceApplier) Apply() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), r.timeout)
	defer cancelFunc()
	err := r.client.V1Alpha1().Resilience().Create(ctx, r.object)
	if err != nil && meshclient.IsConflictError(err) {
		err = r.client.V1Alpha1().Resilience().Patch(ctx, r.object)
		if err != nil {
			return errors.Wrapf(err, "update resilience %s", r.object.Name())
		}
	} else if err != nil {
		return errors.Wrapf(err, "apply resilience %s", r.object.Name())
	}
	return nil
}

type mockApplier struct {
	baseApplier
	object *resource.Mock
}

func (m *mockApplier) Apply() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), m.timeout)
	defer cancelFunc()
	err := m.client.V1Alpha1().Mock().Create(ctx, m.object)
	if err != nil && meshclient.IsConflictError(err) {
		err = m.client.V1Alpha1().Mock().Patch(ctx, m.object)
		if err != nil {
			return errors.Wrapf(err, "update mock %s", m.object.Name())
		}
	} else if err != nil {
		return errors.Wrapf(err, "apply mock %s", m.object.Name())
	}
	return nil
}

type ingressApplier struct {
	baseApplier
	object *resource.Ingress
}

func (i *ingressApplier) Apply() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), i.timeout)
	defer cancelFunc()
	err := i.client.V1Alpha1().Ingress().Create(ctx, i.object)
	if err != nil && meshclient.IsConflictError(err) {
		err = i.client.V1Alpha1().Ingress().Patch(ctx, i.object)
		if err != nil {
			return errors.Wrapf(err, "update ingress %s", i.object.Name())
		}
	} else if err != nil {
		return errors.Wrapf(err, "apply ingress %s", i.object.Name())
	}
	return nil
}

type httpRouteGroupApplier struct {
	baseApplier
	object *resource.HTTPRouteGroup
}

func (g *httpRouteGroupApplier) Apply() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), g.timeout)
	defer cancelFunc()
	err := g.client.V1Alpha1().HTTPRouteGroup().Create(ctx, g.object)
	if err != nil && meshclient.IsConflictError(err) {
		err = g.client.V1Alpha1().HTTPRouteGroup().Patch(ctx, g.object)
		if err != nil {
			return errors.Wrapf(err, "update httpRouteGroup %s", g.object.Name())
		}
	} else if err != nil {
		return errors.Wrapf(err, "apply httpRouteGroup %s", g.object.Name())
	}
	return nil
}

type trafficTargetApplier struct {
	baseApplier
	object *resource.TrafficTarget
}

func (tt *trafficTargetApplier) Apply() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), tt.timeout)
	defer cancelFunc()
	err := tt.client.V1Alpha1().TrafficTarget().Create(ctx, tt.object)
	if err != nil && meshclient.IsConflictError(err) {
		err = tt.client.V1Alpha1().TrafficTarget().Patch(ctx, tt.object)
		if err != nil {
			return errors.Wrapf(err, "update trafficTarget %s", tt.object.Name())
		}
	} else if err != nil {
		return errors.Wrapf(err, "apply trafficTarget %s", tt.object.Name())
	}
	return nil
}

type serviceCanaryApplier struct {
	baseApplier
	object *resource.ServiceCanary
}

func (sc *serviceCanaryApplier) Apply() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), sc.timeout)
	defer cancelFunc()
	err := sc.client.V1Alpha1().ServiceCanary().Create(ctx, sc.object)
	if err != nil && meshclient.IsConflictError(err) {
		err = sc.client.V1Alpha1().ServiceCanary().Patch(ctx, sc.object)
		if err != nil {
			return errors.Wrapf(err, "update serviceCanary %s", sc.object.Name())
		}
	} else if err != nil {
		return errors.Wrapf(err, "apply serviceCanary %s", sc.object.Name())
	}
	return nil
}

type customResourceKindApplier struct {
	baseApplier
	object *resource.CustomResourceKind
}

func (k *customResourceKindApplier) Apply() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), k.timeout)
	defer cancelFunc()
	err := k.client.V1Alpha1().CustomResourceKind().Create(ctx, k.object)
	if err != nil && meshclient.IsConflictError(err) {
		err = k.client.V1Alpha1().CustomResourceKind().Patch(ctx, k.object)
		if err != nil {
			return errors.Wrapf(err, "update custom resource kind %s", k.object.Name())
		}
	} else if err != nil {
		return errors.Wrapf(err, "apply custom resource kind %s", k.object.Name())
	}
	return nil
}

type customResourceApplier struct {
	baseApplier
	object *resource.CustomResource
}

func (cra *customResourceApplier) Apply() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), cra.timeout)
	defer cancelFunc()
	err := cra.client.V1Alpha1().CustomResource().Create(ctx, cra.object)
	if err != nil && meshclient.IsConflictError(err) {
		err = cra.client.V1Alpha1().CustomResource().Patch(ctx, cra.object)
		if err != nil {
			return errors.Wrapf(err, "update custom resource %s", cra.object.Name())
		}
	} else if err != nil {
		return errors.Wrapf(err, "apply custom resource %s", cra.object.Name())
	}
	return nil
}
