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

package meshclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/megaease/easemesh-api/v1alpha1"
	"github.com/megaease/easemeshctl/cmd/client/resource"
	"github.com/megaease/easemeshctl/cmd/common/client"

	"github.com/pkg/errors"
)

type observabilityGetter struct {
	client *meshClient
}

func (o *observabilityGetter) ObservabilityTracings() ObservabilityTracingsInterface {
	return &observabilityTracingInterface{client: o.client}
}
func (o *observabilityGetter) ObservabilityMetrics() ObservabilityMetricsInterface {
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
	url := fmt.Sprintf("http://"+o.client.server+MeshServiceTracingsURL, serviceID)
	r, err := jsonClient.
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "get observabilityTracings %s", serviceID)
			}

			if statusCode >= 300 {
				return nil, errors.Errorf("call %s failed, return status code: %d text:%s", url, statusCode, string(b))
			}
			tracing := &v1alpha1.ObservabilityTracings{}
			err := json.Unmarshal(b, tracing)
			if err != nil {
				return nil, errors.Wrap(err, "unmarshal data to v1alpha1.ObservabilityTracings")
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
	url := fmt.Sprintf("http://"+o.client.server+MeshServiceTracingsURL, tracings.Name())
	update := tracings.ToV1Alpha1()
	_, err := jsonClient.
		PutByContext(ctx, url, update, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "patch observabilityTracings %s", tracings.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call PUT %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (o *observabilityTracingInterface) Create(ctx context.Context, tracings *resource.ObservabilityTracings) error {
	created := tracings.ToV1Alpha1()
	url := fmt.Sprintf("http://"+o.client.server+MeshServiceTracingsURL, tracings.Name())
	_, err := client.NewHTTPJSON().
		PostByContext(ctx, url, created, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusConflict {
				return nil, errors.Wrapf(ConflictError, "create observabilityTracings %s", tracings.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call Post %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (o *observabilityTracingInterface) Delete(ctx context.Context, serviceID string) error {
	url := fmt.Sprintf("http://"+o.client.server+MeshServiceTracingsURL, serviceID)
	_, err := client.NewHTTPJSON().
		DeleteByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "delete observabilityTracings %s", serviceID)
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call DELETE %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (o *observabilityTracingInterface) List(ctx context.Context) ([]*resource.ObservabilityTracings, error) {
	url := "http://" + o.client.server + MeshServicesURL
	result, err := client.NewHTTPJSON().
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrap(NotFoundError, "list service")
			}

			if statusCode >= 300 || statusCode < 200 {
				return nil, errors.Errorf("call GET %s failed, return statuscode %d text %s", url, statusCode, string(b))
			}

			services := []v1alpha1.Service{}
			err := json.Unmarshal(b, &services)
			if err != nil {
				return nil, errors.Wrapf(err, "unmarshal services result")
			}
			results := []*resource.ObservabilityTracings{}
			for _, service := range services {
				if service.Observability != nil && service.Observability.Tracings != nil {
					results = append(results, resource.ToObservabilityTracings(service.Name, service.Observability.Tracings))
				}
			}
			return results, nil
		})

	return result.([]*resource.ObservabilityTracings), err
}

type observabilityMetricInterface struct {
	client *meshClient
}

func (o *observabilityMetricInterface) Get(ctx context.Context, serviceID string) (*resource.ObservabilityMetrics, error) {
	jsonClient := client.NewHTTPJSON()
	url := fmt.Sprintf("http://"+o.client.server+MeshServiceMetricsURL, serviceID)
	r, err := jsonClient.
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "get observabilityMetrics %s", serviceID)
			}

			if statusCode >= 300 {
				return nil, errors.Errorf("call %s failed, return status code: %d text:%s", url, statusCode, string(b))
			}
			metrics := &v1alpha1.ObservabilityMetrics{}
			err := json.Unmarshal(b, metrics)
			if err != nil {
				return nil, errors.Wrap(err, "unmarshal data to v1alpha1.Service")
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
	url := fmt.Sprintf("http://"+o.client.server+MeshServiceMetricsURL, metrics.Name())
	update := metrics.ToV1Alpha1()
	_, err := jsonClient.
		PutByContext(ctx, url, update, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "patch observabilityMetrics %s", metrics.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call PUT %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (o *observabilityMetricInterface) Create(ctx context.Context, metrics *resource.ObservabilityMetrics) error {
	url := fmt.Sprintf("http://"+o.client.server+MeshServiceMetricsURL, metrics.Name())
	created := metrics.ToV1Alpha1()
	_, err := client.NewHTTPJSON().
		PostByContext(ctx, url, created, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusConflict {
				return nil, errors.Wrapf(ConflictError, "create observabilityMetrics %s", metrics.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call Post %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (o *observabilityMetricInterface) Delete(ctx context.Context, serviceID string) error {
	url := fmt.Sprintf("http://"+o.client.server+MeshServiceMetricsURL, serviceID)
	_, err := client.NewHTTPJSON().
		DeleteByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "delete observabilityMetrics %s", serviceID)
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call DELETE %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (o *observabilityMetricInterface) List(ctx context.Context) ([]*resource.ObservabilityMetrics, error) {
	url := fmt.Sprintf("http://" + o.client.server + MeshServicesURL)
	result, err := client.NewHTTPJSON().
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrap(NotFoundError, "list service")
			}

			if statusCode >= 300 || statusCode < 200 {
				return nil, errors.Errorf("call GET %s failed, return statuscode %d text %s", url, statusCode, string(b))
			}

			services := []v1alpha1.Service{}
			err := json.Unmarshal(b, &services)
			if err != nil {
				return nil, errors.Wrapf(err, "unmarshal services result")
			}
			results := []*resource.ObservabilityMetrics{}
			for _, service := range services {
				if service.Observability != nil && service.Observability.Metrics != nil {
					results = append(results, resource.ToObservabilityMetrics(service.Name, service.Observability.Metrics))
				}
			}
			return results, nil
		})

	return result.([]*resource.ObservabilityMetrics), err
}

type observabilityOutputServerInterface struct {
	client *meshClient
}

func (o *observabilityOutputServerInterface) Get(ctx context.Context, serviceID string) (*resource.ObservabilityOutputServer, error) {
	jsonClient := client.NewHTTPJSON()
	url := fmt.Sprintf("http://"+o.client.server+MeshServiceOutputServerURL, serviceID)
	r, err := jsonClient.
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "get observabilityOutputServer %s", serviceID)
			}

			if statusCode >= 300 {
				return nil, errors.Errorf("call %s failed, return status code: %d text:%s", url, statusCode, string(b))
			}
			output := &v1alpha1.ObservabilityOutputServer{}
			err := json.Unmarshal(b, output)
			if err != nil {
				return nil, errors.Wrap(err, "unmarshal data to v1alpha1.ObservabilityOutputServer")
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
	url := fmt.Sprintf("http://"+o.client.server+MeshServiceOutputServerURL, output.Name())
	update := output.ToV1Alpha1()
	_, err := jsonClient.
		PutByContext(ctx, url, update, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "patch observabilityOutputServer %s", output.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call PUT %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (o *observabilityOutputServerInterface) Create(ctx context.Context, output *resource.ObservabilityOutputServer) error {
	url := fmt.Sprintf("http://"+o.client.server+MeshServiceOutputServerURL, output.Name())
	created := output.ToV1Alpha1()
	_, err := client.NewHTTPJSON().
		PostByContext(ctx, url, created, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusConflict {
				return nil, errors.Wrapf(ConflictError, "create observabilityOutputServer %s", output.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call Post %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (o *observabilityOutputServerInterface) Delete(ctx context.Context, serviceID string) error {
	url := fmt.Sprintf("http://"+o.client.server+MeshServiceOutputServerURL, serviceID)
	_, err := client.NewHTTPJSON().
		DeleteByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "delete observabilityOutputServer %s", serviceID)
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call DELETE %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (o *observabilityOutputServerInterface) List(ctx context.Context) ([]*resource.ObservabilityOutputServer, error) {
	url := fmt.Sprintf("http://" + o.client.server + MeshServicesURL)
	result, err := client.NewHTTPJSON().
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrap(NotFoundError, "list service")
			}

			if statusCode >= 300 || statusCode < 200 {
				return nil, errors.Errorf("call GET %s failed, return statuscode %d text %s", url, statusCode, string(b))
			}

			services := []v1alpha1.Service{}
			err := json.Unmarshal(b, &services)
			if err != nil {
				return nil, errors.Wrapf(err, "unmarshal services result")
			}
			results := []*resource.ObservabilityOutputServer{}
			for _, service := range services {
				if service.Observability != nil && service.Observability.OutputServer != nil {
					results = append(results, resource.ToObservabilityOutputServer(service.Name, service.Observability.OutputServer))
				}
			}
			return results, nil
		})
	if err != nil {
		return nil, err
	}
	return result.([]*resource.ObservabilityOutputServer), err
}
