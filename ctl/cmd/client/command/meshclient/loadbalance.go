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

type loadbalanceGetter struct {
	client *meshClient
}

func (s *loadbalanceGetter) LoadBalance() LoadBalanceInterface {
	return &loadbalanceInterface{client: s.client}
}

type loadbalanceInterface struct {
	client *meshClient
}

func (s *loadbalanceInterface) Get(ctx context.Context, serviceID string) (*resource.LoadBalance, error) {
	jsonClient := client.NewHTTPJSON()
	url := fmt.Sprintf("http://"+s.client.server+MeshServiceLoadBalanceURL, serviceID)
	r, err := jsonClient.
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "get loadbalance %s", serviceID)
			}
			if statusCode >= 300 {
				return nil, errors.Errorf("call %s failed, return status code: %d text:%s", url, statusCode, string(b))
			}
			loadbalance := &v1alpha1.LoadBalance{}
			err := json.Unmarshal(b, loadbalance)
			if err != nil {
				return nil, errors.Wrap(err, "unmarshal data to v1alpha1.LoadBalance")
			}
			return resource.ToLoadbalance(serviceID, loadbalance), nil
		})
	if err != nil {
		return nil, err
	}

	return r.(*resource.LoadBalance), nil
}

func (s *loadbalanceInterface) Patch(ctx context.Context, loadbalance *resource.LoadBalance) error {
	jsonClient := client.NewHTTPJSON()
	url := fmt.Sprintf("http://"+s.client.server+MeshServiceLoadBalanceURL, loadbalance.Name())
	update := loadbalance.ToV1Alpha1()
	_, err := jsonClient.
		PutByContext(ctx, url, update, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "patch loadbalance %s", loadbalance.Name())
			}
			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call PUT %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (s *loadbalanceInterface) Create(ctx context.Context, loadbalance *resource.LoadBalance) error {
	url := fmt.Sprintf("http://"+s.client.server+MeshServiceLoadBalanceURL, loadbalance.Name())
	created := loadbalance.ToV1Alpha1()
	_, err := client.NewHTTPJSON().
		PostByContext(ctx, url, created, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusConflict {
				return nil, errors.Wrapf(ConflictError, "create loadbalance %s", loadbalance.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call Post %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (s *loadbalanceInterface) Delete(ctx context.Context, serviceID string) error {
	url := fmt.Sprintf("http://"+s.client.server+MeshServiceLoadBalanceURL, serviceID)
	_, err := client.NewHTTPJSON().
		DeleteByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "delete loadbalance %s", serviceID)
			}
			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call DELETE %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (s *loadbalanceInterface) List(ctx context.Context) ([]*resource.LoadBalance, error) {
	url := "http://" + s.client.server + MeshServicesURL
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
			results := []*resource.LoadBalance{}
			for _, ss := range services {
				if ss.LoadBalance != nil {
					results = append(results, resource.ToLoadbalance(ss.Name, ss.LoadBalance))
				}
			}
			return results, nil
		})
	if err != nil {
		return nil, err
	}
	return result.([]*resource.LoadBalance), err
}
