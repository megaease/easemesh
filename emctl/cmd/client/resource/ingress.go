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
	// Ingress describes ingress resource of the EaseMesh
	Ingress struct {
		meta.MeshResource `yaml:",inline"`
		Spec              *IngressSpec `yaml:"spec" jsonschema:"required"`
	}

	// IngressSpec wraps all route rules
	IngressSpec struct {
		Rules []*v1alpha1.IngressRule `yaml:"rules" jsonschema:"required"`
	}
)

// ToV1Alpha1 converts an Ingress resource to v1alpha1.Ingress
func (ing *Ingress) ToV1Alpha1() *v1alpha1.Ingress {
	result := &v1alpha1.Ingress{}
	result.Name = ing.Name()
	if ing.Spec != nil {
		result.Rules = ing.Spec.Rules
	}
	return result
}

// ToIngress converts a v1alpha1.Ingress resource to an Ingress resource
func ToIngress(ingress *v1alpha1.Ingress) *Ingress {
	result := &Ingress{
		Spec: &IngressSpec{},
	}
	result.MeshResource = NewIngressResource(DefaultAPIVersion, ingress.Name)
	result.Spec.Rules = ingress.Rules
	return result
}
