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
	"time"

	"github.com/megaease/easemesh/mesh-shadow/pkg/handler/fake"
	"github.com/megaease/easemesh/mesh-shadow/pkg/object"
	"github.com/megaease/easemeshctl/cmd/client/resource"
)

// MockServer represents the server of the easemesh control plane for test.
type MockServer struct {
	RequestTimeout time.Duration
	MeshServer     string
}

//NewMockServer create MockServer for test.
func NewMockServer() *MockServer {
	return &MockServer{
		RequestTimeout: time.Second * 10,
		MeshServer:     "",
	}
}

// GetServiceCanary query ServiceCanary by name from EaseMesh control plane.
func (server *MockServer) GetServiceCanary(name string) (*resource.ServiceCanary, error) {
	return fake.NewServiceCanary(), nil
}

// List query MeshCustomObject list from Server according to kind.
func (server *MockServer) List(ctx context.Context, kind string) ([]object.ShadowService, error) {
	return nil, nil
}

// Watch listens to the custom objects of the server according to kind.
func (server *MockServer) Watch(kind string) (*bufio.Reader, error) {
	return nil, nil
}

// PatchServiceCanary update ServiceCanary by name.
func (server *MockServer) PatchServiceCanary(serviceCanary *resource.ServiceCanary) error {
	return nil
}

// CreateServiceCanary create a new ServiceCanary.
func (server *MockServer) CreateServiceCanary(args1 *resource.ServiceCanary) error {
	return nil
}

// DeleteServiceCanary delete ServiceCanary by name.
func (server *MockServer) DeleteServiceCanary(name string) error {
	return nil
}
