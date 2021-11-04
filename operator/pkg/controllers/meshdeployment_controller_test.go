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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/megaease/easemesh/mesh-operator/pkg/api/v1beta1"
	"github.com/megaease/easemesh/mesh-operator/pkg/base"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	meshDeploymentName = "test-server-v1"
	namespace          = "default"
)

var (
	key = types.NamespacedName{Namespace: namespace, Name: meshDeploymentName}
)

var _ = Describe("meshdeployment controller", func() {
	var meshDeployment v1beta1.MeshDeployment
	var log logr.Logger
	BeforeEach(func() {
		log = ctrl.Log.WithName("controllers").WithName("MeshDeployment")
		baseRuntime := &base.Runtime{
			Name:     "mesh-controller-test",
			Client:   k8sClient,
			Scheme:   scheme.Scheme,
			Recorder: &mockRecorder{},
			Log:      log,
		}
		meshDeployment = v1beta1.MeshDeployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      meshDeploymentName,
				Namespace: namespace,
			},
			Spec: v1beta1.MeshDeploymentSpec{
				Service: v1beta1.ServiceSpec{
					Name: "test-server",
					Labels: map[string]string{
						"canary": "internal",
					},
				},
				Deploy: v1beta1.DeploySpec{
					DeploymentSpec: v1.DeploymentSpec{
						Replicas: fromInt32(2),
						Selector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"app": "test-server",
							},
						},
						Template: corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									{
										Name:  "test-server",
										Image: "megaease/non-existent:1.0-alpine",
										Ports: []corev1.ContainerPort{
											{
												Name:          "test-port",
												ContainerPort: 8080,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		Expect(k8sClient.Create(context.TODO(), &meshDeployment)).To(Succeed())

		meshDeploymentReconciler := &MeshDeploymentReconciler{
			Runtime: baseRuntime,
		}
		req := ctrl.Request{NamespacedName: key}

		_, err := meshDeploymentReconciler.Reconcile(context.TODO(), req)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		// remove created cluster
		k8sClient.Delete(context.TODO(), &meshDeployment)
	})

	Context("normal deploy meshdeployment", func() {
		deploy := v1.Deployment{}
		It("should has a deployment", func() {
			Expect(k8sClient.Get(context.TODO(), key, &deploy)).To(Succeed())
		})
	})

	Context("normal deploy meshdeployment check injected container ", func() {
		deploy := v1.Deployment{}
		It("should has a injected container named with easemesh-sidecar", func() {
			Expect(k8sClient.Get(context.TODO(), key, &deploy)).To(Succeed())
			Expect(len(deploy.Spec.Template.Spec.Containers)).To(Equal(2))
		})
	})

})

func fromInt32(i int32) *int32 {
	return &i
}

type mockRecorder struct {
}

var _ record.EventRecorder = &mockRecorder{}

func (m *mockRecorder) Event(object runtime.Object, eventtype, reason, message string) {
}

func (m *mockRecorder) Eventf(object runtime.Object, eventtype, reason, messageFmt string, args ...interface{}) {
}

func (m *mockRecorder) AnnotatedEventf(object runtime.Object, annotations map[string]string, eventtype, reason, messageFmt string, args ...interface{}) {
}
