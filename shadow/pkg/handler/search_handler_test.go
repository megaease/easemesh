package handler

import (
	"testing"

	v1 "k8s.io/api/apps/v1"
)

func Test_isShadowDeployment(t *testing.T) {
	type args struct {
		spec v1.DeploymentSpec
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isShadowDeployment(tt.args.spec); got != tt.want {
				t.Errorf("isShadowDeployment() = %v, want %v", got, tt.want)
			}
		})
	}
}
