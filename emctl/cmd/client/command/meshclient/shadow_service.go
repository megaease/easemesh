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
	"k8s.io/kubernetes/test/e2e/storage/drivers/csi-test/mock/service"

	"github.com/megaease/easemesh-api/v1alpha1"
	"github.com/pkg/errors"
)

type shadowServiceGetter struct {
	client *meshClient
}

func (s *shadowServiceGetter) Service() ShadowServiceInterface {
	return &shadowServiceInterface{client: s.client}
}

var _ ShadowServiceInterface = &shadowServiceInterface{}

type shadowServiceInterface struct {
	client *meshClient
}

func (s *shadowServiceInterface) Get(ctx context.Context, serviceID string) (*resource.ShadowService, error) {
	jsonClient := client.NewHTTPJSON()
	url := fmt.Sprintf("http://"+s.client.server+MeshShadowServiceURL, serviceID)
	r, err := jsonClient.
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "get shadow service %s not found", serviceID)
			}
			if statusCode >= 300 {
				return nil, errors.Errorf("call %s failed, return status code: %d text:%s", url, statusCode, string(b))
			}
			shadowService := &v1alpha1.ShadowService{}
			err := json.Unmarshal(b, shadowService)
			if err != nil {
				return nil, errors.Wrap(err, "unmarshal data to v1alpha1.ShadowService")
			}
			return resource.ToShadowService(shadowService), nil
		})
	if err != nil {
		return nil, err
	}

	return r.(*resource.ShadowService), nil
}

func (s *shadowServiceInterface) Patch(ctx context.Context, shadowService *resource.ShadowService) error {
	jsonClient := client.NewHTTPJSON()
	update := shadowService.ToV1Alpha1()
	url := fmt.Sprintf("http://"+s.client.server+MeshShadowServiceURL, shadowService.Name())
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

func (s *shadowServiceInterface) Create(ctx context.Context, shadowService *resource.ShadowService) error {
	created := shadowService.ToV1Alpha1()
	url := fmt.Sprintf("http://"+s.client.server+MeshShadowServiceURL, shadowService.Name())
	_, err := client.NewHTTPJSON().
		PostByContext(ctx, url, created, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusConflict {
				return nil, errors.Wrapf(ConflictError, "create shadow service %s", shadowService.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call Post %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (s *shadowServiceInterface) Delete(ctx context.Context, shadowServiceID string) error {
	url := fmt.Sprintf("http://"+s.client.server+MeshShadowServiceURL, shadowServiceID)
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

func (s *shadowServiceInterface) List(ctx context.Context) ([]*resource.ShadowService, error) {
	url := fmt.Sprintf("http://" + s.client.server + MeshShadowServicesURL)
	result, err := client.NewHTTPJSON().
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, NotFoundError
			}

			if statusCode >= 300 || statusCode < 200 {
				return nil, errors.Errorf("call GET %s failed, return statuscode %d text %s", url, statusCode, string(b))
			}

			shadowServices := []v1alpha1.ShadowService{}
			err := json.Unmarshal(b, &shadowServices)
			if err != nil {
				return nil, errors.Wrapf(err, "unmarshal services result")
			}
			results := []*resource.ShadowService{}
			for _, ss := range shadowServices {
				results = append(results, resource.ToShadowService(&ss))
			}
			return results, nil
		})
	if err != nil {
		return nil, err
	}
	return result.([]*resource.ShadowService), err
}
