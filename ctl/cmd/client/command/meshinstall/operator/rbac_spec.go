package operator

import (
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	RoleVerbGet    = "get"
	RoleVerbList   = "list"
	RoleVerbWatch  = "watch"
	RoleVerbCreate = "create"
	RoleVerbUpdate = "update"
	RoleVerbPatch  = "patch"
	RoleVerbDelete = "delete"
)

func roleSpec(args *installbase.InstallArgs) installbase.InstallFunc {

	operatorLeaderElectionRole := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: args.MeshNameSpace,
			Name:      meshOperatorLeaderElectionRole,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"configmaps", "leases"},
				Verbs:     []string{RoleVerbGet, RoleVerbList, RoleVerbWatch, RoleVerbCreate, RoleVerbUpdate, RoleVerbPatch, RoleVerbDelete},
			},
			{
				APIGroups: []string{"", "coordination.k8s.io"},
				Resources: []string{"events"},
				Verbs:     []string{RoleVerbCreate, RoleVerbPatch},
			},
		},
	}

	return func(cmd *cobra.Command, kubeClient *kubernetes.Clientset, args *installbase.InstallArgs) error {
		err := installbase.DeployRole(operatorLeaderElectionRole, kubeClient, args.MeshNameSpace)
		if err != nil {
			return err
		}
		return nil
	}
}

func clusterRoleSpec(args *installbase.InstallArgs) installbase.InstallFunc {
	operatorManagerClusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: meshOperatorManagerClusterRole,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"apps"},
				Resources: []string{"deployments"},
				Verbs:     []string{RoleVerbGet, RoleVerbList, RoleVerbWatch, RoleVerbCreate, RoleVerbUpdate, RoleVerbPatch, RoleVerbDelete},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"pods"},
				Verbs:     []string{RoleVerbGet, RoleVerbList},
			},
			{
				APIGroups: []string{"mesh.megaease.com"},
				Resources: []string{"meshdeployments"},
				Verbs:     []string{RoleVerbGet, RoleVerbList, RoleVerbWatch, RoleVerbCreate, RoleVerbUpdate, RoleVerbPatch, RoleVerbDelete},
			},
			{
				APIGroups: []string{"mesh.megaease.com"},
				Resources: []string{"meshdeployments/finalizers"},
				Verbs:     []string{RoleVerbUpdate},
			},
			{
				APIGroups: []string{"mesh.megaease.com"},
				Resources: []string{"meshdeployments/status"},
				Verbs:     []string{RoleVerbGet, RoleVerbPatch, RoleVerbUpdate},
			},
		},
	}

	metricsReaderClusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: meshOperatorMetricsReaderClusterRole,
		},
		Rules: []rbacv1.PolicyRule{
			{
				NonResourceURLs: []string{"/metrics"},
				Verbs:           []string{RoleVerbGet},
			},
		},
	}

	operatorProxyClusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: meshOperatorProxyClusterRole,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"authentication.k8s.io"},
				Resources: []string{"tokenreviews"},
				Verbs:     []string{RoleVerbCreate},
			},
			{
				APIGroups: []string{"authentication.k8s.io"},
				Resources: []string{"subjectaccessreviews"},
				Verbs:     []string{RoleVerbCreate},
			},
		},
	}

	return func(cmd *cobra.Command, kubeClient *kubernetes.Clientset, args *installbase.InstallArgs) error {
		for _, clusterRole := range []*rbacv1.ClusterRole{operatorManagerClusterRole, metricsReaderClusterRole, operatorProxyClusterRole} {
			err := installbase.DeployClusterRole(clusterRole, kubeClient)
			if err != nil {
				return errors.Wrapf(err, "createClusterRole role %s error", clusterRole.Name)
			}

		}
		return nil
	}
}

func roleBindingSpec(args *installbase.InstallArgs) installbase.InstallFunc {
	operatorLeaderElectionRoleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      meshOperatorLeaderElectionRoleBinding,
			Namespace: args.MeshNameSpace,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     meshOperatorLeaderElectionRole,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "default",
				Namespace: args.MeshNameSpace,
			},
		},
	}

	return func(cmd *cobra.Command, kubeClient *kubernetes.Clientset, args *installbase.InstallArgs) error {
		err := installbase.DeployRoleBinding(operatorLeaderElectionRoleBinding, kubeClient, args.MeshNameSpace)
		if err != nil {
			return err
		}
		return nil
	}
}

func clusterRoleBindingSpec(args *installbase.InstallArgs) installbase.InstallFunc {
	operatorManagerClusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: meshOperatorManagerClusterRoleBinding,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     meshOperatorManagerClusterRole,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "default",
				Namespace: args.MeshNameSpace,
			},
		},
	}

	operatorProxyClusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: meshOperatorProxyClusterRoleBinding,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     meshOperatorProxyClusterRole,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "default",
				Namespace: args.MeshNameSpace,
			},
		},
	}

	operatorMetricsReaderClusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: meshOperatorMetricsReaderClusterRoleBinding,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     meshOperatorMetricsReaderClusterRole,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "default",
				Namespace: args.MeshNameSpace,
			},
		},
	}

	return func(cmd *cobra.Command, kubeClient *kubernetes.Clientset, args *installbase.InstallArgs) error {

		clusterRoleBindings := []*rbacv1.ClusterRoleBinding{
			operatorManagerClusterRoleBinding,
			operatorProxyClusterRoleBinding,
			operatorMetricsReaderClusterRoleBinding,
		}

		for _, clusterRoleBinding := range clusterRoleBindings {
			err := installbase.DeployClusterRoleBinding(clusterRoleBinding, kubeClient)
			if err != nil {
				return errors.Wrapf(err, "Create roleBinding %s error", clusterRoleBinding.Name)
			}
		}
		return nil
	}
}
