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

type resilienceGetter struct {
	client *meshClient
}

func (s *resilienceGetter) Resilience() ResilienceInterface {
	return &resilienceInterface{client: s.client}
}

type resilienceInterface struct {
	client *meshClient
}

func (r *resilienceInterface) Get(ctx context.Context, serviceID string) (*resource.Resilience, error) {
	url := fmt.Sprintf("http://"+r.client.server+MeshServiceResilienceURL, serviceID)
	re, err := client.NewHTTPJSON().
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "get resilience %s", serviceID)
			}

			if statusCode >= 300 {
				return nil, errors.Errorf("call %s failed, return status code %d text %s", url, statusCode, string(b))
			}
			resilience := &v1alpha1.Resilience{}
			err := json.Unmarshal(b, resilience)
			if err != nil {
				return nil, errors.Wrap(err, "unmarshal data to v1alpha1.Resilience")
			}
			return resource.ToResilience(serviceID, resilience), nil
		})
	if err != nil {
		return nil, err
	}

	return re.(*resource.Resilience), nil
}

func (r *resilienceInterface) Patch(ctx context.Context, resilience *resource.Resilience) error {
	jsonClient := client.NewHTTPJSON()
	url := fmt.Sprintf("http://"+r.client.server+MeshServiceResilienceURL, resilience.Name())
	update := resilience.ToV1Alpha1()
	_, err := jsonClient.
		PutByContext(ctx, url, update, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "patch resilience %s", resilience.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call PUT %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (r *resilienceInterface) Create(ctx context.Context, resilience *resource.Resilience) error {
	url := fmt.Sprintf("http://"+r.client.server+MeshServiceResilienceURL, resilience.Name())
	created := resilience.ToV1Alpha1()
	_, err := client.NewHTTPJSON().
		PostByContext(ctx, url, created, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusConflict {
				return nil, errors.Wrapf(ConflictError, "create resilience %s", resilience.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call Post %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (r *resilienceInterface) Delete(ctx context.Context, serviceID string) error {
	url := fmt.Sprintf("http://"+r.client.server+MeshServiceResilienceURL, serviceID)
	_, err := client.NewHTTPJSON().
		DeleteByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "delete resilience %s", serviceID)
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call DELETE %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (r *resilienceInterface) List(ctx context.Context) ([]*resource.Resilience, error) {
	url := fmt.Sprintf("http://" + r.client.server + MeshServicesURL)
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
			results := []*resource.Resilience{}
			for _, ss := range services {
				if ss.Resilience != nil {
					results = append(results, resource.ToResilience(ss.Name, ss.Resilience))
				}
			}
			return results, nil
		})
	if err != nil {
		return nil, err
	}
	return result.([]*resource.Resilience), err
}
