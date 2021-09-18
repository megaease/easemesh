package operator

import (
	"testing"

	"github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base/fake"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestDeployOperatorConfigMap(t *testing.T) {

	client := testclient.NewSimpleClientset()
	stageContext := fake.NewStageContextForApply(client, nil)

	err := configMapSpec(stageContext).Deploy(stageContext)
	if err != nil {
		t.Fatalf("deployment operator configmap err %s", err)
	}

}
