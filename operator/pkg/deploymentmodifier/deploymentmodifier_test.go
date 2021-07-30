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

package deploymentmodifier

import (
	_ "embed"

	"github.com/megaease/easemesh/mesh-operator/pkg/base"
	v1 "k8s.io/api/apps/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/yaml"
)

var (
	//go:embed original_deployment.yaml
	originalDeployStr string

	//go:embed want_deployment.yaml
	wantDeployStr string
)

var _ = Describe("DeploymentModifier", func() {
	It("modify deployment", func() {
		originalDeploy := &v1.Deployment{}
		wantDeploy := &v1.Deployment{}

		Expect(yaml.Unmarshal([]byte(originalDeployStr), originalDeploy)).To(Succeed())
		Expect(yaml.Unmarshal([]byte(wantDeployStr), wantDeploy)).To(Succeed())

		baseRuntime := &base.Runtime{
			Name:            "test-runtime-name",
			ImagePullPolicy: "IfNotPresent",
		}

		service := &MeshService{
			Name: "vets-service",
			Labels: map[string]string{
				"app":     "vets-service",
				"version": "beta",
			},
			AppContainerName: "vets-service",
			ApplicationPort:  9000,
			AliveProbeURL:    "http://localhost:9000/health",
		}

		// gotDeployStr, err := yaml.Marshal(originalDeploy)
		// Expect(err).ShouldNot(HaveOccurred())
		// fmt.Printf("%s\n", gotDeployStr)

		modifier := New(baseRuntime, service, originalDeploy)
		Expect(modifier.Modify()).To(Succeed())
		Expect(originalDeploy).To(Equal(wantDeploy))
	})
})
