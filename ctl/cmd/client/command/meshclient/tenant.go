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
	re, err := client.NewHTTPJSON().
		GetByContext(fmt.Sprintf("http://"+t.client.server+MeshTenantURL, tenantID), nil, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "get tenant %s error", tenantID)
			}

			if statusCode >= 300 {
				return nil, errors.Errorf("call %s%s failed, return status code: %d text:%s", t.client.server, MeshTenantURL, statusCode, string(b))
			}
			tenant := &v1alpha1.Tenant{}
			err := json.Unmarshal(b, tenant)
			if err != nil {
				return nil, errors.Wrap(err, "unmarshal data to v1alpha1.Tanent error")
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
	update := tenant.ToV1Alpha1()
	_, err := jsonClient.
		PutByContext(fmt.Sprintf("http://"+t.client.server+MeshTenantURL, tenant.Name()), update, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "patch tenant %s error", tenant.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call PUT %s%s failed, return statuscode %d text %s", t.client.server, MeshTenantURL, statusCode, string(b))
		})
	return err
}

func (t *tenantInterface) Create(ctx context.Context, tenant *resource.Tenant) error {
	created := tenant.ToV1Alpha1()
	_, err := client.NewHTTPJSON().
		// FIXME: the standard RESTful URL of create resource is POST /v1/api/{resources} instead of POST /v1/api/{resources}/{id}.
		// Current URL form should be corrected in the feature
		PostByContext(fmt.Sprintf("http://"+t.client.server+MeshTenantURL, tenant.Name()), created, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusConflict {
				return nil, errors.Wrapf(ConflictError, "create tenant %s error", tenant.Name())
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call Post %s%s failed, return statuscode %d text %s", t.client.server, MeshTenantsURL, statusCode, string(b))
		})
	return err
}

func (t *tenantInterface) Delete(ctx context.Context, tenantID string) error {
	_, err := client.NewHTTPJSON().
		DeleteByContext(fmt.Sprintf("http://"+t.client.server+MeshTenantURL, tenantID), nil, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrapf(NotFoundError, "delete tenant %s error", tenantID)
			}

			if statusCode < 300 && statusCode >= 200 {
				return nil, nil
			}
			return nil, errors.Errorf("call DELETE %s%s failed, return statuscode %d text %s", t.client.server, MeshTenantURL, statusCode, string(b))
		})
	return err
}

func (t *tenantInterface) List(ctx context.Context) ([]*resource.Tenant, error) {
	result, err := client.NewHTTPJSON().
		GetByContext(fmt.Sprintf("http://"+t.client.server+MeshTenantsURL), nil, ctx, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrap(NotFoundError, "list tanent error")
			}

			if statusCode >= 300 || statusCode < 200 {
				return nil, errors.Errorf("call GET %s%s failed, return statuscode %d text %s", t.client.server, MeshTenantsURL, statusCode, string(b))
			}

			tenants := []v1alpha1.Tenant{}
			err := json.Unmarshal(b, &tenants)
			if err != nil {
				return nil, errors.Wrapf(err, "unmarshal tanent result error")
			}

			results := []*resource.Tenant{}
			for _, tenant := range tenants {
				results = append(results, resource.ToTenant(&tenant))
			}
			return results, nil
		})
	if err != nil {
		return nil, err
	}
	return result.([]*resource.Tenant), err
}
