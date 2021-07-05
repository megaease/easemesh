package get

import (
	"context"
	"time"

	"github.com/megaease/easemeshctl/cmd/client/command/meshclient"
	"github.com/megaease/easemeshctl/cmd/client/resource"
	"github.com/megaease/easemeshctl/cmd/common"
)

type (
	Getter interface {
		Get() error
	}

	baseGetter struct {
		client  meshclient.MeshClient
		timeout time.Duration
		printer *printer
	}
)

var _ Getter = &serviceGetter{}

type serviceGetter struct {
	baseGetter
	object *resource.Service
}

func WrapGetterByMeshObject(object resource.MeshObject,
	client meshclient.MeshClient, timeout time.Duration, outputFormat string) Getter {

	base := baseGetter{
		client:  client,
		timeout: timeout,
		printer: newPrinter(outputFormat),
	}

	switch object.Kind() {
	case resource.KindService:
		return &serviceGetter{object: object.(*resource.Service), baseGetter: base}
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
	default:
		common.ExitWithErrorf("BUG: unsupported kind: %s", object.Kind())
	}

	return nil
}

func (s *serviceGetter) Get() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), s.timeout)
	defer cancelFunc()

	if s.object.Name() != "" {
		service, err := s.client.V1Alpha1().Service().Get(ctx, s.object.Name())
		if err != nil {
			return err
		}

		s.printer.printObjects([]resource.MeshObject{service})
	} else {
		services, err := s.client.V1Alpha1().Service().List(ctx)
		if err != nil {
			return err
		}

		objects := make([]resource.MeshObject, len(services))
		for i := range services {
			objects[i] = services[i]
		}

		s.printer.printObjects(objects)
	}

	return nil
}

type canaryGetter struct {
	baseGetter
	object *resource.Canary
}

func (c *canaryGetter) Get() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFunc()

	if c.object.Name() != "" {
		canary, err := c.client.V1Alpha1().Canary().Get(ctx, c.object.Name())
		if err != nil {
			return err
		}

		c.printer.printObjects([]resource.MeshObject{canary})
	} else {
		canaries, err := c.client.V1Alpha1().Canary().List(ctx)
		if err != nil {
			return err
		}

		objects := make([]resource.MeshObject, len(canaries))
		for i := range canaries {
			objects[i] = canaries[i]
		}

		c.printer.printObjects(objects)
	}

	return nil
}

type observabilityTracingsGetter struct {
	baseGetter
	object *resource.ObservabilityTracings
}

func (o *observabilityTracingsGetter) Get() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), o.timeout)
	defer cancelFunc()

	if o.object.Name() != "" {
		tracings, err := o.client.V1Alpha1().ObservabilityTracings().Get(ctx, o.object.Name())
		if err != nil {
			return err
		}

		o.printer.printObjects([]resource.MeshObject{tracings})
	} else {
		tracings, err := o.client.V1Alpha1().ObservabilityTracings().List(ctx)
		if err != nil {
			return err
		}

		objects := make([]resource.MeshObject, len(tracings))
		for i := range tracings {
			objects[i] = tracings[i]
		}

		o.printer.printObjects(objects)
	}

	return nil
}

type observabilityMetricsGetter struct {
	baseGetter
	object *resource.ObservabilityMetrics
}

func (o *observabilityMetricsGetter) Get() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), o.timeout)
	defer cancelFunc()

	if o.object.Name() != "" {
		metrics, err := o.client.V1Alpha1().ObservabilityMetrics().Get(ctx, o.object.Name())
		if err != nil {
			return err
		}

		o.printer.printObjects([]resource.MeshObject{metrics})
	} else {
		metrics, err := o.client.V1Alpha1().ObservabilityMetrics().List(ctx)
		if err != nil {
			return err
		}

		objects := make([]resource.MeshObject, len(metrics))
		for i := range metrics {
			objects[i] = metrics[i]
		}

		o.printer.printObjects(objects)
	}

	return nil
}

type observabilityOutputServerGetter struct {
	baseGetter
	object *resource.ObservabilityOutputServer
}

func (o *observabilityOutputServerGetter) Get() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), o.timeout)
	defer cancelFunc()

	if o.object.Name() != "" {
		server, err := o.client.V1Alpha1().ObservabilityOutputServer().Get(ctx, o.object.Name())
		if err != nil {
			return err
		}

		o.printer.printObjects([]resource.MeshObject{server})
	} else {
		servers, err := o.client.V1Alpha1().ObservabilityOutputServer().List(ctx)
		if err != nil {
			return err
		}

		objects := make([]resource.MeshObject, len(servers))
		for i := range servers {
			objects[i] = servers[i]
		}

		o.printer.printObjects(objects)
	}

	return nil
}

type loadBalanceGetter struct {
	baseGetter
	object *resource.LoadBalance
}

func (l *loadBalanceGetter) Get() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), l.timeout)
	defer cancelFunc()

	if l.object.Name() != "" {
		lb, err := l.client.V1Alpha1().LoadBalance().Get(ctx, l.object.Name())
		if err != nil {
			return err
		}

		l.printer.printObjects([]resource.MeshObject{lb})
	} else {
		lbs, err := l.client.V1Alpha1().LoadBalance().List(ctx)
		if err != nil {
			return err
		}

		objects := make([]resource.MeshObject, len(lbs))
		for i := range lbs {
			objects[i] = lbs[i]
		}

		l.printer.printObjects(objects)
	}

	return nil
}

type tenantGetter struct {
	baseGetter
	object *resource.Tenant
}

func (t *tenantGetter) Get() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), t.timeout)
	defer cancelFunc()

	if t.object.Name() != "" {
		tenant, err := t.client.V1Alpha1().Tenant().Get(ctx, t.object.Name())
		if err != nil {
			return err
		}

		t.printer.printObjects([]resource.MeshObject{tenant})
	} else {
		tenants, err := t.client.V1Alpha1().Tenant().List(ctx)
		if err != nil {
			return err
		}

		objects := make([]resource.MeshObject, len(tenants))
		for i := range tenants {
			objects[i] = tenants[i]
		}

		t.printer.printObjects(objects)
	}

	return nil
}

type resilienceGetter struct {
	baseGetter
	object *resource.Resilience
}

func (r *resilienceGetter) Get() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), r.timeout)
	defer cancelFunc()

	if r.object.Name() != "" {
		resilience, err := r.client.V1Alpha1().Resilience().Get(ctx, r.object.Name())
		if err != nil {
			return err
		}

		r.printer.printObjects([]resource.MeshObject{resilience})
	} else {
		resiliences, err := r.client.V1Alpha1().Resilience().List(ctx)
		if err != nil {
			return err
		}

		objects := make([]resource.MeshObject, len(resiliences))
		for i := range resiliences {
			objects[i] = resiliences[i]
		}

		r.printer.printObjects(objects)
	}

	return nil
}

type ingressGetter struct {
	baseGetter
	object *resource.Ingress
}

func (i *ingressGetter) Get() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), i.timeout)
	defer cancelFunc()

	if i.object.Name() != "" {
		ingress, err := i.client.V1Alpha1().Ingress().Get(ctx, i.object.Name())
		if err != nil {
			return err
		}

		i.printer.printObjects([]resource.MeshObject{ingress})
	} else {
		ingresses, err := i.client.V1Alpha1().Ingress().List(ctx)
		if err != nil {
			return err
		}

		objects := make([]resource.MeshObject, len(ingresses))
		for i := range ingresses {
			objects[i] = ingresses[i]
		}

		i.printer.printObjects(objects)
	}

	return nil
}
