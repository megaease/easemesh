package utils

import (
	"reflect"
	"testing"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func TestDeleteDeployment(t *testing.T) {
	type args struct {
		namespace string
		name      string
		clientSet kubernetes.Interface
		options   v1.DeleteOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteDeployment(tt.args.namespace, tt.args.name, tt.args.clientSet, tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("DeleteDeployment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestListDeployments(t *testing.T) {
	type args struct {
		namespace string
		clientSet kubernetes.Interface
		options   v1.ListOptions
	}
	tests := []struct {
		name    string
		args    args
		want    []appsV1.Deployment
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListDeployments(tt.args.namespace, tt.args.clientSet, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListDeployments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListDeployments() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListNameSpaces(t *testing.T) {
	type args struct {
		clientSet kubernetes.Interface
	}
	tests := []struct {
		name    string
		args    args
		want    []corev1.Namespace
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListNameSpaces(tt.args.clientSet)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListNameSpaces() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListNameSpaces() got = %v, want %v", got, tt.want)
			}
		})
	}
}
