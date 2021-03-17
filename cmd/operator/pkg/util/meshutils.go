package util

import (
	"context"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pkg/errors"
)

type EaseMeshOperator struct {
	*v1.Deployment
}

func (operator EaseMeshOperator) DeepCopyObject() client.Object {
	return nil
}

const (
	EaseMeshOperatorNameSpace string = "mesh-operator-system"
	EaseMeshOperatorName      string = "mesh-operator-controller-manager"
)

func GetEaseMeshOperator(client client.Client) (*EaseMeshOperator, error) {

	deployment := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      EaseMeshOperatorName,
			Namespace: EaseMeshOperatorNameSpace,
		},
	}

	namespacedName := types.NamespacedName{
		Name:      EaseMeshOperatorName,
		Namespace: EaseMeshOperatorNameSpace,
	}
	err := client.Get(context.TODO(), namespacedName, deployment)
	if err != nil {
		return nil, errors.Wrap(err, "EaseMesh Operator not found")
	}

	operator := &EaseMeshOperator{
		deployment,
	}

	return operator, nil
}

func (operator *EaseMeshOperator) GetEGMasterJoinURL(client client.Client) string {

	for name, value := range operator.Labels {
		if name == "cluster-join-url" {
			return value
		}
	}
	return ""
}
