package controllers

import (
	"context"

	"github.com/megaease/easemesh/mesh-operator/pkg/api/v1beta1"
	"github.com/megaease/easemesh/mesh-operator/pkg/controllers/resourcesyncer"
	"github.com/megaease/easemesh/mesh-operator/pkg/syncer"

	"github.com/go-logr/logr"
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
										Image: "zhaokundev/easestack-test-server:1.0-alpine",
									},
								},
							},
						},
					},
				},
			},
		}
		Expect(k8sClient.Create(context.TODO(), &meshDeployment)).To(Succeed())
		deploySyncer := resourcesyncer.NewDeploymentSyncer(k8sClient, &meshDeployment, scheme.Scheme, log)
		Expect(syncer.Sync(context.TODO(), deploySyncer, &mockRecorder{})).To(Succeed())

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
			Expect(deploy.Spec.Template.Spec.Containers[1].Name).To(Equal("easemesh-sidecar"))
		})
	})

})

func fromInt(i int) *int {
	return &i
}

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
