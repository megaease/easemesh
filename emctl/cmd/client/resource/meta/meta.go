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

package meta

type (

	// VersionKind holds version and kind information for APIs
	VersionKind struct {
		APIVersion string `yaml:"apiVersion" yaml:"apiVersion" jsonschema:"omitempty"`
		Kind       string `yaml:"kind" yaml:"kind" jsonschema:"required"`
	}

	// MetaData is meta data for resources of the EaseMesh
	MetaData struct {
		Name   string            `yaml:"name" yaml:"name" jsonschema:"required"`
		Labels map[string]string `yaml:"labels,omitempty" yaml:"labels,omitempty" jsonschema:"omitempty"`
	}

	// MeshResource holds common information for a resource of the EaseMesh
	MeshResource struct {
		VersionKind `yaml:",inline" yaml:",inline"`
		MetaData    MetaData `yaml:"metadata" yaml:"metadata" jsonschema:"required"`
	}

	// MeshObject describes what's feature of a comman EaseMesh object
	MeshObject interface {
		Name() string
		Kind() string
		APIVersion() string
		Labels() map[string]string
	}
	// TableColumn is the user-defined table column.
	TableColumn struct {
		Name  string
		Value string
	}

	// TableObject is the object which wants to
	// customize its own output in format table.
	TableObject interface {
		Columns() []*TableColumn
	}
)

// Name returns name of the EaseMesh resource
func (m *MeshResource) Name() string {
	return m.MetaData.Name
}

// Kind returns kind of the EaseMesh resource
func (m *MeshResource) Kind() string {
	return m.VersionKind.Kind
}

// APIVersion returns api version of the EaseMesh resource
func (m *MeshResource) APIVersion() string {
	return m.VersionKind.APIVersion
}

// Labels returns labels of the EaseMesh resource
func (m *MeshResource) Labels() map[string]string {
	return m.MetaData.Labels
}
