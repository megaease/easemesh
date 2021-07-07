package meshclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/megaease/easemeshctl/cmd/client/resource"
	"github.com/megaease/easemeshctl/cmd/common/client"

	"github.com/megaease/easemesh-api/v1alpha1"
	"github.com/pkg/errors"
)

type serviceGetter struct {
	client *meshClient
}

func (s *serviceGetter) Service() ServiceInterface {
	return &serviceInterface{client: s.client}
}

var _ ServiceInterface = &serviceInterface{}

type serviceInterface struct {
	client *meshClient
}

func (s *serviceInterface) Get(ctx context.Context, serviceID string) (*resource.Service, error) {
	jsonClient := client.NewHTTPJSON()
	r, err := jsonClient.
		GetByContext(fmt.Sprintf("http://"+s.client.server+MeshServiceURL, serviceID), nil, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "get service %s not found", serviceID)
			}
			if statusCode >= 300 {
				return nil, errors.Errorf("call %s%s failed, return status code: %d text:%s", s.client.server, MeshServiceURL, statusCode, string(b))
			}
			service := &v1alpha1.Service{}
			err := json.Unmarshal(b, service)
			if err != nil {
				return nil, errors.Wrap(err, "unmarshal data to v1alpha1.Service error")
			}
			return resource.ToService(service), nil
		})
	if err != nil {
		return nil, err
	}

	return r.(*resource.Service), nil
}

func (s *serviceInterface) Patch(ctx context.Context, service *resource.Service) error {
	jsonClient := client.NewHTTPJSON()
	update := service.ToV1Alpha1()
	_, err := jsonClient.
		PutByContext(fmt.Sprintf("http://"+s.client.server+MeshServiceURL, service.Name()), update, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, NotFoundError
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call PUT %s%s failed, return statuscode %d text %s", s.client.server, MeshServiceURL, statusCode, string(b))
		})
	return err
}

func (s *serviceInterface) Create(ctx context.Context, service *resource.Service) error {
	created := service.ToV1Alpha1()

	_, err := client.NewHTTPJSON().
		PostByContext(fmt.Sprintf("http://"+s.client.server+MeshServiceURL, service.Name()), created, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusConflict {
				return nil, errors.Wrapf(ConflictError, "create service %s error", service.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call Post %s%s failed, return statuscode %d text %s", s.client.server, MeshServiceURL, statusCode, string(b))
		})
	return err
}

func (s *serviceInterface) Delete(ctx context.Context, serviceID string) error {
	_, err := client.NewHTTPJSON().
		DeleteByContext(fmt.Sprintf("http://"+s.client.server+MeshServiceURL, serviceID), nil, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, NotFoundError
			}
			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call DELETE %s%s failed, return statuscode %d text %s", s.client.server, MeshServicesURL, statusCode, string(b))
		})
	return err
}

func (s *serviceInterface) List(ctx context.Context) ([]*resource.Service, error) {
	result, err := client.NewHTTPJSON().
		GetByContext(fmt.Sprintf("http://"+s.client.server+MeshServicesURL), nil, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, NotFoundError
			}

			if statusCode >= 300 || statusCode < 200 {
				return nil, errors.Errorf("call GET %s%s failed, return statuscode %d text %s", s.client.server, MeshServicesURL, statusCode, string(b))
			}

			services := []v1alpha1.Service{}
			err := json.Unmarshal(b, &services)
			if err != nil {
				return nil, errors.Wrapf(err, "unmarshal services result error")
			}
			results := []*resource.Service{}
			for _, ss := range services {
				results = append(results, resource.ToService(&ss))
			}
			return results, nil
		})
	if err != nil {
		return nil, err
	}
	return result.([]*resource.Service), err
}
