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

type customObjectKindGetter struct {
	client *meshClient
}

func (t *customObjectKindGetter) CustomObjectKind() CustomObjectKindInterface {
	return &customObjectKindInterface{client: t.client}
}

type customObjectKindInterface struct {
	client *meshClient
}

func (k *customObjectKindInterface) Get(ctx context.Context, customObjectKindID string) (*resource.CustomObjectKind, error) {
	url := fmt.Sprintf("http://"+k.client.server+MeshCustomObjectKindURL, customObjectKindID)
	re, err := client.NewHTTPJSON().
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "get custom object kind %s", customObjectKindID)
			}

			if statusCode >= 300 {
				return nil, errors.Errorf("call %s failed, return status code: %d text:%s", url, statusCode, string(b))
			}
			customObjectKind := &v1alpha1.CustomObjectKind{}
			err := json.Unmarshal(b, customObjectKind)
			if err != nil {
				return nil, errors.Wrap(err, "unmarshal data to v1alpha1.CustomObjectKind")
			}
			return resource.ToCustomObjectKind(customObjectKind), nil
		})
	if err != nil {
		return nil, err
	}

	return re.(*resource.CustomObjectKind), nil
}

func (k *customObjectKindInterface) Patch(ctx context.Context, customObjectKind *resource.CustomObjectKind) error {
	jsonClient := client.NewHTTPJSON()
	url := "http://" + k.client.server + MeshCustomObjectKindsURL
	update := customObjectKind.ToV1Alpha1()
	_, err := jsonClient.
		PutByContext(ctx, url, update, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "patch custom object kind %s", customObjectKind.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call PUT %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (k *customObjectKindInterface) Create(ctx context.Context, customObjectKind *resource.CustomObjectKind) error {
	created := customObjectKind.ToV1Alpha1()
	url := "http://" + k.client.server + MeshCustomObjectKindsURL
	_, err := client.NewHTTPJSON().
		// FIXME: the standard RESTful URL of create resource is POST /v1/api/{resources} instead of POST /v1/api/{resources}/{id}.
		// Current URL form should be corrected in the feature
		PostByContext(ctx, url, created, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusConflict {
				return nil, errors.Wrapf(ConflictError, "create custom object kind %s", customObjectKind.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call Post %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (k *customObjectKindInterface) Delete(ctx context.Context, customObjectKindID string) error {
	url := fmt.Sprintf("http://"+k.client.server+MeshCustomObjectKindURL, customObjectKindID)
	_, err := client.NewHTTPJSON().
		DeleteByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "delete custom object kind %s", customObjectKindID)
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call DELETE %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (k *customObjectKindInterface) List(ctx context.Context) ([]*resource.CustomObjectKind, error) {
	url := "http://" + k.client.server + MeshCustomObjectKindsURL
	result, err := client.NewHTTPJSON().
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrap(NotFoundError, "list customObjectKind")
			}

			if statusCode >= 300 || statusCode < 200 {
				return nil, errors.Errorf("call GET %s failed, return statuscode %d text %s", url, statusCode, string(b))
			}

			customObjectKinds := []v1alpha1.CustomObjectKind{}
			err := json.Unmarshal(b, &customObjectKinds)
			if err != nil {
				return nil, errors.Wrapf(err, "unmarshal custom object kind result")
			}

			results := []*resource.CustomObjectKind{}
			for _, customObjectKind := range customObjectKinds {
				copy := customObjectKind
				results = append(results, resource.ToCustomObjectKind(&copy))
			}
			return results, nil
		})
	if err != nil {
		return nil, err
	}
	return result.([]*resource.CustomObjectKind), err
}

type customObjectGetter struct {
	client *meshClient
}

func (t *customObjectGetter) CustomObject() CustomObjectInterface {
	return &customObjectInterface{client: t.client}
}

type customObjectInterface struct {
	client *meshClient
}

func (o *customObjectInterface) Get(ctx context.Context, kind, customObjectID string) (*resource.CustomObject, error) {
	url := fmt.Sprintf("http://"+o.client.server+MeshCustomObjectURL, kind, customObjectID)
	re, err := client.NewHTTPJSON().
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "get custom object %s", customObjectID)
			}

			if statusCode >= 300 {
				return nil, errors.Errorf("call %s failed, return status code: %d text:%s", url, statusCode, string(b))
			}
			customObject := map[string]interface{}{}
			err := json.Unmarshal(b, &customObject)
			if err != nil {
				return nil, errors.Wrap(err, "unmarshal data to v1alpha1.CustomObject")
			}
			return resource.ToCustomObject(customObject), nil
		})
	if err != nil {
		return nil, err
	}

	return re.(*resource.CustomObject), nil
}

func (o *customObjectInterface) Patch(ctx context.Context, customObject *resource.CustomObject) error {
	jsonClient := client.NewHTTPJSON()
	url := "http://" + o.client.server + MeshAllCustomObjectsURL
	update := customObject.ToV1Alpha1()
	_, err := jsonClient.
		PutByContext(ctx, url, update, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "patch custom object %s", customObject.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call PUT %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (o *customObjectInterface) Create(ctx context.Context, customObject *resource.CustomObject) error {
	created := customObject.ToV1Alpha1()
	url := "http://" + o.client.server + MeshAllCustomObjectsURL
	_, err := client.NewHTTPJSON().
		// FIXME: the standard RESTful URL of create resource is POST /v1/api/{resources} instead of POST /v1/api/{resources}/{id}.
		// Current URL form should be corrected in the feature
		PostByContext(ctx, url, created, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusConflict {
				return nil, errors.Wrapf(ConflictError, "create custom object %s", customObject.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call Post %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (o *customObjectInterface) Delete(ctx context.Context, kind, customObjectID string) error {
	url := fmt.Sprintf("http://"+o.client.server+MeshCustomObjectURL, kind, customObjectID)
	_, err := client.NewHTTPJSON().
		DeleteByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "delete custom object %s", customObjectID)
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call DELETE %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (o *customObjectInterface) List(ctx context.Context, kind string) ([]*resource.CustomObject, error) {
	url := fmt.Sprintf("http://"+o.client.server+MeshCustomObjectsURL, kind)
	result, err := client.NewHTTPJSON().
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrap(NotFoundError, "list custom object")
			}

			if statusCode >= 300 || statusCode < 200 {
				return nil, errors.Errorf("call GET %s failed, return statuscode %d text %s", url, statusCode, string(b))
			}

			customObjects := []map[string]interface{}{}
			err := json.Unmarshal(b, &customObjects)
			if err != nil {
				return nil, errors.Wrapf(err, "unmarshal custom object result")
			}

			results := []*resource.CustomObject{}
			for _, customObject := range customObjects {
				copy := customObject
				results = append(results, resource.ToCustomObject(copy))
			}
			return results, nil
		})
	if err != nil {
		return nil, err
	}
	return result.([]*resource.CustomObject), err
}
