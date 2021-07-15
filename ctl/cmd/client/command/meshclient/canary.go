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

type canaryGetter struct {
	client *meshClient
}

func (c *canaryGetter) Canary() CanaryInterface {
	return &canaryInterface{client: c.client}
}

var _ CanaryInterface = &canaryInterface{}

type canaryInterface struct {
	client *meshClient
}

func (c *canaryInterface) Get(ctx context.Context, serviceID string) (*resource.Canary, error) {
	jsonClient := client.NewHTTPJSON()
	url := fmt.Sprintf("http://"+c.client.server+MeshServiceCanaryURL, serviceID)
	r, err := jsonClient.
		GetByContext(url, nil, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "get canary %s", serviceID)
			}

			if statusCode >= 300 {
				return nil, errors.Errorf("call %s failed, return status code %d text %s", url, statusCode, string(b))
			}
			canary := &v1alpha1.Canary{}
			err := json.Unmarshal(b, canary)
			if err != nil {
				return nil, errors.Wrap(err, "unmarshal data to v1alpha1.Canary")
			}
			return resource.ToCanary(serviceID, canary), nil
		})
	if err != nil {
		return nil, err
	}

	return r.(*resource.Canary), nil
}

func (c *canaryInterface) Patch(ctx context.Context, canary *resource.Canary) error {
	jsonClient := client.NewHTTPJSON()
	url := fmt.Sprintf("http://"+c.client.server+MeshServiceCanaryURL, canary.Name())
	update := canary.ToV1Alpha1()
	_, err := jsonClient.
		PutByContext(url, update, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "patch canary %s", canary.Name())
			}
			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call PUT %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (c *canaryInterface) Create(ctx context.Context, canary *resource.Canary) error {
	url := fmt.Sprintf("http://"+c.client.server+MeshServiceCanaryURL, canary.Name())
	created := canary.ToV1Alpha1()
	_, err := client.NewHTTPJSON().
		PostByContext(url, created, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusConflict {
				return nil, errors.Wrapf(ConflictError, "create canary %s", canary.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call Post %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (c *canaryInterface) Delete(ctx context.Context, serviceID string) error {
	url := fmt.Sprintf("http://"+c.client.server+MeshServiceCanaryURL, serviceID)
	_, err := client.NewHTTPJSON().
		DeleteByContext(url, nil, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "delete canary %s", serviceID)
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call DELETE %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (c *canaryInterface) List(ctx context.Context) ([]*resource.Canary, error) {
	url := "http://" + c.client.server + MeshServicesURL
	result, err := client.NewHTTPJSON().
		GetByContext(url, nil, ctx, nil).
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
			results := []*resource.Canary{}
			for _, ss := range services {
				if ss.Canary != nil {
					results = append(results, resource.ToCanary(ss.Name, ss.Canary))
				}
			}
			return results, nil
		})
	if err != nil {
		return nil, err
	}
	return result.([]*resource.Canary), err
}
