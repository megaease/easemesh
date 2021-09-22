package syncer

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/megaease/easemesh/mesh-shadow/pkg/common/client"
	"github.com/megaease/easemesh/mesh-shadow/pkg/object"
	"github.com/pkg/errors"
)

const (
	MeshServiceAnnotation   = "mesh.megaease.com/service-name"
	apiURL                  = "/apis/v1"
	MeshShadowServicesURL   = apiURL + "/mesh/shadowservices"
	MeshCustomObjetWatchURL = apiURL + "/mesh/watchcustomobjects/%s"
	MeshCustomObjectsURL    = apiURL + "/mesh/customobjects/%s"
)

var (
	// ConflictError indicate that the resource already exists
	ConflictError = errors.Errorf("resource already exists")
	// NotFoundError indicate that the resource does not exist
	NotFoundError = errors.Errorf("resource not found")
)

type serverInterface interface {
	List(ctx context.Context) ([]object.CustomObject, error)
	Watch(ctx context.Context) error
}

type Server struct {
	RequestTimeout time.Duration
	MeshServer     string
}

func (server *Server) List(ctx context.Context, kind string) ([]object.CustomObject, error) {
	jsonClient := client.NewHTTPJSON()
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

			services := []object.CustomObject{}
			err := json.Unmarshal(b, &services)
			if err != nil {
				return nil, errors.Wrapf(err, "unmarshal CustomObject result")
			}
			return services, nil
		})
	if err != nil {
		return nil, err
	}
	return result.([]object.CustomObject), err
}

func (server *Server) Watch(ctx context.Context, kind string) (*bufio.Reader, error) {
	url := fmt.Sprintf("http://"+server.MeshServer+MeshCustomObjetWatchURL, kind)
	httpResp, err := resty.New().R().SetContext(ctx).Get(url)
	if err != nil {
		return nil, errors.Errorf("list %s objects failed,", kind)
	}
	statusCode := httpResp.StatusCode()
	if statusCode == http.StatusNotFound {
		return nil, errors.Wrap(NotFoundError, "list service")
	}

	if statusCode >= 300 || statusCode < 200 {
		return nil, errors.Errorf("call GET %s failed, return statuscode %d text %s", url, statusCode, string(httpResp.Body()))
	}

	reader := bufio.NewReader(httpResp.RawResponse.Body)
	return reader, nil
}
