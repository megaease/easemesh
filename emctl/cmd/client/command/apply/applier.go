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
	case resource.KindObservabilityMetrics:
		return &observabilityMetricsApplier{object: object.(*resource.ObservabilityMetrics), baseApplier: baseApplier{client: client, timeout: timeout}}
	case resource.KindObservabilityOutputServer:
		return &observabilityOutputServerApplier{object: object.(*resource.ObservabilityOutputServer), baseApplier: baseApplier{client: client, timeout: timeout}}
	case resource.KindObservabilityTracings:
		return &observabilityTracingsApplier{object: object.(*resource.ObservabilityTracings), baseApplier: baseApplier{client: client, timeout: timeout}}
	case resource.KindIngress:
		return &ingressApplier{object: object.(*resource.Ingress), baseApplier: baseApplier{client: client, timeout: timeout}}
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
	for {
		switch {
		case err == nil:
			return nil
		case meshclient.IsConflictError(err):
			err = mc.client.V1Alpha1().MeshController().Patch(ctx, mc.object)
			if err != nil {
				return errors.Wrapf(err, "update meshController %s", mc.object.Name())
			}
		case meshclient.IsNotFoundError(err):
			err = mc.client.V1Alpha1().MeshController().Create(ctx, mc.object)
			if err != nil {
				return errors.Wrapf(err, "create meshController %s", mc.object.Name())
			}
		default:
			return errors.Wrapf(err, "apply meshController %s", mc.object.Name())
		}

	}
}

type serviceApplier struct {
	baseApplier
	object *resource.Service
}

func (s *serviceApplier) Apply() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), s.timeout)
	defer cancelFunc()
	err := s.client.V1Alpha1().Service().Create(ctx, s.object)
	for {
		switch {
		case err == nil:
			return nil
		case meshclient.IsConflictError(err):
			err = s.client.V1Alpha1().Service().Patch(ctx, s.object)
			if err != nil {
				return errors.Wrapf(err, "update service %s", s.object.Name())
			}
		case meshclient.IsNotFoundError(err):
			err = s.client.V1Alpha1().Service().Create(ctx, s.object)
			if err != nil {
				return errors.Wrapf(err, "create service %s", s.object.Name())
			}
		default:
			return errors.Wrapf(err, "apply service %s", s.object.Name())
		}

	}
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
	for {
		switch {
		case err == nil:
			return nil
		case meshclient.IsConflictError(err):
			err = c.client.V1Alpha1().Canary().Patch(ctx, c.object)
			if err != nil {
				return errors.Wrapf(err, "update canary %s", c.object.Name())
			}
		case meshclient.IsNotFoundError(err):
			err = c.client.V1Alpha1().Canary().Create(ctx, c.object)
			if err != nil {
				return errors.Wrapf(err, "create canary %s", c.object.Name())
			}
		default:
			return errors.Wrapf(err, "apply canary %s", c.object.Name())
		}

	}
}

type observabilityTracingsApplier struct {
	baseApplier
	object *resource.ObservabilityTracings
}

func (o *observabilityTracingsApplier) Apply() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), o.timeout)
	defer cancelFunc()
	err := o.client.V1Alpha1().ObservabilityTracings().Create(ctx, o.object)
	for {
		switch {
		case err == nil:
			return nil
		case meshclient.IsConflictError(err):
			err = o.client.V1Alpha1().ObservabilityTracings().Patch(ctx, o.object)
			if err != nil {
				return errors.Wrapf(err, "update observabilityTracings %s", o.object.Name())
			}
		case meshclient.IsNotFoundError(err):
			err = o.client.V1Alpha1().ObservabilityTracings().Create(ctx, o.object)
			if err != nil {
				return errors.Wrapf(err, "create observabilityTracings %s", o.object.Name())
			}
		default:
			return errors.Wrapf(err, "apply observabilityTracings %s", o.object.Name())
		}

	}
}

type observabilityMetricsApplier struct {
	baseApplier
	object *resource.ObservabilityMetrics
}

func (o *observabilityMetricsApplier) Apply() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), o.timeout)
	defer cancelFunc()
	err := o.client.V1Alpha1().ObservabilityMetrics().Create(ctx, o.object)
	for {
		switch {
		case err == nil:
			return nil
		case meshclient.IsConflictError(err):
			err = o.client.V1Alpha1().ObservabilityMetrics().Patch(ctx, o.object)
			if err != nil {
				return errors.Wrapf(err, "update observabilityMetrics %s", o.object.Name())
			}
		case meshclient.IsNotFoundError(err):
			err = o.client.V1Alpha1().ObservabilityMetrics().Create(ctx, o.object)
			if err != nil {
				return errors.Wrapf(err, "create observabilityMetrics %s", o.object.Name())
			}
		default:
			return errors.Wrapf(err, "apply observabilityMetrics %s", o.object.Name())
		}

	}
}

type observabilityOutputServerApplier struct {
	baseApplier
	object *resource.ObservabilityOutputServer
}

func (o *observabilityOutputServerApplier) Apply() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), o.timeout)
	defer cancelFunc()
	err := o.client.V1Alpha1().ObservabilityOutputServer().Create(ctx, o.object)
	for {
		switch {
		case err == nil:
			return nil
		case meshclient.IsConflictError(err):
			err = o.client.V1Alpha1().ObservabilityOutputServer().Patch(ctx, o.object)
			if err != nil {
				return errors.Wrapf(err, "update observabilityOutputServer %s", o.object.Name())
			}
		case meshclient.IsNotFoundError(err):
			err = o.client.V1Alpha1().ObservabilityOutputServer().Create(ctx, o.object)
			if err != nil {
				return errors.Wrapf(err, "create observabilityOutputServer %s", o.object.Name())
			}
		default:
			return errors.Wrapf(err, "apply observabilityOutputServer %s", o.object.Name())
		}

	}
}

type loadBalanceApplier struct {
	baseApplier
	object *resource.LoadBalance
}

func (l *loadBalanceApplier) Apply() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), l.timeout)
	defer cancelFunc()
	err := l.client.V1Alpha1().LoadBalance().Create(ctx, l.object)
	for {
		switch {
		case err == nil:
			return nil
		case meshclient.IsConflictError(err):
			err = l.client.V1Alpha1().LoadBalance().Patch(ctx, l.object)
			if err != nil {
				return errors.Wrapf(err, "update loadbalance %s", l.object.Name())
			}
		case meshclient.IsNotFoundError(err):
			err = l.client.V1Alpha1().LoadBalance().Create(ctx, l.object)
			if err != nil {
				return errors.Wrapf(err, "create loadbalance %s", l.object.Name())
			}
		default:
			return errors.Wrapf(err, "apply loadbalance %s", l.object.Name())
		}
	}
}

type tenantApplier struct {
	baseApplier
	object *resource.Tenant
}

func (t *tenantApplier) Apply() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), t.timeout)
	defer cancelFunc()
	err := t.client.V1Alpha1().Tenant().Create(ctx, t.object)
	for {
		switch {
		case err == nil:
			return nil
		case meshclient.IsConflictError(err):
			err = t.client.V1Alpha1().Tenant().Patch(ctx, t.object)
			if err != nil {
				return errors.Wrapf(err, "update tenant %s", t.object.Name())
			}
		case meshclient.IsNotFoundError(err):
			err = t.client.V1Alpha1().Tenant().Create(ctx, t.object)
			if err != nil {
				return errors.Wrapf(err, "create tenant %s", t.object.Name())
			}
		default:
			return errors.Wrapf(err, "apply tenant %s", t.object.Name())
		}
	}
}

type resilienceApplier struct {
	baseApplier
	object *resource.Resilience
}

func (r *resilienceApplier) Apply() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), r.timeout)
	defer cancelFunc()
	err := r.client.V1Alpha1().Resilience().Create(ctx, r.object)
	for {
		switch {
		case err == nil:
			return nil
		case meshclient.IsConflictError(err):
			err = r.client.V1Alpha1().Resilience().Patch(ctx, r.object)
			if err != nil {
				return errors.Wrapf(err, "update resilience %s", r.object.Name())
			}
		case meshclient.IsNotFoundError(err):
			err = r.client.V1Alpha1().Resilience().Create(ctx, r.object)
			if err != nil {
				return errors.Wrapf(err, "create resilience %s", r.object.Name())
			}
		default:
			return errors.Wrapf(err, "apply resilience %s", r.object.Name())
		}
	}
}

type ingressApplier struct {
	baseApplier
	object *resource.Ingress
}

func (i *ingressApplier) Apply() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), i.timeout)
	defer cancelFunc()
	err := i.client.V1Alpha1().Ingress().Create(ctx, i.object)
	for {
		switch {
		case err == nil:
			return nil
		case meshclient.IsConflictError(err):
			err = i.client.V1Alpha1().Ingress().Patch(ctx, i.object)
			if err != nil {
				return errors.Wrapf(err, "update resilience %s", i.object.Name())
			}
		case meshclient.IsNotFoundError(err):
			err = i.client.V1Alpha1().Ingress().Create(ctx, i.object)
			if err != nil {
				return errors.Wrapf(err, "create resilience %s", i.object.Name())
			}
		default:
			return errors.Wrapf(err, "apply resilience %s", i.object.Name())
		}
	}
}

type customResourceKindApplier struct {
	baseApplier
	object *resource.CustomResourceKind
}

func (k *customResourceKindApplier) Apply() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), k.timeout)
	defer cancelFunc()
	err := k.client.V1Alpha1().CustomResourceKind().Create(ctx, k.object)
>>>>>>> mai
	for {
		switch {
		case err == nil:
			return nil
		case meshclient.IsConflictError(err):
			err = k.client.V1Alpha1().CustomResourceKind().Patch(ctx, k.object)
			if err != nil {
				return errors.Wrapf(err, "update custom resource kind %s", k.object.Name())
			}
		case meshclient.IsNotFoundError(err):
			err = k.client.V1Alpha1().CustomResourceKind().Create(ctx, k.object)
			if err != nil {
				return errors.Wrapf(err, "create custom resource kind %s", k.object.Name())
			}
		default:
			return errors.Wrapf(err, "apply custom resource kind %s", k.object.Name())
		}
	}
}

type customResourceApplier struct {
	baseApplier
	object *resource.CustomResource
}

func (cra *customResourceApplier) Apply() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), cra.timeout)
	defer cancelFunc()
	err := cra.client.V1Alpha1().CustomResource().Create(ctx, cra.object)
	for {
		switch {
		case err == nil:
			return nil
		case meshclient.IsConflictError(err):
			err = cra.client.V1Alpha1().CustomResource().Patch(ctx, cra.object)
			if err != nil {
				return errors.Wrapf(err, "update custom resource %s", cra.object.Name())
			}
		case meshclient.IsNotFoundError(err):
			err = cra.client.V1Alpha1().CustomResource().Create(ctx, cra.object)
			if err != nil {
				return errors.Wrapf(err, "create custom resource %s", cra.object.Name())
			}
		default:
			return errors.Wrapf(err, "apply custom resource %s", cra.object.Name())
		}
	}
}