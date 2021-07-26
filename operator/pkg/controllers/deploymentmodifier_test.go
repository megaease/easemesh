package controllers

import (
	_ "embed"
	"fmt"

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
	defer GinkgoRecover()

	originalDeploy := &v1.Deployment{}
	wantDeploy := &v1.Deployment{}

	Expect(yaml.Unmarshal([]byte(originalDeployStr), originalDeploy)).To(Succeed())
	Expect(yaml.Unmarshal([]byte(wantDeployStr), wantDeploy)).To(Succeed())

	baseRuntime := &base.Runtime{
		Name:            "test-runtime-name",
		ImagePullPolicy: "IfNotPresent",
	}

	service := &meshService{
		Name: "vets-service",
		Labels: map[string]string{
			"app":     "vets-service",
			"version": "beta",
		},
		AppContainerName: "vets-service",
		ApplicationPort:  9000,
		AliveProbeURL:    "http://localhost:9000/health",
	}

	modifier := newDeploymentModifier(baseRuntime, service, originalDeploy)
	Expect(modifier.modify()).To(Succeed())

	gotDeployStr, err := yaml.Marshal(originalDeploy)
	Expect(err).ShouldNot(HaveOccurred())
	fmt.Printf("%s\n", gotDeployStr)

	Expect(originalDeploy).To(Equal(wantDeploy))
})
