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
	"github.com/megaease/easemesh-api/v2alpha1"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
)

type (
	// TrafficTarget describes ingress resource of the EaseMesh
	TrafficTarget struct {
		meta.MeshResource `yaml:",inline"`
		Spec              *TrafficTargetSpec `yaml:"spec" jsonschema:"required"`
	}

	// TrafficTargetSpec wraps all route rules
	TrafficTargetSpec struct {
		Destination *v2alpha1.IdentityBindingSubject   `yaml:"destination" jsonschema:"required"`
		Sources     []*v2alpha1.IdentityBindingSubject `yaml:"sources" jsonschema:"required"`
		Rules       []*v2alpha1.TrafficTargetRule      `yaml:"rules" jsonschema:"required"`
	}
)

// ToV2Alpha1 converts an TrafficTarget resource to v2alpha1.TrafficTarget
func (tt *TrafficTarget) ToV2Alpha1() *v2alpha1.TrafficTarget {
	result := &v2alpha1.TrafficTarget{}
	result.Name = tt.Name()
	if tt.Spec != nil {
		result.Destination = tt.Spec.Destination
		result.Sources = tt.Spec.Sources
		result.Rules = tt.Spec.Rules
	}
	return result
}

// ToTrafficTarget converts a v2alpha1.TrafficTarget resource to an TrafficTarget resource
func ToTrafficTarget(tt *v2alpha1.TrafficTarget) *TrafficTarget {
	result := &TrafficTarget{
		Spec: &TrafficTargetSpec{},
	}
	result.MeshResource = NewTrafficTargetResource(DefaultAPIVersion, tt.Name)
	result.Spec.Destination = tt.Destination
	result.Spec.Sources = tt.Sources
	result.Spec.Rules = tt.Rules
	return result
}
