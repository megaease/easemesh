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
	// LoadBalance describes loadbalance resource of the EaseMesh
	LoadBalance struct {
		meta.MeshResource `yaml:",inline"`
		Spec              *v2alpha1.LoadBalance `yaml:"spec" jsonschema:"required"`
	}
)

var _ meta.TableObject = &LoadBalance{}

// Columns returns the columns of LoadBalance.
func (l *LoadBalance) Columns() []*meta.TableColumn {
	if l.Spec == nil {
		return nil
	}

	return []*meta.TableColumn{
		{
			Name:  "Policy",
			Value: l.Spec.Policy,
		},
		{
			Name:  "HeaderHashKey",
			Value: l.Spec.HeaderHashKey,
		},
	}
}

// ToV2Alpha1 converts a loadbalance resource to v2alpha1.LoadBalance
func (l *LoadBalance) ToV2Alpha1() *v2alpha1.LoadBalance {
	return l.Spec
}

// ToLoadBalance converts a v2alpha1.LoadBalance resource to a LoadBalance resource
func ToLoadBalance(name string, loadBalance *v2alpha1.LoadBalance) *LoadBalance {
	result := &LoadBalance{
		Spec: &v2alpha1.LoadBalance{},
	}
	result.MeshResource = NewLoadBalanceResource(DefaultAPIVersion, name)
	result.Spec = loadBalance
	return result
}
