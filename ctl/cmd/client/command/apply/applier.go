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
	"github.com/megaease/easemeshctl/cmd/common"

	"github.com/pkg/errors"
)

type Applier interface {
	Apply() error
}

var _ Applier = &serviceApplier{}

type baseApplier struct {
	client  meshclient.MeshClient
	timeout time.Duration
}

type serviceApplier struct {
	baseApplier
	object *resource.Service
}

func WrapApplierByMeshObject(object resource.MeshObject,
	client meshclient.MeshClient, timeout time.Duration) Applier {
	switch object.Kind() {
	case resource.KindService:
		return &serviceApplier{object: object.(*resource.Service), baseApplier: baseApplier{client: client, timeout: timeout}}
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
	default:
		common.ExitWithErrorf("BUG: unsupported kind: %s", object.Kind())
	}

	return nil
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
				return errors.Wrapf(err, "update service %s error", s.object.Name())
			}
		case meshclient.IsNotFoundError(err):
			err = s.client.V1Alpha1().Service().Create(ctx, s.object)
			if err != nil {
				return errors.Wrapf(err, "create service %s error", s.object.Name())
			}
		default:
			return errors.Wrapf(err, "apply service %s error", s.object.Name())
		}

	}
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
				return errors.Wrapf(err, "update canary %s error", c.object.Name())
			}
		case meshclient.IsNotFoundError(err):
			err = c.client.V1Alpha1().Canary().Create(ctx, c.object)
			if err != nil {
				return errors.Wrapf(err, "create canary %s error", c.object.Name())
			}
		default:
			return errors.Wrapf(err, "apply canary %s error", c.object.Name())
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
				return errors.Wrapf(err, "update observabilityTracings %s error", o.object.Name())
			}
		case meshclient.IsNotFoundError(err):
			err = o.client.V1Alpha1().ObservabilityTracings().Create(ctx, o.object)
			if err != nil {
				return errors.Wrapf(err, "create observabilityTracings %s error", o.object.Name())
			}
		default:
			return errors.Wrapf(err, "apply observabilityTracings %s error", o.object.Name())
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
				return errors.Wrapf(err, "update observabilityMetrics %s error", o.object.Name())
			}
		case meshclient.IsNotFoundError(err):
			err = o.client.V1Alpha1().ObservabilityMetrics().Create(ctx, o.object)
			if err != nil {
				return errors.Wrapf(err, "create observabilityMetrics %s error", o.object.Name())
			}
		default:
			return errors.Wrapf(err, "apply observabilityMetrics %s error", o.object.Name())
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
				return errors.Wrapf(err, "update observabilityOutputServer %s error", o.object.Name())
			}
		case meshclient.IsNotFoundError(err):
			err = o.client.V1Alpha1().ObservabilityOutputServer().Create(ctx, o.object)
			if err != nil {
				return errors.Wrapf(err, "create observabilityOutputServer %s error", o.object.Name())
			}
		default:
			return errors.Wrapf(err, "apply observabilityOutputServer %s error", o.object.Name())
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
				return errors.Wrapf(err, "update loadbalance %s error", l.object.Name())
			}
		case meshclient.IsNotFoundError(err):
			err = l.client.V1Alpha1().LoadBalance().Create(ctx, l.object)
			if err != nil {
				return errors.Wrapf(err, "create loadbalance %s error", l.object.Name())
			}
		default:
			return errors.Wrapf(err, "apply loadbalance %s error", l.object.Name())
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
				return errors.Wrapf(err, "update tenant %s error", t.object.Name())
			}
		case meshclient.IsNotFoundError(err):
			err = t.client.V1Alpha1().Tenant().Create(ctx, t.object)
			if err != nil {
				return errors.Wrapf(err, "create tenant %s error", t.object.Name())
			}
		default:
			return errors.Wrapf(err, "apply tenant %s error", t.object.Name())
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
				return errors.Wrapf(err, "update resilience %s error", r.object.Name())
			}
		case meshclient.IsNotFoundError(err):
			err = r.client.V1Alpha1().Resilience().Create(ctx, r.object)
			if err != nil {
				return errors.Wrapf(err, "create resilience %s error", r.object.Name())
			}
		default:
			return errors.Wrapf(err, "apply resilience %s error", r.object.Name())
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
				return errors.Wrapf(err, "update resilience %s error", i.object.Name())
			}
		case meshclient.IsNotFoundError(err):
			err = i.client.V1Alpha1().Ingress().Create(ctx, i.object)
			if err != nil {
				return errors.Wrapf(err, "create resilience %s error", i.object.Name())
			}
		default:
			return errors.Wrapf(err, "apply resilience %s error", i.object.Name())
		}
	}
}
