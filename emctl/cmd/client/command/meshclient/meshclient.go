/*
 * Copyright (c) 2021, MegaEase
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
	"os"
	"strings"

	"github.com/megaease/easemeshctl/cmd/client/command/meshclient/fake"
)

var isTest bool

func init() {
	// For test, if the code detect it run in test, could set isTest to true
	if strings.Contains(os.Args[0], "/_test/") ||
		strings.Contains(os.Args[0], ".test") {
		isTest = true
	}
}

type meshClient struct {
	server   string
	v1Alpha1 V1Alpha1Interface
}

func (m *meshClient) V1Alpha1() V1Alpha1Interface {
	return m.v1Alpha1
}

type v1alpha1Interface struct {
	meshControllerGetter
	loadbalanceGetter
	canaryGetter
	resilienceGetter
	serviceGetter
	serviceInstanceGetter
	tenantGetter
	observabilityGetter
	ingressGetter
	serviceCanaryGetter
	customResourceKindGetter
	customResourceGetter
}

var _ V1Alpha1Interface = &v1alpha1Interface{}

// New initials a new MeshClient
func New(server string) MeshClient {
	server = strings.TrimPrefix(server, "http://")

	if isTest {
		// This is for test, in the unit test we will create a mock MeshClient
		if fake.ResourceReactorForType(server) != nil {
			return &fakeMeshClient{reactorType: server}
		}
	}

	client := &meshClient{server: server}
	alpha1 := v1alpha1Interface{
		meshControllerGetter:     meshControllerGetter{client: client},
		loadbalanceGetter:        loadbalanceGetter{client: client},
		canaryGetter:             canaryGetter{client: client},
		resilienceGetter:         resilienceGetter{client: client},
		tenantGetter:             tenantGetter{client: client},
		observabilityGetter:      observabilityGetter{client: client},
		serviceGetter:            serviceGetter{client: client},
		serviceInstanceGetter:    serviceInstanceGetter{client: client},
		ingressGetter:            ingressGetter{client: client},
		serviceCanaryGetter:      serviceCanaryGetter{client: client},
		customResourceKindGetter: customResourceKindGetter{client: client},
		customResourceGetter:     customResourceGetter{client: client},
	}
	client.v1Alpha1 = &alpha1
	return client
}
