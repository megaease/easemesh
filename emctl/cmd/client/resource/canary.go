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
	// Canary describes canary resource of the EaseMesh
	Canary struct {
		meta.MeshResource `yaml:",inline"`
		Spec              *v1alpha1.Canary `yaml:"spec" jsonschema:"required"`
	}
)

// ToV1Alpha1 converts a Canary resource to v1alpha1.Canary
func (c *Canary) ToV1Alpha1() *v1alpha1.Canary {
	return c.Spec
}

// ToCanary converts a v1alpha1.Canary resource to a Canary resource
func ToCanary(name string, canary *v1alpha1.Canary) *Canary {
	result := &Canary{
		Spec: &v1alpha1.Canary{},
	}
	result.MeshResource = NewCanaryResource(DefaultAPIVersion, name)
	result.Spec.CanaryRules = canary.CanaryRules
	return result
}
