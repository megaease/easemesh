package handler

import (
	"testing"

	"github.com/megaease/easemesh/mesh-shadow/pkg/object"
)

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
