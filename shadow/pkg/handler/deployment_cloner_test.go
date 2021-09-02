package handler

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/megaease/easemesh/mesh-shadow/pkg/object"
	"github.com/megaease/easemesh/mesh-shadow/pkg/object/v1beta1"
	appv1 "k8s.io/api/apps/v1"
	k8Yaml "k8s.io/apimachinery/pkg/util/yaml"
)

func TestCloneDeploymentSpec(t *testing.T) {

	data, err := os.ReadFile("./original_deployment.yaml")

	sourceDeployment := &appv1.Deployment{}
	dec := k8Yaml.NewYAMLOrJSONDecoder(bytes.NewReader(data), 1000)
	err = dec.Decode(sourceDeployment)

	shadowService := object.ShadowService{
		Name: "first-shadow-service",
		NameSpace: "default",
		ServiceName: "visits-service",
		MySQL: &object.MySQL{
			Hosts: []string{"127.0.0.1:3306"},
		},
		ElasticSearch: &object.ElasticSearch{
			Hosts: []string{"127.0.0.1:9200"},
		},
		Kafka: &object.Kafka{
			Hosts: []string{"127.0.0.1:9092"},
		},
	}
	if err != nil {
		fmt.Println(shadowService, err)
	}


	sourceMeshDeployment := &v1beta1.MeshDeployment{}

	meshData, err := os.ReadFile("./original_meshdeployment.yaml")

	meshDec := k8Yaml.NewYAMLOrJSONDecoder(bytes.NewReader(meshData), 1000)
	err = meshDec.Decode(sourceMeshDeployment)

	if err != nil {
		fmt.Println(err)
	}

}
