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

package handler

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/megaease/easemesh/mesh-shadow/pkg/object"
	appv1 "k8s.io/api/apps/v1"
	k8Yaml "k8s.io/apimachinery/pkg/util/yaml"
)

func TestCloneHandler_CloneDeploymentSpec(t *testing.T) {

	data, _ := os.ReadFile("./original_deployment.yaml")

	decoder := k8Yaml.NewYAMLOrJSONDecoder(bytes.NewReader(data), 1000)

	sourceDeployment := &appv1.Deployment{}
	err := decoder.Decode(sourceDeployment)
	if err != nil {
		fmt.Println(err)
	}

	handler := CloneHandler{}
	service := &object.ShadowService{
		ServiceName: "visits-service",
		Namespace:   "default",
		Name:        "visits-service-shadow",
	}
	handler.CloneDeploymentSpec(sourceDeployment, service)

}
