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

type tenantGetter struct {
	client *meshClient
}

func (t *tenantGetter) Tenant() TenantInterface {
	return &tenantInterface{client: t.client}
}

type tenantInterface struct {
	client *meshClient
}

func (t *tenantInterface) Get(ctx context.Context, tenantID string) (*resource.Tenant, error) {
	url := fmt.Sprintf("http://"+t.client.server+MeshTenantURL, tenantID)
	re, err := client.NewHTTPJSON().
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "get tenant %s", tenantID)
			}

			if statusCode >= 300 {
				return nil, errors.Errorf("call %s failed, return status code: %d text:%s", url, statusCode, string(b))
			}
			tenant := &v1alpha1.Tenant{}
			err := json.Unmarshal(b, tenant)
			if err != nil {
				return nil, errors.Wrap(err, "unmarshal data to v1alpha1.Tenant")
			}
			return resource.ToTenant(tenant), nil
		})
	if err != nil {
		return nil, err
	}

	return re.(*resource.Tenant), nil
}

func (t *tenantInterface) Patch(ctx context.Context, tenant *resource.Tenant) error {
	jsonClient := client.NewHTTPJSON()
	url := fmt.Sprintf("http://"+t.client.server+MeshTenantURL, tenant.Name())
	update := tenant.ToV1Alpha1()
	_, err := jsonClient.
		PutByContext(ctx, url, update, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "patch tenant %s", tenant.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call PUT %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (t *tenantInterface) Create(ctx context.Context, tenant *resource.Tenant) error {
	created := tenant.ToV1Alpha1()
	url := fmt.Sprintf("http://" + t.client.server + MeshTenantsURL)
	_, err := client.NewHTTPJSON().
		PostByContext(ctx, url, created, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusConflict {
				return nil, errors.Wrapf(ConflictError, "create tenant %s", tenant.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call Post %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (t *tenantInterface) Delete(ctx context.Context, tenantID string) error {
	url := fmt.Sprintf("http://"+t.client.server+MeshTenantURL, tenantID)
	_, err := client.NewHTTPJSON().
		DeleteByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "delete tenant %s", tenantID)
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call DELETE %s failed, return statuscode %d text %s", url, statusCode, string(b))
		})
	return err
}

func (t *tenantInterface) List(ctx context.Context) ([]*resource.Tenant, error) {
	url := fmt.Sprintf("http://" + t.client.server + MeshTenantsURL)
	result, err := client.NewHTTPJSON().
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrap(NotFoundError, "list tenant")
			}

			if statusCode >= 300 || statusCode < 200 {
				return nil, errors.Errorf("call GET %s failed, return statuscode %d text %s", url, statusCode, string(b))
			}

			tenants := []v1alpha1.Tenant{}
			err := json.Unmarshal(b, &tenants)
			if err != nil {
				return nil, errors.Wrapf(err, "unmarshal tenant result")
			}

			results := []*resource.Tenant{}
			for _, tenant := range tenants {
				copy := tenant
				results = append(results, resource.ToTenant(&copy))
			}
			return results, nil
		})
	if err != nil {
		return nil, err
	}
	return result.([]*resource.Tenant), err
}
