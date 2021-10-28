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
	"fmt"
	"net/http"

	"github.com/megaease/easemeshctl/cmd/client/resource"
	"github.com/megaease/easemeshctl/cmd/common/client"
	"gopkg.in/yaml.v2"

	"github.com/pkg/errors"
)

type consulServiceRegistryGetter struct {
	client *meshClient
}

func (t *consulServiceRegistryGetter) ConsulServiceRegistry() ConsulServiceRegistryInterface {
	return &consulServiceRegistryInterface{client: t.client}
}

type consulServiceRegistryInterface struct {
	client *meshClient
}

func (t *consulServiceRegistryInterface) Get(ctx context.Context, consulServiceRegistryID string) (*resource.ConsulServiceRegistry, error) {
	url := fmt.Sprintf("http://"+t.client.server+ConsulServiceRegistryURL, consulServiceRegistryID)
	re, err := client.NewHTTPJSON().
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "get consulServiceRegistry %s", consulServiceRegistryID)
			}

			if statusCode >= 300 {
				return nil, errors.Errorf("call %s failed: return status code: %d text: %s", url, statusCode, string(b))
			}
			consulServiceRegistry := &resource.ConsulServiceRegistryV1Alpha1{}
			err := yaml.Unmarshal(b, consulServiceRegistry)
			if err != nil {
				return nil, errors.Wrap(err, "unmarshal data to ConsulServiceRegistry")
			}
			return resource.ToConsulServiceRegistry(consulServiceRegistry), nil
		})
	if err != nil {
		return nil, err
	}

	return re.(*resource.ConsulServiceRegistry), nil
}

func (t *consulServiceRegistryInterface) Patch(ctx context.Context, consulServiceRegistry *resource.ConsulServiceRegistry) error {
	url := fmt.Sprintf("http://"+t.client.server+ConsulServiceRegistryURL, consulServiceRegistry.Name())
	update, err := yaml.Marshal(consulServiceRegistry.ToV1Alpha1())
	if err != nil {
		return fmt.Errorf("marshal %#v to yaml failed: %v", consulServiceRegistry, err)
	}

	_, err = client.NewHTTPJSON().
		PutByContext(ctx, url, update, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "patch consulServiceRegistry %s", consulServiceRegistry.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call PUT %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})

	return err
}

func (t *consulServiceRegistryInterface) Create(ctx context.Context, consulServiceRegistry *resource.ConsulServiceRegistry) error {
	url := fmt.Sprintf("http://" + t.client.server + ConsulServiceRegistrysURL)
	create, err := yaml.Marshal(consulServiceRegistry.ToV1Alpha1())
	if err != nil {
		return fmt.Errorf("marshal %#v to yaml failed: %v", consulServiceRegistry, err)
	}

	_, err = client.NewHTTPJSON().
		PostByContext(ctx, url, create, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusConflict {
				return nil, errors.Wrapf(ConflictError, "create consulServiceRegistry %s", consulServiceRegistry.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call Post %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})

	return err
}

func (t *consulServiceRegistryInterface) Delete(ctx context.Context, consulServiceRegistryID string) error {
	url := fmt.Sprintf("http://"+t.client.server+ConsulServiceRegistryURL, consulServiceRegistryID)
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

func (t *consulServiceRegistryInterface) List(ctx context.Context) ([]*resource.ConsulServiceRegistry, error) {
	url := fmt.Sprintf("http://" + t.client.server + ConsulServiceRegistrysURL)
	result, err := client.NewHTTPJSON().
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrap(NotFoundError, "list consulServiceRegistry")
			}

			if statusCode >= 300 || statusCode < 200 {
				return nil, errors.Errorf("call GET %s failed, return statuscode %d text %s", url, statusCode, string(b))
			}

			objects := []map[string]interface{}{}
			err := yaml.Unmarshal(b, &objects)
			if err != nil {
				return nil, errors.Wrapf(err, "unmarshal objects")
			}

			results := []*resource.ConsulServiceRegistry{}
			for _, object := range objects {
				if object["kind"] != resource.KindConsulServiceRegistry {
					continue
				}

				buff, err := yaml.Marshal(object)
				if err != nil {
					return nil, errors.Wrapf(err, "marshal %#v to yaml", object)
				}

				consulServiceRegistry := &resource.ConsulServiceRegistryV1Alpha1{}
				err = yaml.Unmarshal(buff, consulServiceRegistry)
				if err != nil {
					return nil, fmt.Errorf("unmarshal %s to yaml failed: %v", buff, err)
				}

				results = append(results, resource.ToConsulServiceRegistry(consulServiceRegistry))
			}
			return results, nil
		})
	if err != nil {
		return nil, err
	}
	return result.([]*resource.ConsulServiceRegistry), err
}
