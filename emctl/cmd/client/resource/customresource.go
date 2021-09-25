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

package resource

import (
	"github.com/megaease/easemesh-api/v1alpha1"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
	"google.golang.org/protobuf/types/known/structpb"
)

type (
	// CustomResourceKind describes custom resource kind of the EaseMesh
	CustomResourceKind struct {
		meta.MeshResource `yaml:",inline"`
		Spec              *CustomResourceKindSpec `yaml:"spec" jsonschema:"required"`
	}

	// CustomResourceKindSpec describes the spec of a custom resource kind
	CustomResourceKindSpec struct {
		JSONSchema map[string]interface{} `yaml:"jsonSchema" jsonschema:"omitempty"`
	}

	// CustomResource describes custom resource of the EaseMesh
	CustomResource struct {
		meta.MeshResource `yaml:",inline"`
		Spec              map[string]interface{} `yaml:"spec" jsonschema:"required"`
	}
)

// ToV1Alpha1 converts an Ingress resource to v1alpha1.Ingress
func (k *CustomResourceKind) ToV1Alpha1() *v1alpha1.CustomResourceKind {
	result := &v1alpha1.CustomResourceKind{}
	result.Name = k.Name()
	if k.Spec != nil {
		s, _ := structpb.NewStruct(k.Spec.JSONSchema)
		result.JsonSchema = s
	}
	return result
}

// ToCustomResourceKind converts a v1alpha1.CustomResourceKind resource to a CustomResourceKind resource
func ToCustomResourceKind(k *v1alpha1.CustomResourceKind) *CustomResourceKind {
	result := &CustomResourceKind{
		Spec: &CustomResourceKindSpec{},
	}
	result.MeshResource = NewCustomResourceKindResource(DefaultAPIVersion, k.Name)
	if k.JsonSchema != nil {
		result.Spec.JSONSchema = k.JsonSchema.AsMap()
	}
	return result
}

// ToV1Alpha1 converts an Ingress resource to v1alpha1.Ingress
func (r *CustomResource) ToV1Alpha1() map[string]interface{} {
	result := map[string]interface{}{}
	result["name"] = r.Name()
	result["kind"] = r.Kind()
	for k, v := range r.Spec {
		result[k] = v
	}
	return result
}

// ToCustomResource converts a v1alpha1.CustomResource resource to a CustomResource resource
func ToCustomResource(r map[string]interface{}) *CustomResource {
	result := &CustomResource{
		Spec: map[string]interface{}{},
	}
	name := r["name"].(string)
	kind := r["kind"].(string)
	result.MeshResource = NewMeshResource(DefaultAPIVersion, kind, name)
	delete(r, "name")
	delete(r, "kind")
	result.Spec = r
	return result
}
