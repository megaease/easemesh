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

type ingressGetter struct {
	client *meshClient
}

var _ IngressGetter = &ingressGetter{}

func (i *ingressGetter) Ingress() IngressInterface {
	return &ingressInterface{client: i.client}
}

type ingressInterface struct {
	client *meshClient
}

func (i *ingressInterface) Get(ctx context.Context, ingressID string) (*resource.Ingress, error) {
	re, err := client.NewHTTPJSON().
		GetByContext(fmt.Sprintf("http://"+i.client.server+MeshIngressURL, ingressID), nil, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "get ingress %s error", ingressID)
			}

			if statusCode >= 300 {
				return nil, errors.Errorf("call %s%s failed, return status code: %d text:%s", i.client.server, MeshIngressURL, statusCode, string(b))
			}
			ingress := &v1alpha1.Ingress{}
			err := json.Unmarshal(b, ingress)
			if err != nil {
				return nil, errors.Wrap(err, "unmarshal data to v1alpha1.Tanent error")
			}
			return resource.ToIngress(ingress), nil
		})
	if err != nil {
		return nil, err
	}

	return re.(*resource.Ingress), nil
}

func (i *ingressInterface) Patch(ctx context.Context, ingress *resource.Ingress) error {
	jsonClient := client.NewHTTPJSON()
	update := ingress.ToV1Alpha1()
	_, err := jsonClient.
		PutByContext(fmt.Sprintf("http://"+i.client.server+MeshIngressURL, ingress.Name()), &update, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "patch ingress %s error", ingress.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call PUT %s%s failed, return statuscode %d text %s", i.client.server, MeshIngressURL, statusCode, string(b))
		})
	return err
}

func (t *ingressInterface) Create(ctx context.Context, ingress *resource.Ingress) error {
	created := ingress.ToV1Alpha1()
	_, err := client.NewHTTPJSON().
		// FIXME: the standard RESTful URL of create resource is POST /v1/api/{resources} instead of POST /v1/api/{resources}/{id}.
		// Current URL form should be corrected in the feature
		PostByContext(fmt.Sprintf("http://"+t.client.server+MeshIngressURL, ingress.Name()), &created, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusConflict {
				return nil, errors.Wrapf(ConflictError, "create ingress %s error", ingress.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call Post %s%s failed, return statuscode %d text %s", t.client.server, MeshIngressURL, statusCode, string(b))
		})
	return err
}

func (i *ingressInterface) Delete(ctx context.Context, ingressID string) error {
	_, err := client.NewHTTPJSON().
		DeleteByContext(fmt.Sprintf("http://"+i.client.server+MeshIngressURL, ingressID), nil, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "delete ingress %s error", ingressID)
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call DELETE %s%s failed, return statuscode %d text %s", i.client.server, MeshIngressURL, statusCode, string(b))
		})
	return err
}

func (i *ingressInterface) List(ctx context.Context) ([]*resource.Ingress, error) {
	result, err := client.NewHTTPJSON().
		GetByContext(fmt.Sprintf("http://"+i.client.server+MeshIngressesURL), nil, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrap(NotFoundError, "list tanent error")
			}

			if statusCode >= 300 || statusCode < 200 {
				return nil, errors.Errorf("call GET %s%s failed, return statuscode %d text %s", i.client.server, MeshIngressesURL, statusCode, string(b))
			}

			ingresses := []v1alpha1.Ingress{}
			err := json.Unmarshal(b, &ingresses)
			if err != nil {
				return nil, errors.Wrapf(err, "unmarshal ingress result error")
			}

			results := []*resource.Ingress{}
			for _, ingress := range ingresses {
				results = append(results, resource.ToIngress(&ingress))
			}
			return results, nil
		})
	if err != nil {
		return nil, err
	}
	return result.([]*resource.Ingress), err
}
