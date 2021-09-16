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

package delete

import (
	"context"
	"time"

	"github.com/megaease/easemeshctl/cmd/client/command/meshclient"
	"github.com/megaease/easemeshctl/cmd/client/resource"
	"github.com/megaease/easemeshctl/cmd/common"

	"github.com/pkg/errors"
)

// WrapDeleterByMeshObject returns a new Deleter from a MeshObject
func WrapDeleterByMeshObject(object resource.MeshObject,
	client meshclient.MeshClient, timeout time.Duration) Deleter {
	switch object.Kind() {
	case resource.KindMeshController:
		return &meshControllerDeleter{object: object.(*resource.MeshController), baseDeleter: baseDeleter{client: client, timeout: timeout}}
	case resource.KindService:
		return &serviceDeleter{object: object.(*resource.Service), baseDeleter: baseDeleter{client: client, timeout: timeout}}
	case resource.KindServiceInstance:
		return &serviceInstanceDeleter{object: object.(*resource.ServiceInstance), baseDeleter: baseDeleter{client: client, timeout: timeout}}
	case resource.KindCanary:
		return &canaryDeleter{object: object.(*resource.Canary), baseDeleter: baseDeleter{client: client, timeout: timeout}}
	case resource.KindLoadBalance:
		return &loadBalanceDeleter{object: object.(*resource.LoadBalance), baseDeleter: baseDeleter{client: client, timeout: timeout}}
	case resource.KindTenant:
		return &tenantDeleter{object: object.(*resource.Tenant), baseDeleter: baseDeleter{client: client, timeout: timeout}}
	case resource.KindResilience:
		return &resilienceDeleter{object: object.(*resource.Resilience), baseDeleter: baseDeleter{client: client, timeout: timeout}}
	case resource.KindObservabilityMetrics:
		return &observabilityMetricsDeleter{object: object.(*resource.ObservabilityMetrics), baseDeleter: baseDeleter{client: client, timeout: timeout}}
	case resource.KindObservabilityOutputServer:
		return &observabilityOutputServerDeleter{object: object.(*resource.ObservabilityOutputServer), baseDeleter: baseDeleter{client: client, timeout: timeout}}
	case resource.KindObservabilityTracings:
		return &observabilityTracingsDeleter{object: object.(*resource.ObservabilityTracings), baseDeleter: baseDeleter{client: client, timeout: timeout}}
	case resource.KindIngress:
		return &ingressDeleter{object: object.(*resource.Ingress), baseDeleter: baseDeleter{client: client, timeout: timeout}}
	default:
		common.ExitWithErrorf("BUG: unsupported kind: %s", object.Kind())
	}

	return nil
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
	return s.client.V1Alpha1().MeshController().Delete(ctx, s.object.Name())
}

type serviceDeleter struct {
	baseDeleter
	object *resource.Service
}

func (s *serviceDeleter) Delete() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), s.timeout)
	defer cancelFunc()
	return s.client.V1Alpha1().Service().Delete(ctx, s.object.Name())
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

	return s.client.V1Alpha1().ServiceInstance().Delete(ctx,
		serviceName, instanceID)
}

type canaryDeleter struct {
	baseDeleter
	object *resource.Canary
}

func (c *canaryDeleter) Delete() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFunc()

	err := c.client.V1Alpha1().Canary().Delete(ctx, c.object.Name())
	if meshclient.IsNotFoundError(err) {
		return errors.Wrapf(err, "delete canary %s", c.object.Name())
	}

	return err
}

type observabilityTracingsDeleter struct {
	baseDeleter
	object *resource.ObservabilityTracings
}

func (o *observabilityTracingsDeleter) Delete() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), o.timeout)
	defer cancelFunc()

	err := o.client.V1Alpha1().ObservabilityTracings().Delete(ctx, o.object.Name())
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

	err := o.client.V1Alpha1().ObservabilityMetrics().Delete(ctx, o.object.Name())
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

	err := o.client.V1Alpha1().ObservabilityOutputServer().Delete(ctx, o.object.Name())
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

	err := l.client.V1Alpha1().LoadBalance().Delete(ctx, l.object.Name())
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

	err := t.client.V1Alpha1().Tenant().Delete(ctx, t.object.Name())
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

	err := r.client.V1Alpha1().Resilience().Delete(ctx, r.object.Name())
	if meshclient.IsNotFoundError(err) {
		return errors.Wrapf(err, "delete resilience %s", r.object.Name())
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

	err := i.client.V1Alpha1().Ingress().Delete(ctx, i.object.Name())
	if meshclient.IsNotFoundError(err) {
		return errors.Wrapf(err, "delete ingress %s", i.object.Name())
	}

	return err
}
