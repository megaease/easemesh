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

type meshControllerGetter struct {
	client *meshClient
}

func (t *meshControllerGetter) MeshController() MeshControllerInterface {
	return &meshControllerInterface{client: t.client}
}

type meshControllerInterface struct {
	client *meshClient
}

func (t *meshControllerInterface) Get(ctx context.Context, meshControllerID string) (*resource.MeshController, error) {
	url := fmt.Sprintf("http://"+t.client.server+MeshControllerURL, meshControllerID)
	re, err := client.NewHTTPJSON().
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "get meshController %s", meshControllerID)
			}

			if statusCode >= 300 {
				return nil, errors.Errorf("call %s failed: return status code: %d text: %s", url, statusCode, string(b))
			}
			meshController := &resource.MeshControllerV1Alpha1{}
			err := yaml.Unmarshal(b, meshController)
			if err != nil {
				return nil, errors.Wrap(err, "unmarshal data to MeshController")
			}
			return resource.ToMeshController(meshController), nil
		})
	if err != nil {
		return nil, err
	}

	return re.(*resource.MeshController), nil
}

func (t *meshControllerInterface) Patch(ctx context.Context, meshController *resource.MeshController) error {
	jsonClient := client.NewHTTPJSON()
	url := fmt.Sprintf("http://"+t.client.server+MeshControllerURL, meshController.Name())
	update, err := yaml.Marshal(meshController.ToV1Alpha1())
	if err != nil {
		return fmt.Errorf("marshal %#v to yaml failed: %v", meshController, err)
	}

	_, err = jsonClient.
		PutByContext(ctx, url, update, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "patch meshController %s", meshController.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call PUT %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})

	return err
}

func (t *meshControllerInterface) Create(ctx context.Context, meshController *resource.MeshController) error {
	url := fmt.Sprintf("http://" + t.client.server + MeshControllersURL)
	create, err := yaml.Marshal(meshController.ToV1Alpha1())
	if err != nil {
		return fmt.Errorf("marshal %#v to yaml failed: %v", meshController, err)
	}

	_, err = client.NewHTTPJSON().
		PostByContext(ctx, url, create, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusConflict {
				return nil, errors.Wrapf(ConflictError, "create meshController %s", meshController.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call Post %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})

	return err
}

func (t *meshControllerInterface) Delete(ctx context.Context, meshControllerID string) error {
	return errors.New("not support deleting mesh controller, use easegress client to do it")
}

func (t *meshControllerInterface) List(ctx context.Context) ([]*resource.MeshController, error) {
	url := fmt.Sprintf("http://" + t.client.server + MeshControllersURL)
	result, err := client.NewHTTPJSON().
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrap(NotFoundError, "list meshController")
			}

			if statusCode >= 300 || statusCode < 200 {
				return nil, errors.Errorf("call GET %s failed, return statuscode %d text %s", url, statusCode, string(b))
			}

			objects := []map[string]interface{}{}
			err := yaml.Unmarshal(b, &objects)
			if err != nil {
				return nil, errors.Wrapf(err, "unmarshal objects")
			}

			results := []*resource.MeshController{}
			for _, object := range objects {
				if object["kind"] != resource.KindMeshController {
					continue
				}

				buff, err := yaml.Marshal(object)
				if err != nil {
					return nil, errors.Wrapf(err, "marshal %#v to yaml", object)
				}

				meshController := &resource.MeshControllerV1Alpha1{}
				err = yaml.Unmarshal(buff, meshController)
				if err != nil {
					return nil, fmt.Errorf("unmarshal %s to yaml failed: %v", buff, err)
				}

				results = append(results, resource.ToMeshController(meshController))
			}
			return results, nil
		})
	if err != nil {
		return nil, err
	}
	return result.([]*resource.MeshController), err
}
