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

package util

import (
	"fmt"

	"github.com/megaease/easemeshctl/cmd/client/resource"

	yamljsontool "github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Decoder interface {
	Decode(data []byte) (resource.MeshObject, *resource.VersionKind, error)
}

type decoder struct {
	oc resource.ObjectCreator
}

func (d *decoder) Decode(jsonBuff []byte) (resource.MeshObject, *resource.VersionKind, error) {
	yamlBuff, err := yamljsontool.JSONToYAML(jsonBuff)
	if err != nil {
		return nil, nil, fmt.Errorf("transform json %s to yaml failed: %v", jsonBuff, err)
	}

	vk := &resource.VersionKind{}
	err = yaml.Unmarshal(yamlBuff, vk)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unmarshal data to resource.VersionKind failed")
	}

	meshObject, err := d.oc.NewFromKind(*vk)
	if err != nil {
		return nil, vk, err
	}

	err = yaml.Unmarshal(yamlBuff, meshObject)
	if err != nil {
		return nil, vk, errors.Wrap(err, "unmarshal data to MeshObject error")
	}
	return meshObject, vk, nil
}

func newDefaultDecoder() Decoder {
	return &decoder{oc: resource.NewObjectCreator()}
}
