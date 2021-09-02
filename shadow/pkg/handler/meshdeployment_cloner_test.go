package handler

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/megaease/easemesh/mesh-shadow/pkg/object/v1beta1"
	k8Yaml "k8s.io/apimachinery/pkg/util/yaml"
)

func TestCloneMeshDeployment(t *testing.T) {

	sourceMeshDeployment := &v1beta1.MeshDeployment{}

	meshData, err := os.ReadFile("./original_meshdeployment.yaml")

	meshDec := k8Yaml.NewYAMLOrJSONDecoder(bytes.NewReader(meshData), 1000)
	err = meshDec.Decode(sourceMeshDeployment)

	// shadowService := object.ShadowService{
	// 	Name: "first-shadow-service",
	// 	NameSpace: "default",
	// 	ServiceName: "visits-service",
	// 	MySQL: &object.MySQL{
	// 		Hosts: []string{"127.0.0.1:3306"},
	// 	},
	// 	ElasticSearch: &object.ElasticSearch{
	// 		Hosts: []string{"127.0.0.1:9200"},
	// 	},
	// 	Kafka: &object.Kafka{
	// 		Hosts: []string{"127.0.0.1:9092"},
	// 	},
	// }

	// CloneMeshDeployment(nil, sourceMeshDeployment, &shadowService)
	if err != nil {
		fmt.Println(err)
	}

}
