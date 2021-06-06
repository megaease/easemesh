package operator

import (
	"fmt"
	"time"

	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
)

const (
	meshOperatorConfigMap = "easemesh-operator-config"
	//
	meshOperatorLeaderElectionRole        = "mesh-operator-leader-election-role"
	meshOperatorLeaderElectionRoleBinding = "mesh-operator-leader-election-rolebinding"
	//
	meshOperatorManagerClusterRole        = "mesh-operator-manager-role"
	meshOperatorManagerClusterRoleBinding = "mesh-operator-manager-rolebinding"

	//
	meshOperatorMetricsReaderClusterRole        = "mesh-operator-metrics-reader-role"
	meshOperatorMetricsReaderClusterRoleBinding = "mesh-operator-metrics-reader-rolebinding"

	//
	meshOperatorProxyClusterRole        = "mesh-operator-proxy-role"
	meshOperatorProxyClusterRoleBinding = "mesh-operator-proxy-rolebinding"
)

// Deploy deploy resources of operator
func Deploy(context *installbase.StageContext) error {
	err := installbase.BatchDeployResources(context.Cmd, context.Client, &context.Arguments, []installbase.InstallFunc{
		configMapSpec(&context.Arguments),
		serviceSpec(&context.Arguments),
		roleSpec(&context.Arguments),
		clusterRoleSpec(&context.Arguments),
		roleBindingSpec(&context.Arguments),
		clusterRoleBindingSpec(&context.Arguments),
		operatorDeploymentSpec(&context.Arguments),
	})
	if err != nil {
		return err
	}

	return checkOperatorStatus(context.Client, &context.Arguments)
}

// PreCheck check prerequisite for installing mesh operator
func PreCheck(context *installbase.StageContext) error {
	// Do nothing
	return nil
}

// Clear clears all k8s resources about operator
func Clear(context *installbase.StageContext) error {

	appsV1Resources := [][]string{
		{"deployments", installbase.DefaultMeshOperatorName},
	}

	coreV1Resources := [][]string{
		{"services", installbase.DefaultMeshOperatorControllerManagerServiceName},
		{"configmap", meshOperatorConfigMap},
	}

	rbacV1Resources := [][]string{
		{"rolebindings", meshOperatorLeaderElectionRoleBinding},
		{"roles", meshOperatorLeaderElectionRole},
		{"clusterrolebindings", meshOperatorManagerClusterRoleBinding},
		{"clusterroles", meshOperatorManagerClusterRole},
		{"clusterrolebindings", meshOperatorMetricsReaderClusterRoleBinding},
		{"clusterroles", meshOperatorMetricsReaderClusterRole},
		{"clusterrolebindings", meshOperatorProxyClusterRoleBinding},
		{"clusterroles", meshOperatorProxyClusterRole},
	}
	installbase.DeleteResources(context.Client, appsV1Resources, context.Arguments.MeshNameSpace, installbase.DeleteAppsV1Resource)
	installbase.DeleteResources(context.Client, coreV1Resources, context.Arguments.MeshNameSpace, installbase.DeleteCoreV1Resource)
	installbase.DeleteResources(context.Client, rbacV1Resources, context.Arguments.MeshNameSpace, installbase.DeleteRbacV1Resources)

	return nil
}

// Describe leverage human-readable text to describe different phase
// in the process of the mesh operator
func Describe(context *installbase.StageContext, phase installbase.InstallPhase) string {
	switch phase {
	case installbase.BeginPhase:
		return fmt.Sprintf("Begin to install mesh operator in the namespace: %s", context.Arguments.MeshNameSpace)
	case installbase.EndPhase:
		return fmt.Sprintf("\nMesh operator deployed successfully, deployment: %s\n%s", installbase.DefaultMeshOperatorName,
			installbase.FormatPodStatus(context.Client, context.Arguments.MeshNameSpace,
				installbase.AdaptListPodFunc(meshOperatorLabels())))
	}
	return ""
}

func checkOperatorStatus(client *kubernetes.Clientset, args *installbase.InstallArgs) error {
	i := 0
	for {
		time.Sleep(time.Millisecond * 100)
		i++
		if i > 600 {
			return errors.Errorf("easemesh operator deploy failed, mesh operator (EG deployment) not ready")
		}
		ready, err := installbase.CheckDeploymentResourceStatus(client, args.MeshNameSpace,
			installbase.DefaultMeshOperatorName,
			installbase.DeploymentReadyPredict)
		if ready {
			return nil
		}
		if err != nil {
			return err
		}
	}
}
