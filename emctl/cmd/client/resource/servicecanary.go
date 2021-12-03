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
	"sort"
	"strconv"
	"strings"

	"github.com/megaease/easemesh-api/v1alpha1"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
)

type (
	// ServiceCanary describes canary resource of the EaseMesh.
	ServiceCanary struct {
		meta.MeshResource `yaml:",inline"`
		Spec              *ServiceCanarySpec `yaml:"spec" jsonschema:"required"`
	}

	// ServiceCanarySpec is the service canary spec.
	ServiceCanarySpec struct {
		Priority     int32
		Selector     *v1alpha1.ServiceSelector `yaml:"selector" jsonschema:"required"`
		TrafficRules *v1alpha1.TrafficRules    `yaml:"trafficRules" jsonschema:"required"`
	}
)

var _ meta.TableObject = &ServiceCanary{}

// Columns returns the columns of ServiceCanary.
func (sc *ServiceCanary) Columns() []*meta.TableColumn {
	if sc.Spec == nil {
		return nil
	}

	labels := []string{}
	if sc.Spec.Selector != nil {
		for k, v := range sc.Spec.Selector.MatchInstanceLabels {
			labels = append(labels, strings.Join([]string{k, v}, "="))
		}
	}
	sort.Strings(labels)

	return []*meta.TableColumn{
		{
			Name:  "Services",
			Value: strings.Join(sc.Spec.Selector.MatchServices, ","),
		},
		{
			Name:  "InstanceLabels",
			Value: strings.Join(labels, ","),
		},
		{
			Name:  "Priority",
			Value: strconv.Itoa(int(sc.Spec.Priority)),
		},
	}
}

// ToV1Alpha1 converts a ServiceCanary resource to v1alpha1.ServiceCanary.
func (sc *ServiceCanary) ToV1Alpha1() *v1alpha1.ServiceCanary {
	result := &v1alpha1.ServiceCanary{}
	result.Name = sc.Name()
	if sc.Spec != nil {
		result.Selector = sc.Spec.Selector
		result.TrafficRules = sc.Spec.TrafficRules
		result.Priority = sc.Spec.Priority
	}

	return result
}

// ToServiceCanary converts a v1alpha1.ServiceCanary resource to a ServiceCanary resource.
func ToServiceCanary(serviceCanary *v1alpha1.ServiceCanary) *ServiceCanary {
	result := &ServiceCanary{
		Spec: &ServiceCanarySpec{},
	}

	result.MeshResource = NewServiceCanaryResource(DefaultAPIVersion, serviceCanary.Name)
	result.Spec.Priority = serviceCanary.Priority
	result.Spec.Selector = serviceCanary.Selector
	result.Spec.TrafficRules = serviceCanary.TrafficRules

	return result
}
