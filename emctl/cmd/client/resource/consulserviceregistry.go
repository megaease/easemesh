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
	"strings"

	"github.com/megaease/easemeshctl/cmd/client/command/printer"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
)

type (

	// ConsulServiceRegistry is the spec of ConsulServiceRegistry on Easegress.
	ConsulServiceRegistry struct {
		meta.MeshResource         `yaml:",inline"`
		ConsulServiceRegistrySpec `yaml:",inline"`
	}

	// ConsulServiceRegistryV1Alpha1 is the v1alphv1 version of mesh controller.
	ConsulServiceRegistryV1Alpha1 struct {
		Kind                      string `yaml:"kind"`
		Name                      string `yaml:"name"`
		ConsulServiceRegistrySpec `yaml:",inline"`
	}

	// ConsulServiceRegistrySpec is the admin config of mesh controller.
	ConsulServiceRegistrySpec struct {
		Address      string   `yaml:"address" jsonschema:"required"`
		Scheme       string   `yaml:"scheme" jsonschema:"required,enum=http,enum=https"`
		Datacenter   string   `yaml:"datacenter" jsonschema:"omitempty"`
		Token        string   `yaml:"token" jsonschema:"omitempty"`
		Namespace    string   `yaml:"namespace" jsonschema:"omitempty"`
		SyncInterval string   `yaml:"syncInterval" jsonschema:"required,format=duration"`
		ServiceTags  []string `yaml:"serviceTags" jsonschema:"omitempty"`
	}
)

var _ printer.TableObject = &ConsulServiceRegistry{}

// Columns returns the columns of ConsulServiceRegistry.
func (mc *ConsulServiceRegistry) Columns() []*printer.TableColumn {
	return []*printer.TableColumn{
		{
			Name:  "Address",
			Value: mc.Address,
		},
		{
			Name:  "SyncInterval",
			Value: mc.SyncInterval,
		},
		{
			Name:  "ServiceTags",
			Value: strings.Join(mc.ServiceTags, ","),
		},
	}
}

// ToV1Alpha1 converts ConsulServiceRegistry resouce to v1alpha1.
func (mc *ConsulServiceRegistry) ToV1Alpha1() *ConsulServiceRegistryV1Alpha1 {
	return &ConsulServiceRegistryV1Alpha1{
		Kind:                      mc.Kind(),
		Name:                      mc.Name(),
		ConsulServiceRegistrySpec: mc.ConsulServiceRegistrySpec,
	}
}

// ToConsulServiceRegistry converts a ConsulServiceRegistryV1Alpha1 resouce to a ConsulServiceRegistry resource.
func ToConsulServiceRegistry(consulServiceRegistry *ConsulServiceRegistryV1Alpha1) *ConsulServiceRegistry {
	return &ConsulServiceRegistry{
		MeshResource:              NewMeshResource(DefaultAPIVersion, consulServiceRegistry.Kind, consulServiceRegistry.Name),
		ConsulServiceRegistrySpec: consulServiceRegistry.ConsulServiceRegistrySpec,
	}
}
