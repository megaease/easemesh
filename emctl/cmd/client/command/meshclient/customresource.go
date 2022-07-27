/*
 * Copyright (c) 2021, MegaEase
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

	"github.com/megaease/easemesh-api/v2alpha1"
	"github.com/megaease/easemeshctl/cmd/client/resource"
	"github.com/megaease/easemeshctl/cmd/common/client"

	"github.com/pkg/errors"
)

type customResourceKindGetter struct {
	client *meshClient
}

func (t *customResourceKindGetter) CustomResourceKind() CustomResourceKindInterface {
	return &customResourceKindInterface{client: t.client}
}

type customResourceKindInterface struct {
	client *meshClient
}

func (k *customResourceKindInterface) Get(ctx context.Context, customResourceKindID string) (*resource.CustomResourceKind, error) {
	url := fmt.Sprintf("http://"+k.client.server+MeshCustomResourceKindURL, customResourceKindID)
	re, err := client.NewHTTPJSON().
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "get custom resource kind %s", customResourceKindID)
			}

			if statusCode >= 300 {
				return nil, errors.Errorf("call %s failed, return status code: %d text:%s", url, statusCode, string(b))
			}
			customResourceKind := &v2alpha1.CustomResourceKind{}
			err := json.Unmarshal(b, customResourceKind)
			if err != nil {
				return nil, errors.Wrap(err, "unmarshal data to v2alpha1.CustomResourceKind")
			}
			return resource.ToCustomResourceKind(customResourceKind), nil
		})
	if err != nil {
		return nil, err
	}

	return re.(*resource.CustomResourceKind), nil
}

func (k *customResourceKindInterface) Patch(ctx context.Context, customResourceKind *resource.CustomResourceKind) error {
	jsonClient := client.NewHTTPJSON()
	url := "http://" + k.client.server + MeshCustomResourceKindsURL
	update := customResourceKind.ToV2Alpha1()
	_, err := jsonClient.
		PutByContext(ctx, url, update, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "patch custom resource kind %s", customResourceKind.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call PUT %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (k *customResourceKindInterface) Create(ctx context.Context, customResourceKind *resource.CustomResourceKind) error {
	created := customResourceKind.ToV2Alpha1()
	url := "http://" + k.client.server + MeshCustomResourceKindsURL
	_, err := client.NewHTTPJSON().
		// FIXME: the standard RESTful URL of create resource is POST /v1/api/{resources} instead of POST /v1/api/{resources}/{id}.
		// Current URL form should be corrected in the feature
		PostByContext(ctx, url, created, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusConflict {
				return nil, errors.Wrapf(ConflictError, "create custom resource kind %s", customResourceKind.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call Post %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (k *customResourceKindInterface) Delete(ctx context.Context, customResourceKindID string) error {
	url := fmt.Sprintf("http://"+k.client.server+MeshCustomResourceKindURL, customResourceKindID)
	_, err := client.NewHTTPJSON().
		DeleteByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "delete custom resource kind %s", customResourceKindID)
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call DELETE %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (k *customResourceKindInterface) List(ctx context.Context) ([]*resource.CustomResourceKind, error) {
	url := "http://" + k.client.server + MeshCustomResourceKindsURL
	result, err := client.NewHTTPJSON().
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrap(NotFoundError, "list custom resource kind")
			}

			if statusCode >= 300 || statusCode < 200 {
				return nil, errors.Errorf("call GET %s failed, return statuscode %d text %s", url, statusCode, string(b))
			}

			customResourceKinds := []v2alpha1.CustomResourceKind{}
			err := json.Unmarshal(b, &customResourceKinds)
			if err != nil {
				return nil, errors.Wrapf(err, "unmarshal custom resource kind result")
			}

			results := []*resource.CustomResourceKind{}
			for _, customResourceKind := range customResourceKinds {
				copy := customResourceKind
				results = append(results, resource.ToCustomResourceKind(&copy))
			}
			return results, nil
		})
	if err != nil {
		return nil, err
	}
	return result.([]*resource.CustomResourceKind), err
}

type customResourceGetter struct {
	client *meshClient
}

func (t *customResourceGetter) CustomResource() CustomResourceInterface {
	return &customResourceInterface{client: t.client}
}

type customResourceInterface struct {
	client *meshClient
}

func (o *customResourceInterface) Get(ctx context.Context, kind, customResourceID string) (*resource.CustomResource, error) {
	url := fmt.Sprintf("http://"+o.client.server+MeshCustomResourceURL, kind, customResourceID)
	re, err := client.NewHTTPJSON().
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "get custom resource %s", customResourceID)
			}

			if statusCode >= 300 {
				return nil, errors.Errorf("call %s failed, return status code: %d text:%s", url, statusCode, string(b))
			}
			customResource := map[string]interface{}{}
			err := json.Unmarshal(b, &customResource)
			if err != nil {
				return nil, errors.Wrap(err, "unmarshal data to v2alpha1.CustomResource")
			}
			return resource.ToCustomResource(customResource), nil
		})
	if err != nil {
		return nil, err
	}

	return re.(*resource.CustomResource), nil
}

func (o *customResourceInterface) Patch(ctx context.Context, customResource *resource.CustomResource) error {
	jsonClient := client.NewHTTPJSON()
	url := "http://" + o.client.server + MeshAllCustomResourcesURL
	update := customResource.ToV2Alpha1()
	_, err := jsonClient.
		PutByContext(ctx, url, update, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "patch custom resource %s", customResource.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call PUT %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (o *customResourceInterface) Create(ctx context.Context, customResource *resource.CustomResource) error {
	created := customResource.ToV2Alpha1()
	url := "http://" + o.client.server + MeshAllCustomResourcesURL
	_, err := client.NewHTTPJSON().
		// FIXME: the standard RESTful URL of create resource is POST /v1/api/{resources} instead of POST /v1/api/{resources}/{id}.
		// Current URL form should be corrected in the feature
		PostByContext(ctx, url, created, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusConflict {
				return nil, errors.Wrapf(ConflictError, "create custom resource %s", customResource.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call Post %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (o *customResourceInterface) Delete(ctx context.Context, kind, customResourceID string) error {
	url := fmt.Sprintf("http://"+o.client.server+MeshCustomResourceURL, kind, customResourceID)
	_, err := client.NewHTTPJSON().
		DeleteByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "delete custom resource %s", customResourceID)
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call DELETE %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (o *customResourceInterface) List(ctx context.Context, kind string) ([]*resource.CustomResource, error) {
	url := fmt.Sprintf("http://"+o.client.server+MeshCustomResourcesURL, kind)
	result, err := client.NewHTTPJSON().
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrap(NotFoundError, "list custom resource")
			}

			if statusCode >= 300 || statusCode < 200 {
				return nil, errors.Errorf("call GET %s failed, return statuscode %d text %s", url, statusCode, string(b))
			}

			customResources := []map[string]interface{}{}
			err := json.Unmarshal(b, &customResources)
			if err != nil {
				return nil, errors.Wrapf(err, "unmarshal custom resource result")
			}

			results := []*resource.CustomResource{}
			for _, customResource := range customResources {
				copy := customResource
				results = append(results, resource.ToCustomResource(copy))
			}
			return results, nil
		})
	if err != nil {
		return nil, err
	}
	return result.([]*resource.CustomResource), err
}
