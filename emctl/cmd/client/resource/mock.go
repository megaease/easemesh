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

package resource

import (
	"github.com/megaease/easemesh-api/v1alpha1"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
)

type (
	// Mock describes mock resource of the EaseMesh
	Mock struct {
		meta.MeshResource `yaml:",inline"`
		Spec              *v1alpha1.Mock `yaml:"spec" jsonschema:"required"`
	}
)

// ToV1Alpha1 converts a Mock resource to v1alpha1.Mock
func (m *Mock) ToV1Alpha1() *v1alpha1.Mock {
	return m.Spec
}

// ToMock converts a v1alpha1.Mock resource to a Mock resource
func ToMock(name string, mock *v1alpha1.Mock) *Mock {
	result := &Mock{
		Spec: &v1alpha1.Mock{},
	}
	result.MeshResource = NewMockResource(DefaultAPIVersion, name)
	result.Spec.Enabled = mock.Enabled
	result.Spec.Rules = mock.Rules
	return result
}
