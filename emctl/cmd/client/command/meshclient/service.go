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

type serviceGetter struct {
	client *meshClient
}

func (s *serviceGetter) Service() ServiceInterface {
	return &serviceInterface{client: s.client}
}

var _ ServiceInterface = &serviceInterface{}

type serviceInterface struct {
	client *meshClient
}

func (s *serviceInterface) Get(ctx context.Context, serviceID string) (*resource.Service, error) {
	jsonClient := client.NewHTTPJSON()
	url := fmt.Sprintf("http://"+s.client.server+MeshServiceURL, serviceID)
	r, err := jsonClient.
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "get service %s not found", serviceID)
			}
			if statusCode >= 300 {
				return nil, errors.Errorf("call %s failed, return status code: %d text:%s", url, statusCode, string(b))
			}
			service := &v1alpha1.Service{}
			err := json.Unmarshal(b, service)
			if err != nil {
				return nil, errors.Wrap(err, "unmarshal data to v1alpha1.Service")
			}
			return resource.ToService(service), nil
		})
	if err != nil {
		return nil, err
	}

	return r.(*resource.Service), nil
}

func (s *serviceInterface) Patch(ctx context.Context, service *resource.Service) error {
	jsonClient := client.NewHTTPJSON()
	update := service.ToV1Alpha1()
	url := fmt.Sprintf("http://"+s.client.server+MeshServiceURL, service.Name())
	_, err := jsonClient.
		PutByContext(ctx, url, update, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, NotFoundError
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call PUT %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (s *serviceInterface) Create(ctx context.Context, service *resource.Service) error {
	created := service.ToV1Alpha1()
	url := fmt.Sprintf("http://"+s.client.server+MeshServicesURL)
	_, err := client.NewHTTPJSON().
		PostByContext(ctx, url, created, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusConflict {
				return nil, errors.Wrapf(ConflictError, "create service %s", service.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call Post %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (s *serviceInterface) Delete(ctx context.Context, serviceID string) error {
	url := fmt.Sprintf("http://"+s.client.server+MeshServiceURL, serviceID)
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

func (s *serviceInterface) List(ctx context.Context) ([]*resource.Service, error) {
	url := fmt.Sprintf("http://" + s.client.server + MeshServicesURL)
	result, err := client.NewHTTPJSON().
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, NotFoundError
			}

			if statusCode >= 300 || statusCode < 200 {
				return nil, errors.Errorf("call GET %s failed, return statuscode %d text %s", url, statusCode, string(b))
			}

			services := []v1alpha1.Service{}
			err := json.Unmarshal(b, &services)
			if err != nil {
				return nil, errors.Wrapf(err, "unmarshal services result")
			}
			results := []*resource.Service{}
			for _, service := range services {
				copy := service
				results = append(results, resource.ToService(&copy))
			}
			return results, nil
		})
	if err != nil {
		return nil, err
	}
	return result.([]*resource.Service), err
}
