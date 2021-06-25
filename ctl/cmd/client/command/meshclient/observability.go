package meshclient

import (
	"context"
	"fmt"
	"net/http"

	"github.com/megaease/easemesh-api/v1alpha1"
	"github.com/megaease/easemeshctl/cmd/client/resource"
	"github.com/megaease/easemeshctl/cmd/common/client"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type observabilityGetter struct {
	client *meshClient
}

func (o *observabilityGetter) ObservabilityTracings() ObservabilityTracingInterface {
	return &observabilityTracingInterface{client: o.client}
}
func (o *observabilityGetter) ObservabilityMetrics() ObservabilityMetricInterface {
	return &observabilityMetricInterface{client: o.client}
}
func (o *observabilityGetter) ObservabilityOutputServer() ObservabilityOutputServerInterface {
	return &observabilityOutputServerInterface{client: o.client}
}

type observabilityTracingInterface struct {
	client *meshClient
}

func (o *observabilityTracingInterface) Get(ctx context.Context, serviceID string) (*resource.ObservabilityTracings, error) {
	jsonClient := client.NewHTTPJSON()
	r, err := jsonClient.
		GetByContext(fmt.Sprintf("http://"+o.client.server+MeshServiceTracingsURL, serviceID), nil, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "get observabilityTracings %s error", serviceID)
			}

			if statusCode >= 300 {
				return nil, errors.Errorf("call %s%s failed, return status code: %d text:%s", o.client.server, MeshServiceTracingsURL, statusCode, string(b))
			}
			tracing := &v1alpha1.ObservabilityTracings{}
			err := yaml.Unmarshal(b, tracing)
			if err != nil {
				return nil, errors.Wrap(err, "unmarshal data to v1alpha1.ObservabilityTracings error")
			}
			return resource.ToObservabilityTracings(serviceID, tracing), nil
		})
	if err != nil {
		return nil, err
	}

	return r.(*resource.ObservabilityTracings), nil
}

func (o *observabilityTracingInterface) Patch(ctx context.Context, tracings *resource.ObservabilityTracings) error {
	jsonClient := client.NewHTTPJSON()
	update := tracings.ToV1Alpha1()
	_, err := jsonClient.
		PutByContext(fmt.Sprintf("http://"+o.client.server+MeshServiceTracingsURL, tracings.Name()), &update, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "patch observabilityTracings %s error", tracings.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call PUT %s%s failed, return statuscode %d text %s", o.client.server, MeshServiceTracingsURL, statusCode, string(b))
		})
	return err
}

func (o *observabilityTracingInterface) Create(ctx context.Context, tracings *resource.ObservabilityTracings) error {
	created := tracings.ToV1Alpha1()
	_, err := client.NewHTTPJSON().
		PostByContext(fmt.Sprintf("http://"+o.client.server+MeshServiceTracingsURL, tracings.Name()), &created, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusConflict {
				return nil, errors.Wrapf(ConflictError, "create observabilityTracings %s error", tracings.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call Post %s%s failed, return statuscode %d text %s", o.client.server, MeshServiceTracingsURL, statusCode, string(b))
		})
	return err
}

func (o *observabilityTracingInterface) Delete(ctx context.Context, serviceID string) error {
	_, err := client.NewHTTPJSON().
		PostByContext(fmt.Sprintf("http://"+o.client.server+MeshServiceTracingsURL, serviceID), nil, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "delete observabilityTracings %s error", serviceID)
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call DELETE %s%s failed, return statuscode %d text %s", o.client.server, MeshServiceTracingsURL, statusCode, string(b))
		})
	return err
}

func (o *observabilityTracingInterface) List(ctx context.Context) ([]resource.ObservabilityTracings, error) {
	result, err := client.NewHTTPJSON().
		GetByContext(fmt.Sprintf("http://"+o.client.server+MeshServicesURL), nil, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrap(NotFoundError, "list service error")
			}

			if statusCode >= 300 || statusCode < 200 {
				return nil, errors.Errorf("call GET %s%s failed, return statuscode %d text %s", o.client.server, MeshServicesURL, statusCode, string(b))
			}

			services := []v1alpha1.Service{}
			err := yaml.Unmarshal(b, services)
			if err != nil {
				return nil, errors.Wrapf(err, "unmarshal services result error")
			}
			results := []resource.ObservabilityTracings{}
			for _, ss := range services {
				if ss.Observability != nil && ss.Observability.Tracings != nil {
					results = append(results, resource.ToObservabilityTracings(ss.Name, ss.Observability.Tracings))
				}
			}
			return results, nil
		})

	return result.([]resource.ObservabilityTracings), err
}

type observabilityMetricInterface struct {
	client *meshClient
}

func (o *observabilityMetricInterface) Get(ctx context.Context, serviceID string) (*resource.ObservabilityMetrics, error) {
	jsonClient := client.NewHTTPJSON()
	r, err := jsonClient.
		GetByContext(fmt.Sprintf("http://"+o.client.server+MeshServiceMetricsURL, serviceID), nil, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "get observabilityMetrics %s error", serviceID)
			}

			if statusCode >= 300 {
				return nil, errors.Errorf("call %s%s failed, return status code: %d text:%s", o.client.server, MeshServiceMetricsURL, statusCode, string(b))
			}
			metrics := &v1alpha1.ObservabilityMetrics{}
			err := yaml.Unmarshal(b, metrics)
			if err != nil {
				return nil, errors.Wrap(err, "unmarshal data to v1alpha1.Service error")
			}
			return resource.ToObservabilityMetrics(serviceID, metrics), nil
		})
	if err != nil {
		return nil, err
	}

	return r.(*resource.ObservabilityMetrics), nil
}

func (o *observabilityMetricInterface) Patch(ctx context.Context, metrics *resource.ObservabilityMetrics) error {
	jsonClient := client.NewHTTPJSON()
	update := metrics.ToV1Alpha1()
	_, err := jsonClient.
		PutByContext(fmt.Sprintf("http://"+o.client.server+MeshServiceMetricsURL, metrics.Name()), &update, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "patch observabilityMetrics %s error", metrics.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call PUT %s%s failed, return statuscode %d text %s", o.client.server, MeshServiceMetricsURL, statusCode, string(b))
		})
	return err
}

func (o *observabilityMetricInterface) Create(ctx context.Context, metrics *resource.ObservabilityMetrics) error {
	created := metrics.ToV1Alpha1()
	_, err := client.NewHTTPJSON().
		PostByContext(fmt.Sprintf("http://"+o.client.server+MeshServiceMetricsURL, metrics.Name()), &created, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusConflict {
				return nil, errors.Wrapf(ConflictError, "create observabilityMetrics %s error", metrics.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call Post %s%s failed, return statuscode %d text %s", o.client.server, MeshServiceMetricsURL, statusCode, string(b))
		})
	return err
}

func (o *observabilityMetricInterface) Delete(ctx context.Context, serviceID string) error {
	_, err := client.NewHTTPJSON().
		PostByContext(fmt.Sprintf("http://"+o.client.server+MeshServiceMetricsURL, serviceID), nil, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "delete observabilityMetrics %s error", serviceID)
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call DELETE %s%s failed, return statuscode %d text %s", o.client.server, MeshServiceMetricsURL, statusCode, string(b))
		})
	return err
}

func (o *observabilityMetricInterface) List(ctx context.Context) ([]resource.ObservabilityMetrics, error) {
	result, err := client.NewHTTPJSON().
		GetByContext(fmt.Sprintf("http://"+o.client.server+MeshServicesURL), nil, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrap(NotFoundError, "list service error")
			}

			if statusCode >= 300 || statusCode < 200 {
				return nil, errors.Errorf("call GET %s%s failed, return statuscode %d text %s", o.client.server, MeshServicesURL, statusCode, string(b))
			}

			services := []v1alpha1.Service{}
			err := yaml.Unmarshal(b, services)
			if err != nil {
				return nil, errors.Wrapf(err, "unmarshal services result error")
			}
			results := []resource.ObservabilityMetrics{}
			for _, ss := range services {
				if ss.Observability != nil && ss.Observability.Metrics != nil {
					results = append(results, resource.ToObservabilityMetrics(ss.Name, ss.Observability.Metrics))
				}
			}
			return results, nil
		})

	return result.([]resource.ObservabilityMetrics), err
}

type observabilityOutputServerInterface struct {
	client *meshClient
}

func (o *observabilityOutputServerInterface) Get(ctx context.Context, serviceID string) (*resource.ObservabilityOutputServer, error) {
	jsonClient := client.NewHTTPJSON()
	r, err := jsonClient.
		GetByContext(fmt.Sprintf("http://"+o.client.server+MeshServiceOutputServerURL, serviceID), nil, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "get observabilityOutputServer %s error", serviceID)
			}

			if statusCode >= 300 {
				return nil, errors.Errorf("call %s%s failed, return status code: %d text:%s", o.client.server, MeshServiceOutputServerURL, statusCode, string(b))
			}
			output := &v1alpha1.ObservabilityOutputServer{}
			err := yaml.Unmarshal(b, output)
			if err != nil {
				return nil, errors.Wrap(err, "unmarshal data to v1alpha1.ObservabilityOutputServer error")
			}
			return resource.ToObservabilityOutputServer(serviceID, output), nil
		})
	if err != nil {
		return nil, err
	}

	return r.(*resource.ObservabilityOutputServer), nil
}

func (o *observabilityOutputServerInterface) Patch(ctx context.Context, output *resource.ObservabilityOutputServer) error {
	jsonClient := client.NewHTTPJSON()
	update := output.ToV1Alpha1()
	_, err := jsonClient.
		PutByContext(fmt.Sprintf("http://"+o.client.server+MeshServiceOutputServerURL, output.Name()), &update, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "patch observabilityOutputServer %s error", output.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call PUT %s%s failed, return statuscode %d text %s", o.client.server, MeshServiceOutputServerURL, statusCode, string(b))
		})
	return err
}

func (o *observabilityOutputServerInterface) Create(ctx context.Context, output *resource.ObservabilityOutputServer) error {
	created := output.ToV1Alpha1()
	_, err := client.NewHTTPJSON().
		PostByContext(fmt.Sprintf("http://"+o.client.server+MeshServiceOutputServerURL, output.Name()), &created, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusConflict {
				return nil, errors.Wrapf(ConflictError, "create observabilityOutputServer %s error", output.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call Post %s%s failed, return statuscode %d text %s", o.client.server, MeshServiceOutputServerURL, statusCode, string(b))
		})
	return err
}

func (o *observabilityOutputServerInterface) Delete(ctx context.Context, serviceID string) error {
	_, err := client.NewHTTPJSON().
		PostByContext(fmt.Sprintf("http://"+o.client.server+MeshServiceOutputServerURL, serviceID), nil, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "delete observabilityOutputServer %s error", serviceID)
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call DELETE %s%s failed, return statuscode %d text %s", o.client.server, MeshServiceOutputServerURL, statusCode, string(b))
		})
	return err
}

func (o *observabilityOutputServerInterface) List(ctx context.Context) ([]resource.ObservabilityOutputServer, error) {
	result, err := client.NewHTTPJSON().
		GetByContext(fmt.Sprintf("http://"+o.client.server+MeshServicesURL), nil, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrap(NotFoundError, "list service error")
			}

			if statusCode >= 300 || statusCode < 200 {
				return nil, errors.Errorf("call GET %s%s failed, return statuscode %d text %s", o.client.server, MeshServicesURL, statusCode, string(b))
			}

			services := []v1alpha1.Service{}
			err := yaml.Unmarshal(b, services)
			if err != nil {
				return nil, errors.Wrapf(err, "unmarshal services result error")
			}
			results := []resource.ObservabilityOutputServer{}
			for _, ss := range services {
				if ss.Observability != nil && ss.Observability.OutputServer != nil {
					results = append(results, resource.ToObservabilityOutputServer(ss.Name, ss.Observability.OutputServer))
				}
			}
			return results, nil
		})

	return result.([]resource.ObservabilityOutputServer), err
}
