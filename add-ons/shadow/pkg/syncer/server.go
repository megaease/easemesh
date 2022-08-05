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

package syncer

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/megaease/easemesh-api/v2alpha1"
	"github.com/megaease/easemesh/mesh-shadow/pkg/object"
	"github.com/megaease/easemeshctl/cmd/client/command/meshclient"
	"github.com/megaease/easemeshctl/cmd/client/resource"
	emctlclient "github.com/megaease/easemeshctl/cmd/common/client"
	"github.com/pkg/errors"
)

const (
	apiURL = "/apis/v2"
	// MeshCustomObjetWatchURL is the mesh custom resource watching path.
	MeshCustomObjetWatchURL = apiURL + "/mesh/watchcustomresources/%s"
	// MeshCustomObjectsURL is the mesh custom resource list path.
	MeshCustomObjectsURL = apiURL + "/mesh/customresources/%s"

	// MeshServiceCanaryPrefix is the service canary prefix.
	MeshServiceCanaryPrefix = "/mesh/servicecanaries"

	// MeshServiceCanaryPath is the service canary path.
	MeshServiceCanaryPath = "/mesh/servicecanaries/%s"
)

var (
	// NotFoundError indicate that the resource does not exist
	NotFoundError = errors.Errorf("resource not found")
)

// MeshControlPlane captures the set of operations for interacting with the EaseMesh REST apis.
type MeshControlPlane interface {
	List(ctx context.Context, kind string) ([]object.ShadowService, error)
	Watch(kind string) (*bufio.Reader, error)
	GetServiceCanary(name string) (*resource.ServiceCanary, error)
	CreateServiceCanary(args1 *resource.ServiceCanary) error
	PatchServiceCanary(serviceCanary *resource.ServiceCanary) error
	DeleteServiceCanary(name string) error
}

// Server represents the server of the easemesh control plane.
type Server struct {
	RequestTimeout time.Duration
	MeshServer     string
}

// NewServer create Server to access EaseMesh control plane.
func NewServer(requestTimeout time.Duration, meshServer string) *Server {
	return &Server{
		RequestTimeout: requestTimeout,
		MeshServer:     meshServer,
	}
}

// List query MeshCustomObject list from Server according to kind.
func (server *Server) List(ctx context.Context, kind string) ([]object.ShadowService, error) {
	jsonClient := emctlclient.NewHTTPJSON()
	url := fmt.Sprintf("http://"+server.MeshServer+MeshCustomObjectsURL, kind)
	result, err := jsonClient.
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrap(NotFoundError, "list service")
			}

			if statusCode >= 300 || statusCode < 200 {
				return nil, errors.Errorf("call GET %s failed, return statuscode %d text %s", url, statusCode, string(b))
			}

			var services []object.ShadowService
			err := json.Unmarshal(b, &services)
			if err != nil {
				return nil, errors.Wrapf(err, "unmarshal CustomObject result")
			}
			return services, nil
		})
	if err != nil {
		return nil, err
	}
	return result.([]object.ShadowService), err
}

// Watch listens to the custom objects of the server according to kind.
func (server *Server) Watch(kind string) (*bufio.Reader, error) {
	url := fmt.Sprintf("http://"+server.MeshServer+MeshCustomObjetWatchURL, kind)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	httpResp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, errors.Errorf("list %s objects failed,", kind)
	}
	statusCode := httpResp.StatusCode
	if statusCode == http.StatusNotFound {
		return nil, errors.Wrap(NotFoundError, "watch service")
	}

	if statusCode >= 300 || statusCode < 200 {
		return nil, errors.Errorf("call GET %s failed, return statuscode %d ", url, statusCode)
	}

	reader := bufio.NewReader(httpResp.Body)
	return reader, nil
}

// GetServiceCanary query ServiceCanary by name from EaseMesh control plane.
func (server *Server) GetServiceCanary(name string) (*resource.ServiceCanary, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), server.RequestTimeout)
	defer cancelFunc()

	url := fmt.Sprintf("http://"+server.MeshServer+apiURL+MeshServiceCanaryPath, name)
	r0, err := emctlclient.NewHTTPJSON().GetByContext(ctx, url, nil, nil).HandleResponse(func(buff []byte, statusCode int) (interface{}, error) {
		if statusCode == http.StatusNotFound {
			return nil, errors.Wrapf(NotFoundError, "get ServiceCanary %s", name)
		}
		if statusCode >= 300 {
			return nil, errors.Errorf("call %s failed, return status code %d text %+v", url, statusCode, string(buff))
		}
		ServiceCanary := &v2alpha1.ServiceCanary{}
		err := json.Unmarshal(buff, ServiceCanary)
		if err != nil {
			return nil, errors.Wrapf(err, "unmarshal data to v2alpha1.ServiceCanary")
		}
		return resource.ToServiceCanary(ServiceCanary), nil
	})
	if err != nil {
		return nil, err
	}
	return r0.(*resource.ServiceCanary), nil
}

// PatchServiceCanary update ServiceCanary by name.
func (server *Server) PatchServiceCanary(serviceCanary *resource.ServiceCanary) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), server.RequestTimeout)
	defer cancelFunc()

	url := fmt.Sprintf("http://"+server.MeshServer+apiURL+MeshServiceCanaryPath, serviceCanary.Name())
	alpha1 := serviceCanary.ToV2Alpha1()
	_, err := emctlclient.NewHTTPJSON().PutByContext(ctx, url, alpha1, nil).HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
		if statusCode == http.StatusNotFound {
			return nil, errors.Wrapf(NotFoundError, "patch ServiceCanary %s", serviceCanary.Name())
		}
		if statusCode < 300 && statusCode >= 200 {
			return nil, nil
		}
		return nil, errors.Errorf("call PUT %s failed, return statuscode %d text %+v", url, statusCode, string(b))
	})
	return err
}

// CreateServiceCanary create a new ServiceCanary.
func (server *Server) CreateServiceCanary(args1 *resource.ServiceCanary) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), server.RequestTimeout)
	defer cancelFunc()

	url := "http://" + server.MeshServer + apiURL + MeshServiceCanaryPrefix
	object := args1.ToV2Alpha1()
	_, err := emctlclient.NewHTTPJSON().PostByContext(ctx, url, object, nil).HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
		if statusCode == http.StatusConflict {
			return nil, errors.Wrapf(meshclient.ConflictError, "create ServiceCanary %s", args1.Name())
		}
		if statusCode < 300 && statusCode >= 200 {
			return nil, nil
		}
		return nil, errors.Errorf("call Post %s failed, return statuscode %d text %+v", url, statusCode, string(b))
	})
	return err
}

// DeleteServiceCanary delete ServiceCanary by name.
func (server *Server) DeleteServiceCanary(name string) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), server.RequestTimeout)
	defer cancelFunc()

	url := fmt.Sprintf("http://"+server.MeshServer+apiURL+MeshServiceCanaryPath, name)
	_, err := emctlclient.NewHTTPJSON().DeleteByContext(ctx, url, nil, nil).HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
		if statusCode == http.StatusNotFound {
			return nil, errors.Wrapf(NotFoundError, "Delete ServiceCanary %s", name)
		}
		if statusCode < 300 && statusCode >= 200 {
			return nil, nil
		}
		return nil, errors.Errorf("call Delete %s failed, return statuscode %d text %+v", url, statusCode, string(b))
	})
	return err
}
