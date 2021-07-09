/*
 * Copyright (c) 2017, MegaEase
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package operator

import (
	"github.com/megaease/easemeshctl/cmd/client/command/flags"
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

func roleSpec(installFlags *flags.Install) installbase.InstallFunc {

	operatorLeaderElectionRole := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: installFlags.MeshNameSpace,
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

	return func(cmd *cobra.Command, kubeClient *kubernetes.Clientset, installFlags *flags.Install) error {
		err := installbase.DeployRole(operatorLeaderElectionRole, kubeClient, installFlags.MeshNameSpace)
		if err != nil {
			return err
		}
		return nil
	}
}

func clusterRoleSpec(installFlags *flags.Install) installbase.InstallFunc {
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

	return func(cmd *cobra.Command, kubeClient *kubernetes.Clientset, installFlags *flags.Install) error {
		for _, clusterRole := range []*rbacv1.ClusterRole{operatorManagerClusterRole, metricsReaderClusterRole, operatorProxyClusterRole} {
			err := installbase.DeployClusterRole(clusterRole, kubeClient)
			if err != nil {
				return errors.Wrapf(err, "createClusterRole role %s error", clusterRole.Name)
			}

		}
		return nil
	}
}

func roleBindingSpec(installFlags *flags.Install) installbase.InstallFunc {
	operatorLeaderElectionRoleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      meshOperatorLeaderElectionRoleBinding,
			Namespace: installFlags.MeshNameSpace,
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
				Namespace: installFlags.MeshNameSpace,
			},
		},
	}

	return func(cmd *cobra.Command, kubeClient *kubernetes.Clientset, installFlags *flags.Install) error {
		err := installbase.DeployRoleBinding(operatorLeaderElectionRoleBinding, kubeClient, installFlags.MeshNameSpace)
		if err != nil {
			return err
		}
		return nil
	}
}

func clusterRoleBindingSpec(installFlags *flags.Install) installbase.InstallFunc {
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
				Namespace: installFlags.MeshNameSpace,
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
				Namespace: installFlags.MeshNameSpace,
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
				Namespace: installFlags.MeshNameSpace,
			},
		},
	}

	return func(cmd *cobra.Command, kubeClient *kubernetes.Clientset, installFlags *flags.Install) error {

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
