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
	"github.com/megaease/easemeshctl/cmd/client/command/printer"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
)

type (
	// CustomObjectKind describes tenant resource of the EaseMesh
	CustomObjectKind struct {
		meta.MeshResource `yaml:",inline"`
		Spec              *CustomObjectKindSpec `yaml:"spec" jsonschema:"required"`
	}

	// CustomObjectKindSpec describes whats service resided in
	CustomObjectKindSpec struct {
		JSONSchema string `yaml:"jsonSchema" jsonschema:"omitempty"`
	}

	// CustomObject describes tenant resource of the EaseMesh
	CustomObject struct {
		meta.MeshResource `yaml:",inline"`
		Spec              map[string]interface{} `yaml:"spec" jsonschema:"required"`
	}
)

// Columns returns the columns of CustomObjectKind.
func (k *CustomObjectKind) Columns() []*printer.TableColumn {
	if k.Spec == nil {
		return nil
	}

	return []*printer.TableColumn{
		{
			Name:  "JSONSchema",
			Value: k.Spec.JSONSchema,
		},
	}
}

// ToV1Alpha1 converts an Ingress resource to v1alpha1.Ingress
func (k *CustomObjectKind) ToV1Alpha1() *v1alpha1.CustomObjectKind {
	result := &v1alpha1.CustomObjectKind{}
	result.Name = k.Name()
	if k.Spec != nil {
		result.JsonSchema = k.Spec.JSONSchema
	}
	return result
}

// ToCustomObjectKind converts a v1alpha1.CustomObjectKind resource to a CustomObjectKind resource
func ToCustomObjectKind(k *v1alpha1.CustomObjectKind) *CustomObjectKind {
	result := &CustomObjectKind{
		Spec: &CustomObjectKindSpec{},
	}
	result.MeshResource = NewCustomObjectKindResource(DefaultAPIVersion, k.Name)
	result.Spec.JSONSchema = k.JsonSchema
	return result
}

// ToV1Alpha1 converts an Ingress resource to v1alpha1.Ingress
func (o *CustomObject) ToV1Alpha1() map[string]interface{} {
	result := map[string]interface{}{}
	result["name"] = o.Name()
	result["kind"] = o.Kind()
	for k, v := range o.Spec {
		result[k] = v
	}
	return result
}

// ToCustomObject converts a v1alpha1.CustomObject resource to a CustomObject resource
func ToCustomObject(o map[string]interface{}) *CustomObject {
	result := &CustomObject{
		Spec: map[string]interface{}{},
	}
	name := o["name"].(string)
	kind := o["kind"].(string)
	result.MeshResource = NewMeshResource(DefaultAPIVersion, kind, name)
	delete(o, "name")
	delete(o, "kind")
	result.Spec = o
	return result
}
