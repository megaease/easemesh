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

	"github.com/megaease/easemesh/mesh-shadow/pkg/object"
	emctlclient "github.com/megaease/easemeshctl/cmd/common/client"
	"github.com/pkg/errors"
)

const (
	apiURL                  = "/apis/v1"
	MeshCustomObjetWatchURL = apiURL + "/mesh/watchcustomresources/%s"
	MeshCustomObjectsURL    = apiURL + "/mesh/customresources/%s"
)

var (
	// NotFoundError indicate that the resource does not exist
	NotFoundError = errors.Errorf("resource not found")
)

type Server struct {
	RequestTimeout time.Duration
	MeshServer     string
}

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

func (server *Server) Watch(kind string) (*bufio.Reader, error) {
	url := fmt.Sprintf("http://"+server.MeshServer+MeshCustomObjetWatchURL, kind)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	httpResp, err := http.DefaultClient.Do(request)
	fmt.Println(err)
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
