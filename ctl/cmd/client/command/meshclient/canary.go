package meshclient

import (
	"context"
	"fmt"
	"net/http"

	"github.com/megaease/easemesh-api/v1alpha1"
	"github.com/megaease/easemeshctl/cmd/client/resource"
	"github.com/megaease/easemeshctl/cmd/common/client"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type canaryGetter struct {
	client *meshClient
}

func (c *canaryGetter) Canary() CanaryInterface {
	return &canaryInterface{client: c.client}
}

var _ CanaryInterface = &canaryInterface{}

type canaryInterface struct {
	client *meshClient
}

func (c *canaryInterface) Get(ctx context.Context, serviceID string) (*resource.Canary, error) {
	jsonClient := client.NewHTTPJSON()
	r, err := jsonClient.
		GetByContext(fmt.Sprintf("http://"+c.client.server+MeshServiceCanaryURL, serviceID), nil, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "get canary %s error", serviceID)
			}

			if statusCode >= 300 {
				return nil, errors.Errorf("call %s%s failed, return status code: %d text:%s", c.client.server, MeshServiceCanaryURL, statusCode, string(b))
			}
			canary := &v1alpha1.Canary{}
			err := yaml.Unmarshal(b, canary)
			if err != nil {
				return nil, errors.Wrap(err, "unmarshal data to v1alpha1.Canary error")
			}
			return resource.ToCanary(serviceID, canary), nil
		})
	if err != nil {
		return nil, err
	}

	return r.(*resource.Canary), nil
}

func (c *canaryInterface) Patch(ctx context.Context, canary *resource.Canary) error {
	jsonClient := client.NewHTTPJSON()
	update := canary.ToV1Alpha1()
	_, err := jsonClient.
		PutByContext(fmt.Sprintf("http://"+c.client.server+MeshServiceCanaryURL, canary.Name()), &update, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "patch canary %s error", canary.Name())
			}
			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call PUT %s%s failed, return statuscode %d text %s", c.client.server, MeshServiceCanaryURL, statusCode, string(b))
		})
	return err
}

func (c *canaryInterface) Create(ctx context.Context, canary *resource.Canary) error {
	created := canary.ToV1Alpha1()
	_, err := client.NewHTTPJSON().
		PostByContext(fmt.Sprintf("http://"+c.client.server+MeshServiceCanaryURL, canary.Name()), &created, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusConflict {
				return nil, errors.Wrapf(ConflictError, "create canary %s error", canary.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call Post %s%s failed, return statuscode %d text %s", c.client.server, MeshServiceCanaryURL, statusCode, string(b))
		})
	return err
}

func (c *canaryInterface) Delete(ctx context.Context, serviceID string) error {
	_, err := client.NewHTTPJSON().
		PostByContext(fmt.Sprintf("http://"+c.client.server+MeshServiceCanaryURL, serviceID), nil, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "delete canary %s error", serviceID)
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call DELETE %s%s failed, return statuscode %d text %s", c.client.server, MeshServiceCanaryURL, statusCode, string(b))
		})
	return err
}

func (c *canaryInterface) List(ctx context.Context) ([]resource.Canary, error) {
	result, err := client.NewHTTPJSON().
		GetByContext(fmt.Sprintf("http://"+c.client.server+MeshServicesURL), nil, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrap(NotFoundError, "list service error")
			}

			if statusCode >= 300 || statusCode < 200 {
				return nil, errors.Errorf("call GET %s%s failed, return statuscode %d text %s", c.client.server, MeshServicesURL, statusCode, string(b))
			}

			services := []v1alpha1.Service{}
			err := yaml.Unmarshal(b, services)
			if err != nil {
				return nil, errors.Wrapf(err, "unmarshal services result error")
			}
			results := []resource.Canary{}
			for _, ss := range services {
				if ss.Canary != nil {
					results = append(results, resource.ToCanary(ss.Name, ss.Canary))
				}
			}
			return results, nil
		})

	return result.([]resource.Canary), err
}
