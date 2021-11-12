/*
Copyright (c) 2021, MegaEase
All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// code generated by github.com/megaease/easemeshctl/cmd/generator, DO NOT EDIT.
package meshclient

import (
	"context"
	"encoding/json"
	"fmt"
	v1alpha1 "github.com/megaease/easemesh-api/v1alpha1"
	resource "github.com/megaease/easemeshctl/cmd/client/resource"
	client "github.com/megaease/easemeshctl/cmd/common/client"
	errors "github.com/pkg/errors"
	"net/http"
)

type serviceGetter struct {
	client *meshClient
}
type serviceInterface struct {
	client *meshClient
}

func (s *serviceGetter) Service() ServiceInterface {
	return &serviceInterface{client: s.client}
}
func (s *serviceInterface) Get(args0 context.Context, args1 string) (*resource.Service, error) {
	url := fmt.Sprintf("http://"+s.client.server+apiURL+"/mesh/"+"services/%s", args1)
	r0, err := client.NewHTTPJSON().GetByContext(args0, url, nil, nil).HandleResponse(func(buff []byte, statusCode int) (interface{}, error) {
		if statusCode == http.StatusNotFound {
			return nil, errors.Wrapf(NotFoundError, "get Service %s", args1)
		}
		if statusCode >= 300 {
			return nil, errors.Errorf("call %s failed, return status code %d text %+v", url, statusCode, string(buff))
		}
		Service := &v1alpha1.Service{}
		err := json.Unmarshal(buff, Service)
		if err != nil {
			return nil, errors.Wrapf(err, "unmarshal data to v1alpha1.Service")
		}
		return resource.ToService(Service), nil
	})
	if err != nil {
		return nil, err
	}
	return r0.(*resource.Service), nil
}
func (s *serviceInterface) Patch(args0 context.Context, args1 *resource.Service) error {
	url := fmt.Sprintf("http://"+s.client.server+apiURL+"/mesh/"+"services/%s", args1.Name())
	object := args1.ToV1Alpha1()
	_, err := client.NewHTTPJSON().PutByContext(args0, url, object, nil).HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
		if statusCode == http.StatusNotFound {
			return nil, errors.Wrapf(NotFoundError, "patch Service %s", args1.Name())
		}
		if statusCode < 300 && statusCode >= 200 {
			return nil, nil
		}
		return nil, errors.Errorf("call PUT %s failed, return statuscode %d text %+v", url, statusCode, string(b))
	})
	return err
}
func (s *serviceInterface) Create(args0 context.Context, args1 *resource.Service) error {
	url := "http://" + s.client.server + apiURL + "/mesh/services"
	object := args1.ToV1Alpha1()
	_, err := client.NewHTTPJSON().PostByContext(args0, url, object, nil).HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
		if statusCode == http.StatusConflict {
			return nil, errors.Wrapf(ConflictError, "create Service %s", args1.Name())
		}
		if statusCode < 300 && statusCode >= 200 {
			return nil, nil
		}
		return nil, errors.Errorf("call Post %s failed, return statuscode %d text %+v", url, statusCode, string(b))
	})
	return err
}
func (s *serviceInterface) Delete(args0 context.Context, args1 string) error {
	url := fmt.Sprintf("http://"+s.client.server+apiURL+"/mesh/"+"services/%s", args1)
	_, err := client.NewHTTPJSON().DeleteByContext(args0, url, nil, nil).HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
		if statusCode == http.StatusNotFound {
			return nil, errors.Wrapf(NotFoundError, "Delete Service %s", args1)
		}
		if statusCode < 300 && statusCode >= 200 {
			return nil, nil
		}
		return nil, errors.Errorf("call Delete %s failed, return statuscode %d text %+v", url, statusCode, string(b))
	})
	return err
}
func (s *serviceInterface) List(args0 context.Context) ([]*resource.Service, error) {
	url := "http://" + s.client.server + apiURL + "/mesh/services"
	result, err := client.NewHTTPJSON().GetByContext(args0, url, nil, nil).HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
		if statusCode == http.StatusNotFound {
			return nil, errors.Wrapf(NotFoundError, "list service")
		}
		if statusCode >= 300 && statusCode < 200 {
			return nil, errors.Errorf("call GET %s failed, return statuscode %d text %+v", url, statusCode, b)
		}
		service := []v1alpha1.Service{}
		err := json.Unmarshal(b, &service)
		if err != nil {
			return nil, errors.Wrapf(err, "unmarshal data to v1alpha1.")
		}
		results := []*resource.Service{}
		for _, item := range service {
			copy := item
			results = append(results, resource.ToService(&copy))
		}
		return results, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*resource.Service), nil
}
