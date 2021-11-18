package handler

import (
	"testing"

	"k8s.io/client-go/kubernetes"
)

func TestShadowServiceCloner_Clone(t *testing.T) {
	type fields struct {
		KubeClient    kubernetes.Interface
		RunTimeClient *client.Client
		CRDClient     *rest.RESTClient
	}
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
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
		})
	}
}
