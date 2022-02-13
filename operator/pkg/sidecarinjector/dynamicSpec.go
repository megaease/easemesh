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

package sidecarinjector

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/megaease/easemesh/mesh-operator/pkg/base"
	"gopkg.in/yaml.v2"
)

type (
	// meshControllerSpec is the part that operator cares of mesh controller spec.
	meshControllerSpec struct {
		ImageRegistryURL          string `yaml:"imageRegistryURL"`
		ImagePullPolicy           string `yaml:"imagePullPolicy"`
		SidecarImageName          string `yaml:"sidecarImageName"`
		AgentInitializerImageName string `yaml:"agentInitializerImageName"`
		Log4jConfigName           string `yaml:"log4jConfigName"`
	}

	dynamicSpec struct {
		runtime            *base.Runtime
		meshControllerSpec *meshControllerSpec
	}
)

const (
	meshControllerName = "easemesh-controller"
)

func newDynamicSpec(runtime *base.Runtime) *dynamicSpec {
	ds := &dynamicSpec{
		runtime: runtime,
	}

	ds.meshControllerSpec = ds.staticSpec()

	url := fmt.Sprintf("http://%s/apis/v1/objects/%s", ds.runtime.APIAddr, meshControllerName)
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		ds.runtime.Log.Error(err, "get mesh controller spec failed", "url", url)
		return ds
	}

	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ds.runtime.Log.Error(err, "read body failed", "url", url)
		return ds
	}
	resp.Body.Close()

	spec := ds.staticSpec()
	err = yaml.Unmarshal(buff, spec)
	if err != nil {
		ds.runtime.Log.Error(err, "unmarshal yaml failed", "body", buff, "spec", spec)
		return ds
	}

	ds.meshControllerSpec = spec

	return ds
}

func (ds *dynamicSpec) spec() *meshControllerSpec {
	return ds.meshControllerSpec
}

func (ds *dynamicSpec) staticSpec() *meshControllerSpec {
	return &meshControllerSpec{
		ImageRegistryURL:          ds.runtime.ImageRegistryURL,
		ImagePullPolicy:           ds.runtime.ImagePullPolicy,
		SidecarImageName:          ds.runtime.SidecarImageName,
		AgentInitializerImageName: ds.runtime.AgentInitializerImageName,
		Log4jConfigName:           ds.runtime.Log4jConfigName,
	}
}
