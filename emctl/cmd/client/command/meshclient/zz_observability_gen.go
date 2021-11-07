/*
Copyright (c) 2021, MegaEase
All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// code generated by github.com/megaease/easemeshctl/cmd/generator, DO NOT EDIT.
package meshclient

import (
	"context"
	"encoding/json"
	"fmt"
	v1alpha1 "github.com/megaease/easemesh-api/v1alpha1"
	resource "github.com/megaease/easemeshctl/cmd/client/resource"
	client "github.com/megaease/easemeshctl/cmd/common/client"
	errors "github.com/pkg/errors"
	"net/http"
)

type observabilityGetter struct {
	client *meshClient
}
type observabilityOutputServerInterface struct {
	client *meshClient
}
type observabilityMetricsInterface struct {
	client *meshClient
}
type observabilityTracingsInterface struct {
	client *meshClient
}

func (o *observabilityGetter) ObservabilityTracings() ObservabilityTracingsInterface {
	return &observabilityTracingsInterface{client: o.client}
}
func (o *observabilityGetter) ObservabilityMetrics() ObservabilityMetricsInterface {
	return &observabilityMetricsInterface{client: o.client}
}
func (o *observabilityGetter) ObservabilityOutputServer() ObservabilityOutputServerInterface {
	return &observabilityOutputServerInterface{client: o.client}
}
func (o *observabilityOutputServerInterface) Get(args0 context.Context, args1 string) (*resource.ObservabilityOutputServer, error) {
	url := fmt.Sprintf("http://"+o.client.server+apiURL+"/mesh/"+"services/%s/outputserver", args1)
	r0, err := client.NewHTTPJSON().GetByContext(args0, url, nil, nil).HandleResponse(func(buff []byte, statusCode int) (interface{}, error) {
		if statusCode == http.StatusNotFound {
			return nil, errors.Wrapf(NotFoundError, "get ObservabilityOutputServer %s", args1)
		}
		if statusCode >= 300 {
			return nil, errors.Errorf("call %s failed, return status code %d text %+v", url, statusCode, string(buff))
		}
		ObservabilityOutputServer := &v1alpha1.ObservabilityOutputServer{}
		err := json.Unmarshal(buff, ObservabilityOutputServer)
		if err != nil {
			return nil, errors.Wrapf(err, "unmarshal data to v1alpha1.ObservabilityOutputServer")
		}
		return resource.ToObservabilityOutputServer(args1, ObservabilityOutputServer), nil
	})
	if err != nil {
		return nil, err
	}
	return r0.(*resource.ObservabilityOutputServer), nil
}
func (o *observabilityOutputServerInterface) Patch(args0 context.Context, args1 *resource.ObservabilityOutputServer) error {
	url := fmt.Sprintf("http://"+o.client.server+apiURL+"/mesh/"+"services/%s/outputserver", args1.Name())
	object := args1.ToV1Alpha1()
	_, err := client.NewHTTPJSON().PutByContext(args0, url, object, nil).HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
		if statusCode == http.StatusNotFound {
			return nil, errors.Wrapf(NotFoundError, "patch ObservabilityOutputServer %s", args1.Name())
		}
		if statusCode < 300 && statusCode >= 200 {
			return nil, nil
		}
		return nil, errors.Errorf("call PUT %s failed, return statuscode %d text %+v", url, statusCode, string(b))
	})
	return err
}
func (o *observabilityOutputServerInterface) Create(args0 context.Context, args1 *resource.ObservabilityOutputServer) error {
	url := fmt.Sprintf("http://"+o.client.server+apiURL+"/mesh/"+"services/%s/outputserver", args1.Name())
	object := args1.ToV1Alpha1()
	_, err := client.NewHTTPJSON().PostByContext(args0, url, object, nil).HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
		if statusCode == http.StatusConflict {
			return nil, errors.Wrapf(ConflictError, "create ObservabilityOutputServer %s", args1.Name())
		}
		if statusCode < 300 && statusCode >= 200 {
			return nil, nil
		}
		return nil, errors.Errorf("call Post %s failed, return statuscode %d text %+v", url, statusCode, string(b))
	})
	return err
}
func (o *observabilityOutputServerInterface) Delete(args0 context.Context, args1 string) error {
	url := fmt.Sprintf("http://"+o.client.server+apiURL+"/mesh/"+"services/%s/outputserver", args1)
	_, err := client.NewHTTPJSON().DeleteByContext(args0, url, nil, nil).HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
		if statusCode == http.StatusNotFound {
			return nil, errors.Wrapf(NotFoundError, "Delete ObservabilityOutputServer %s", args1)
		}
		if statusCode < 300 && statusCode >= 200 {
			return nil, nil
		}
		return nil, errors.Errorf("call Delete %s failed, return statuscode %d text %+v", url, statusCode, string(b))
	})
	return err
}
func (o *observabilityOutputServerInterface) List(args0 context.Context) ([]*resource.ObservabilityOutputServer, error) {
	url := "http://" + o.client.server + apiURL + "/mesh/services"
	result, err := client.NewHTTPJSON().GetByContext(args0, url, nil, nil).HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
		if statusCode == http.StatusNotFound {
			return nil, errors.Wrapf(NotFoundError, "list service")
		}
		if statusCode >= 300 && statusCode < 200 {
			return nil, errors.Errorf("call GET %s failed, return statuscode %d text %+v", url, statusCode, b)
		}
		services := []v1alpha1.Service{}
		err := json.Unmarshal(b, &services)
		if err != nil {
			return nil, errors.Wrapf(err, "unmarshal data to v1alpha1.")
		}
		results := []*resource.ObservabilityOutputServer{}
		for _, service := range services {
			if service.Observability != nil {
				results = append(results, resource.ToObservabilityOutputServer(service.Name, service.Observability.OutputServer))
			}
		}
		return results, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*resource.ObservabilityOutputServer), nil
}
func (o *observabilityMetricsInterface) Get(args0 context.Context, args1 string) (*resource.ObservabilityMetrics, error) {
	url := fmt.Sprintf("http://"+o.client.server+apiURL+"/mesh/"+"services/%s/metrics", args1)
	r0, err := client.NewHTTPJSON().GetByContext(args0, url, nil, nil).HandleResponse(func(buff []byte, statusCode int) (interface{}, error) {
		if statusCode == http.StatusNotFound {
			return nil, errors.Wrapf(NotFoundError, "get ObservabilityMetrics %s", args1)
		}
		if statusCode >= 300 {
			return nil, errors.Errorf("call %s failed, return status code %d text %+v", url, statusCode, string(buff))
		}
		ObservabilityMetrics := &v1alpha1.ObservabilityMetrics{}
		err := json.Unmarshal(buff, ObservabilityMetrics)
		if err != nil {
			return nil, errors.Wrapf(err, "unmarshal data to v1alpha1.ObservabilityMetrics")
		}
		return resource.ToObservabilityMetrics(args1, ObservabilityMetrics), nil
	})
	if err != nil {
		return nil, err
	}
	return r0.(*resource.ObservabilityMetrics), nil
}
func (o *observabilityMetricsInterface) Patch(args0 context.Context, args1 *resource.ObservabilityMetrics) error {
	url := fmt.Sprintf("http://"+o.client.server+apiURL+"/mesh/"+"services/%s/metrics", args1.Name())
	object := args1.ToV1Alpha1()
	_, err := client.NewHTTPJSON().PutByContext(args0, url, object, nil).HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
		if statusCode == http.StatusNotFound {
			return nil, errors.Wrapf(NotFoundError, "patch ObservabilityMetrics %s", args1.Name())
		}
		if statusCode < 300 && statusCode >= 200 {
			return nil, nil
		}
		return nil, errors.Errorf("call PUT %s failed, return statuscode %d text %+v", url, statusCode, string(b))
	})
	return err
}
func (o *observabilityMetricsInterface) Create(args0 context.Context, args1 *resource.ObservabilityMetrics) error {
	url := fmt.Sprintf("http://"+o.client.server+apiURL+"/mesh/"+"services/%s/metrics", args1.Name())
	object := args1.ToV1Alpha1()
	_, err := client.NewHTTPJSON().PostByContext(args0, url, object, nil).HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
		if statusCode == http.StatusConflict {
			return nil, errors.Wrapf(ConflictError, "create ObservabilityMetrics %s", args1.Name())
		}
		if statusCode < 300 && statusCode >= 200 {
			return nil, nil
		}
		return nil, errors.Errorf("call Post %s failed, return statuscode %d text %+v", url, statusCode, string(b))
	})
	return err
}
func (o *observabilityMetricsInterface) Delete(args0 context.Context, args1 string) error {
	url := fmt.Sprintf("http://"+o.client.server+apiURL+"/mesh/"+"services/%s/metrics", args1)
	_, err := client.NewHTTPJSON().DeleteByContext(args0, url, nil, nil).HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
		if statusCode == http.StatusNotFound {
			return nil, errors.Wrapf(NotFoundError, "Delete ObservabilityMetrics %s", args1)
		}
		if statusCode < 300 && statusCode >= 200 {
			return nil, nil
		}
		return nil, errors.Errorf("call Delete %s failed, return statuscode %d text %+v", url, statusCode, string(b))
	})
	return err
}
func (o *observabilityMetricsInterface) List(args0 context.Context) ([]*resource.ObservabilityMetrics, error) {
	url := "http://" + o.client.server + apiURL + "/mesh/services"
	result, err := client.NewHTTPJSON().GetByContext(args0, url, nil, nil).HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
		if statusCode == http.StatusNotFound {
			return nil, errors.Wrapf(NotFoundError, "list service")
		}
		if statusCode >= 300 && statusCode < 200 {
			return nil, errors.Errorf("call GET %s failed, return statuscode %d text %+v", url, statusCode, b)
		}
		services := []v1alpha1.Service{}
		err := json.Unmarshal(b, &services)
		if err != nil {
			return nil, errors.Wrapf(err, "unmarshal data to v1alpha1.")
		}
		results := []*resource.ObservabilityMetrics{}
		for _, service := range services {
			if service.Observability != nil {
				results = append(results, resource.ToObservabilityMetrics(service.Name, service.Observability.Metrics))
			}
		}
		return results, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*resource.ObservabilityMetrics), nil
}
func (o *observabilityTracingsInterface) Get(args0 context.Context, args1 string) (*resource.ObservabilityTracings, error) {
	url := fmt.Sprintf("http://"+o.client.server+apiURL+"/mesh/"+"services/%s/tracings", args1)
	r0, err := client.NewHTTPJSON().GetByContext(args0, url, nil, nil).HandleResponse(func(buff []byte, statusCode int) (interface{}, error) {
		if statusCode == http.StatusNotFound {
			return nil, errors.Wrapf(NotFoundError, "get ObservabilityTracings %s", args1)
		}
		if statusCode >= 300 {
			return nil, errors.Errorf("call %s failed, return status code %d text %+v", url, statusCode, string(buff))
		}
		ObservabilityTracings := &v1alpha1.ObservabilityTracings{}
		err := json.Unmarshal(buff, ObservabilityTracings)
		if err != nil {
			return nil, errors.Wrapf(err, "unmarshal data to v1alpha1.ObservabilityTracings")
		}
		return resource.ToObservabilityTracings(args1, ObservabilityTracings), nil
	})
	if err != nil {
		return nil, err
	}
	return r0.(*resource.ObservabilityTracings), nil
}
func (o *observabilityTracingsInterface) Patch(args0 context.Context, args1 *resource.ObservabilityTracings) error {
	url := fmt.Sprintf("http://"+o.client.server+apiURL+"/mesh/"+"services/%s/tracings", args1.Name())
	object := args1.ToV1Alpha1()
	_, err := client.NewHTTPJSON().PutByContext(args0, url, object, nil).HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
		if statusCode == http.StatusNotFound {
			return nil, errors.Wrapf(NotFoundError, "patch ObservabilityTracings %s", args1.Name())
		}
		if statusCode < 300 && statusCode >= 200 {
			return nil, nil
		}
		return nil, errors.Errorf("call PUT %s failed, return statuscode %d text %+v", url, statusCode, string(b))
	})
	return err
}
func (o *observabilityTracingsInterface) Create(args0 context.Context, args1 *resource.ObservabilityTracings) error {
	url := fmt.Sprintf("http://"+o.client.server+apiURL+"/mesh/"+"services/%s/tracings", args1.Name())
	object := args1.ToV1Alpha1()
	_, err := client.NewHTTPJSON().PostByContext(args0, url, object, nil).HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
		if statusCode == http.StatusConflict {
			return nil, errors.Wrapf(ConflictError, "create ObservabilityTracings %s", args1.Name())
		}
		if statusCode < 300 && statusCode >= 200 {
			return nil, nil
		}
		return nil, errors.Errorf("call Post %s failed, return statuscode %d text %+v", url, statusCode, string(b))
	})
	return err
}
func (o *observabilityTracingsInterface) Delete(args0 context.Context, args1 string) error {
	url := fmt.Sprintf("http://"+o.client.server+apiURL+"/mesh/"+"services/%s/tracings", args1)
	_, err := client.NewHTTPJSON().DeleteByContext(args0, url, nil, nil).HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
		if statusCode == http.StatusNotFound {
			return nil, errors.Wrapf(NotFoundError, "Delete ObservabilityTracings %s", args1)
		}
		if statusCode < 300 && statusCode >= 200 {
			return nil, nil
		}
		return nil, errors.Errorf("call Delete %s failed, return statuscode %d text %+v", url, statusCode, string(b))
	})
	return err
}
func (o *observabilityTracingsInterface) List(args0 context.Context) ([]*resource.ObservabilityTracings, error) {
	url := "http://" + o.client.server + apiURL + "/mesh/services"
	result, err := client.NewHTTPJSON().GetByContext(args0, url, nil, nil).HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
		if statusCode == http.StatusNotFound {
			return nil, errors.Wrapf(NotFoundError, "list service")
		}
		if statusCode >= 300 && statusCode < 200 {
			return nil, errors.Errorf("call GET %s failed, return statuscode %d text %+v", url, statusCode, b)
		}
		services := []v1alpha1.Service{}
		err := json.Unmarshal(b, &services)
		if err != nil {
			return nil, errors.Wrapf(err, "unmarshal data to v1alpha1.")
		}
		results := []*resource.ObservabilityTracings{}
		for _, service := range services {
			if service.Observability != nil {
				results = append(results, resource.ToObservabilityTracings(service.Name, service.Observability.Tracings))
			}
		}
		return results, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*resource.ObservabilityTracings), nil
}
