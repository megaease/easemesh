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
	// HTTPRouteGroup describes ingress resource of the EaseMesh
	HTTPRouteGroup struct {
		meta.MeshResource `yaml:",inline"`
		Spec              *HTTPRouteGroupSpec `yaml:"spec" jsonschema:"required"`
	}

	// HTTPRouteGroupSpec wraps all route rules
	HTTPRouteGroupSpec struct {
		Matches []*v1alpha1.HTTPMatch `yaml:"matches" jsonschema:"required"`
	}
)

// ToV1Alpha1 converts an HTTPRouteGroup resource to v1alpha1.HTTPRouteGroup
func (grp *HTTPRouteGroup) ToV1Alpha1() *v1alpha1.HTTPRouteGroup {
	result := &v1alpha1.HTTPRouteGroup{}
	result.Name = grp.Name()
	if grp.Spec != nil {
		result.Matches = grp.Spec.Matches
	}
	return result
}

// ToHTTPRouteGroup converts a v1alpha1.HTTPRouteGroup resource to an HTTPRouteGroup resource
func ToHTTPRouteGroup(grp *v1alpha1.HTTPRouteGroup) *HTTPRouteGroup {
	result := &HTTPRouteGroup{
		Spec: &HTTPRouteGroupSpec{},
	}
	result.MeshResource = NewHTTPRouteGroupResource(DefaultAPIVersion, grp.Name)
	result.Spec.Matches = grp.Matches
	return result
}
