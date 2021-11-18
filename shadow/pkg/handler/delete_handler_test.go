package handler

import (
	"testing"

	"github.com/megaease/easemesh/mesh-shadow/pkg/object"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func prepareClientForTest() kubernetes.Interface {
	var result runtime.Object
	namespace := fakeNameSpace()
	deployment := fakeDeployment()
	shadowDeployment := fakeClonedDeployment()

	client := fake.NewSimpleClientset(
		namespace,
		deployment,
		shadowDeployment,
	)
	client.PrependReactor("create", "*", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		result = action.(k8stesting.CreateAction).GetObject()

		return true, action.(k8stesting.CreateAction).GetObject(), k8serr.NewAlreadyExists(schema.GroupResource{
			Resource: "Namespace",
			Group:    "v1",
		}, "na")
	})

	client.PrependReactor("update", "*", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, action.(k8stesting.UpdateAction).GetObject(), nil
	})

	client.PrependReactor("get", "*", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, result, nil
	})

	return client

}

func Test_namespacedName(t *testing.T) {
	type args struct {
		namespace string
		name      string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				name:      "test1",
				namespace: "testns",
			},
			want: "testns" + "/" + "test1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := namespacedName(tt.args.namespace, tt.args.name); got != tt.want {
				t.Errorf("namespacedName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_shadowServiceExists(t *testing.T) {
	type args struct {
		namespacedName       string
		shadowServiceNameMap map[string]object.ShadowService
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test1",
			args: args{
				namespacedName: "testns/test1",
				shadowServiceNameMap: map[string]object.ShadowService{
					"testns/test1": object.ShadowService{
						Name:      "test1",
						Namespace: "testns",
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shadowServiceExists(tt.args.namespacedName, tt.args.shadowServiceNameMap); got != tt.want {
				t.Errorf("shadowServiceExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShadowServiceDeleter_FindDeletableObjs(t *testing.T) {
	deleter := &ShadowServiceDeleter{
		KubeClient:    prepareClientForTest(),
		RunTimeClient: nil,
		DeleteChan:    nil,
	}

	shadowService := fakeShadowService()
	objs := []object.ShadowService{shadowService}
	deleter.FindDeletableObjs(objs)
}
