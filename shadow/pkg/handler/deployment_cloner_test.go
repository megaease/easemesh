package handler

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

func TestShadowServiceCloner_cloneDeployment(t *testing.T) {
	type fields struct {
		KubeClient    kubernetes.Interface
		RunTimeClient *client.Client
		CRDClient     *rest.RESTClient
	}
	type args struct {
		sourceDeployment *appsV1.Deployment
		shadowService    *object.ShadowService
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   ShadowDeploymentFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cloner := &ShadowServiceCloner{
				KubeClient:    tt.fields.KubeClient,
				RunTimeClient: tt.fields.RunTimeClient,
				CRDClient:     tt.fields.CRDClient,
			}
			if got := cloner.cloneDeployment(tt.args.sourceDeployment, tt.args.shadowService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cloneDeployment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShadowServiceCloner_cloneDeploymentSpec(t *testing.T) {
	type fields struct {
		KubeClient    kubernetes.Interface
		RunTimeClient *client.Client
		CRDClient     *rest.RESTClient
	}
	type args struct {
		sourceDeployment *appsV1.Deployment
		shadowService    *object.ShadowService
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *appsV1.Deployment
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cloner := &ShadowServiceCloner{
				KubeClient:    tt.fields.KubeClient,
				RunTimeClient: tt.fields.RunTimeClient,
				CRDClient:     tt.fields.CRDClient,
			}
			if got := cloner.cloneDeploymentSpec(tt.args.sourceDeployment, tt.args.shadowService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cloneDeploymentSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShadowServiceCloner_decorateShadowConfiguration(t *testing.T) {
	type fields struct {
		KubeClient    kubernetes.Interface
		RunTimeClient *client.Client
		CRDClient     *rest.RESTClient
	}
	type args struct {
		deployment       *appsV1.Deployment
		sourceDeployment *appsV1.Deployment
		shadowService    *object.ShadowService
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *appsV1.Deployment
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cloner := &ShadowServiceCloner{
				KubeClient:    tt.fields.KubeClient,
				RunTimeClient: tt.fields.RunTimeClient,
				CRDClient:     tt.fields.CRDClient,
			}
			if got := cloner.decorateShadowConfiguration(tt.args.deployment, tt.args.sourceDeployment, tt.args.shadowService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decorateShadowConfiguration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShadowServiceCloner_decorateShadowDeploymentBaseSpec(t *testing.T) {
	type fields struct {
		KubeClient    kubernetes.Interface
		RunTimeClient *client.Client
		CRDClient     *rest.RESTClient
	}
	type args struct {
		deployment       *appsV1.Deployment
		sourceDeployment *appsV1.Deployment
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *appsV1.Deployment
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cloner := &ShadowServiceCloner{
				KubeClient:    tt.fields.KubeClient,
				RunTimeClient: tt.fields.RunTimeClient,
				CRDClient:     tt.fields.CRDClient,
			}
			if got := cloner.decorateShadowDeploymentBaseSpec(tt.args.deployment, tt.args.sourceDeployment); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decorateShadowDeploymentBaseSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShadowServiceCloner_generateShadowDeployment(t *testing.T) {
	type fields struct {
		KubeClient    kubernetes.Interface
		RunTimeClient *client.Client
		CRDClient     *rest.RESTClient
	}
	type args struct {
		sourceDeployment *appsV1.Deployment
		shadowService    *object.ShadowService
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *appsV1.Deployment
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cloner := &ShadowServiceCloner{
				KubeClient:    tt.fields.KubeClient,
				RunTimeClient: tt.fields.RunTimeClient,
				CRDClient:     tt.fields.CRDClient,
			}
			if got := cloner.generateShadowDeployment(tt.args.sourceDeployment, tt.args.shadowService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateShadowDeployment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findContainer(t *testing.T) {
	type args struct {
		containers    []corev1.Container
		containerName string
	}
	tests := []struct {
		name  string
		args  args
		want  *corev1.Container
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := findContainer(tt.args.containers, tt.args.containerName)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findContainer() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("findContainer() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_generateShadowConfigEnv(t *testing.T) {
	type args struct {
		envName string
		config  interface{}
	}
	tests := []struct {
		name string
		args args
		want *corev1.EnvVar
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateShadowConfigEnv(tt.args.envName, tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateShadowConfigEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_injectContainers(t *testing.T) {
	type args struct {
		containers []corev1.Container
		elems      []v1.Container
	}
	tests := []struct {
		name string
		args args
		want []corev1.Container
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := injectContainers(tt.args.containers, tt.args.elems...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("injectContainers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_injectEnvVars(t *testing.T) {
	type args struct {
		envVars []corev1.EnvVar
		elems   []v1.EnvVar
	}
	tests := []struct {
		name string
		args args
		want []corev1.EnvVar
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := injectEnvVars(tt.args.envVars, tt.args.elems...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("injectEnvVars() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_injectShadowAnnotation(t *testing.T) {
	type args struct {
		annotations map[string]string
		service     *object.ShadowService
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_injectShadowLabels(t *testing.T) {
	type args struct {
		labels map[string]string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_shadowConfigurationKeys(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shadowConfigurationKeys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("shadowConfigurationKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_shadowConfigurationMap(t *testing.T) {
	type args struct {
		shadowService *object.ShadowService
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shadowConfigurationMap(tt.args.shadowService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("shadowConfigurationMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_shadowContainer(t *testing.T) {
	type args struct {
		container v1.Container
	}
	tests := []struct {
		name string
		args args
		want v1.Container
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shadowContainer(tt.args.container); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("shadowContainer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_shadowContainers(t *testing.T) {
	type args struct {
		containers []corev1.Container
	}
	tests := []struct {
		name string
		args args
		want []corev1.Container
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shadowContainers(tt.args.containers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("shadowContainers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_shadowInitContainers(t *testing.T) {
	type args struct {
		initContainers []corev1.Container
	}
	tests := []struct {
		name string
		args args
		want []corev1.Container
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shadowInitContainers(tt.args.initContainers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("shadowInitContainers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_shadowName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shadowName(tt.args.name); got != tt.want {
				t.Errorf("shadowName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_shadowServiceLabels(t *testing.T) {
	tests := []struct {
		name string
		want map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shadowServiceLabels(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("shadowServiceLabels() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_shadowVolumes(t *testing.T) {
	type args struct {
		volumes []corev1.Volume
	}
	tests := []struct {
		name string
		args args
		want []corev1.Volume
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shadowVolumes(tt.args.volumes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("shadowVolumes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sourceName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sourceName(tt.args.name); got != tt.want {
				t.Errorf("sourceName() = %v, want %v", got, tt.want)
			}
		})
	}
}
