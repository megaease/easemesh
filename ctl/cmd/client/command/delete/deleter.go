package delete

import (
	"context"
	"time"

	"github.com/megaease/easemeshctl/cmd/client/command/meshclient"
	"github.com/megaease/easemeshctl/cmd/client/resource"
	"github.com/megaease/easemeshctl/cmd/common"

	"github.com/pkg/errors"
)

type Deleter interface {
	Delete() error
}

var _ Deleter = &serviceDeleter{}

type baseDeleter struct {
	client  meshclient.MeshClient
	timeout time.Duration
}

type serviceDeleter struct {
	baseDeleter
	object *resource.Service
}

func WrapDeleterByMeshObject(object resource.MeshObject,
	client meshclient.MeshClient, timeout time.Duration) Deleter {
	switch object.Kind() {
	case resource.KindService:
		return &serviceDeleter{object: object.(*resource.Service), baseDeleter: baseDeleter{client: client, timeout: timeout}}
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

func (s *serviceDeleter) Delete() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), s.timeout)
	defer cancelFunc()
	return s.client.V1Alpha1().Service().Delete(ctx, s.object.Name())
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
		return errors.Wrapf(err, "delete canary %s error", c.object.Name())
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
		return errors.Wrapf(err, "delete observabilityTracings %s error", o.object.Name())
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
		return errors.Wrapf(err, "delete observabilityMetrics %s error", o.object.Name())
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
		return errors.Wrapf(err, "delete observabilityOutputServer %s error", o.object.Name())
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
		return errors.Wrapf(err, "delete loadBalance %s error", l.object.Name())
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
		return errors.Wrapf(err, "delete tenant %s error", t.object.Name())
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
		return errors.Wrapf(err, "delete resilience %s error", r.object.Name())
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
		return errors.Wrapf(err, "delete ingress %s error", i.object.Name())
	}

	return err
}
