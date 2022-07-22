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
	"google.golang.org/protobuf/types/known/structpb"
)

type (
	// CustomResourceKind describes custom resource kind of the EaseMesh
	CustomResourceKind struct {
		meta.MeshResource `yaml:",inline"`
		Spec              *CustomResourceKindSpec `yaml:"spec" jsonschema:"required"`
	}

	// DynamicObject defines a dynamic object which is a map of string to interface{}.
	// The value of this map could also be a dynamic object, but in this case, its type
	// must be `map[string]interface{}`, and should not be `map[interface{}]interface{}`.
	DynamicObject map[string]interface{}

	// CustomResourceKindSpec describes the spec of a custom resource kind
	CustomResourceKindSpec struct {
		JSONSchema DynamicObject `yaml:"jsonSchema" jsonschema:"omitempty"`
	}

	// CustomResource describes custom resource of the EaseMesh
	CustomResource struct {
		meta.MeshResource `yaml:",inline"`
		Spec              map[string]interface{} `yaml:"spec" jsonschema:"required"`
	}
)

// UnmarshalYAML implements yaml.Unmarshaler
// the type of a DynamicObject field could be `map[interface{}]interface{}` if it is
// unmarshaled from yaml, but some packages, like the standard json package could not
// handle this type, so it must be converted to `map[string]interface{}`.
func (do *DynamicObject) UnmarshalYAML(unmarshal func(interface{}) error) error {
	m := map[string]interface{}{}
	if err := unmarshal(&m); err != nil {
		return err
	}

	var convert func(interface{}) interface{}
	convert = func(src interface{}) interface{} {
		switch x := src.(type) {
		case map[interface{}]interface{}:
			x2 := map[string]interface{}{}
			for k, v := range x {
				x2[k.(string)] = convert(v)
			}
			return x2
		case []interface{}:
			x2 := make([]interface{}, len(x))
			for i, v := range x {
				x2[i] = convert(v)
			}
			return x2
		}
		return src
	}

	for k, v := range m {
		m[k] = convert(v)
	}
	*do = m

	return nil
}

// ToV2Alpha1 converts an Ingress resource to v2alpha1.Ingress
func (k *CustomResourceKind) ToV2Alpha1() *v2alpha1.CustomResourceKind {
	result := &v2alpha1.CustomResourceKind{}
	result.Name = k.Name()
	if k.Spec != nil {
		s, _ := structpb.NewStruct(k.Spec.JSONSchema)
		result.JsonSchema = s
	}
	return result
}

// ToCustomResourceKind converts a v2alpha1.CustomResourceKind resource to a CustomResourceKind resource
func ToCustomResourceKind(k *v2alpha1.CustomResourceKind) *CustomResourceKind {
	result := &CustomResourceKind{
		Spec: &CustomResourceKindSpec{},
	}
	result.MeshResource = NewCustomResourceKindResource(DefaultAPIVersion, k.Name)
	if k.JsonSchema != nil {
		result.Spec.JSONSchema = k.JsonSchema.AsMap()
	}
	return result
}

// ToV2Alpha1 converts an Ingress resource to v2alpha1.Ingress
func (r *CustomResource) ToV2Alpha1() map[string]interface{} {
	result := map[string]interface{}{}
	result["name"] = r.Name()
	result["kind"] = r.Kind()
	for k, v := range r.Spec {
		result[k] = v
	}
	return result
}

// ToCustomResource converts a v2alpha1.CustomResource resource to a CustomResource resource
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
