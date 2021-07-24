package controllers

import (
	"fmt"
	"testing"

	"github.com/megaease/easemesh/mesh-operator/pkg/base"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func TestDeploymentModifier(t *testing.T) {
	deploy := &v1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-vets-service",
		},
		Spec: v1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "test-vets-service",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "test-vets-service",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "test-vets-service",
							Image:   "megaease/spring-petclinic-vets-service:latest",
							Command: []string{"/bin/sh"},
							Args:    []string{"-c", "java -server -Xmx1024m -Xms1024m -Dspring.profiles.active=sit -Djava.security.egd=file:/dev/./urandom  org.springframework.boot.loader.JarLauncher"},
						},
					},
				},
			},
		},
	}

	baseRuntime := &base.Runtime{
		Name:             "test-runtime-name",
		ImageRegistryURL: "docker.io",
	}

	service := &meshService{
		Name: "test-service",
		Labels: map[string]string{
			"app":     "test-vets-service",
			"version": "beta",
		},
		AppContainerName: "test-vets-service",
		ApplicationPort:  9000,
		AliveProbeURL:    "http://localhost:9000/health",
	}

	modifier := newDeploymentModifier(baseRuntime, service, deploy)
	err := modifier.modify()
	if err != nil {
		t.Fatal(err)
	}

	buff, err := yaml.Marshal(deploy)
	if err != nil {
		t.Fatalf("mashal failed: %v", err)
	}
	fmt.Printf("%s", buff)
}
