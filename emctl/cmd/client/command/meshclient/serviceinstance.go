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

	"github.com/megaease/easemeshctl/cmd/client/resource"
	"github.com/megaease/easemeshctl/cmd/common/client"

	"github.com/megaease/easemesh-api/v1alpha1"
	"github.com/pkg/errors"
)

type serviceInstanceGetter struct {
	client *meshClient
}

func (s *serviceInstanceGetter) ServiceInstance() ServiceInstanceInterface {
	return &serviceInstanceInterface{client: s.client}
}

var _ ServiceInstanceInterface = &serviceInstanceInterface{}

type serviceInstanceInterface struct {
	client *meshClient
}

func (s *serviceInstanceInterface) Get(ctx context.Context, serviceName, instanceID string) (*resource.ServiceInstance, error) {
	jsonClient := client.NewHTTPJSON()
	url := fmt.Sprintf("http://"+s.client.server+MeshServiceInstanceURL, serviceName, instanceID)
	r, err := jsonClient.
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "get serviceInstance %s/%s not found",
					serviceName, instanceID)
			}
			if statusCode >= 300 {
				return nil, errors.Errorf("call %s failed, return status code: %d text:%s", url, statusCode, string(b))
			}
			serviceInstance := &v1alpha1.ServiceInstance{}
			err := json.Unmarshal(b, serviceInstance)
			if err != nil {
				return nil, errors.Wrap(err, "unmarshal data to v1alpha1.ServiceInstance")
			}
			return resource.ToServiceInstance(serviceInstance), nil
		})
	if err != nil {
		return nil, err
	}

	return r.(*resource.ServiceInstance), nil
}

func (s *serviceInstanceInterface) Delete(ctx context.Context, serviceName, instanceID string) error {
	url := fmt.Sprintf("http://"+s.client.server+MeshServiceInstanceURL, serviceName, instanceID)
	_, err := client.NewHTTPJSON().
		DeleteByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, NotFoundError
			}
			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call DELETE %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (s *serviceInstanceInterface) List(ctx context.Context) ([]*resource.ServiceInstance, error) {
	url := fmt.Sprintf("http://" + s.client.server + MeshServiceInstancesURL)
	result, err := client.NewHTTPJSON().
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, NotFoundError
			}

			if statusCode >= 300 || statusCode < 200 {
				return nil, errors.Errorf("call GET %s failed, return statuscode %d text %s", url, statusCode, string(b))
			}

			serviceInstances := []v1alpha1.ServiceInstance{}
			err := json.Unmarshal(b, &serviceInstances)
			if err != nil {
				return nil, errors.Wrapf(err, "unmarshal serviceInstances result")
			}
			results := []*resource.ServiceInstance{}
			for _, ss := range serviceInstances {
				results = append(results, resource.ToServiceInstance(&ss))
			}
			return results, nil
		})
	if err != nil {
		return nil, err
	}
	return result.([]*resource.ServiceInstance), err
}
