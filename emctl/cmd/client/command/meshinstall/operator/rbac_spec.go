/*
 * Copyright (c) 2021, MegaEase
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
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	"github.com/pkg/errors"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	roleVerbGet    = "get"
	roleVerbList   = "list"
	roleVerbWatch  = "watch"
	roleVerbCreate = "create"
	roleVerbUpdate = "update"
	roleVerbPatch  = "patch"
	roleVerbDelete = "delete"
)

func roleSpec(ctx *installbase.StageContext) installbase.InstallFunc {
	operatorLeaderElectionRole := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ctx.Flags.MeshNamespace,
			Name:      leaderElectionRole,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"configmaps", "leases"},
				Verbs:     []string{roleVerbGet, roleVerbList, roleVerbWatch, roleVerbCreate, roleVerbUpdate, roleVerbPatch, roleVerbDelete},
			},
			{
				APIGroups: []string{"", "coordination.k8s.io"},
				Resources: []string{"events"},
				Verbs:     []string{roleVerbCreate, roleVerbPatch},
			},
		},
	}

	return func(ctx *installbase.StageContext) error {
		return installbase.DeployRole(operatorLeaderElectionRole, ctx.Client, ctx.Flags.MeshNamespace)
	}
}

func clusterRoleSpec(ctx *installbase.StageContext) installbase.InstallFunc {
	operatorManagerClusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: managerClusterRole,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"apps"},
				Resources: []string{"deployments"},
				Verbs:     []string{roleVerbGet, roleVerbList, roleVerbWatch, roleVerbCreate, roleVerbUpdate, roleVerbPatch, roleVerbDelete},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"pods"},
				Verbs:     []string{roleVerbGet, roleVerbList},
			},
			{
				APIGroups: []string{"mesh.megaease.com"},
				Resources: []string{"meshdeployments"},
				Verbs:     []string{roleVerbGet, roleVerbList, roleVerbWatch, roleVerbCreate, roleVerbUpdate, roleVerbPatch, roleVerbDelete},
			},
			{
				APIGroups: []string{"mesh.megaease.com"},
				Resources: []string{"meshdeployments/finalizers"},
				Verbs:     []string{roleVerbUpdate},
			},
			{
				APIGroups: []string{"mesh.megaease.com"},
				Resources: []string{"meshdeployments/status"},
				Verbs:     []string{roleVerbGet, roleVerbPatch, roleVerbUpdate},
			},
		},
	}

	metricsReaderClusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: metricsReaderClusterRole,
		},
		Rules: []rbacv1.PolicyRule{
			{
				NonResourceURLs: []string{"/metrics"},
				Verbs:           []string{roleVerbGet},
			},
		},
	}

	operatorProxyClusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: proxyClusterRole,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"authentication.k8s.io"},
				Resources: []string{"tokenreviews"},
				Verbs:     []string{roleVerbCreate},
			},
			{
				APIGroups: []string{"authentication.k8s.io"},
				Resources: []string{"subjectaccessreviews"},
				Verbs:     []string{roleVerbCreate},
			},
		},
	}

	return func(ctx *installbase.StageContext) error {
		for _, clusterRole := range []*rbacv1.ClusterRole{operatorManagerClusterRole, metricsReaderClusterRole, operatorProxyClusterRole} {
			err := installbase.DeployClusterRole(clusterRole, ctx.Client)
			if err != nil {
				return errors.Wrapf(err, "createClusterRole role %s", clusterRole.Name)
			}

		}
		return nil
	}
}

func roleBindingSpec(ctx *installbase.StageContext) installbase.InstallFunc {
	operatorLeaderElectionRoleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      leaderElectionRoleBinding,
			Namespace: ctx.Flags.MeshNamespace,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     leaderElectionRole,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "default",
				Namespace: ctx.Flags.MeshNamespace,
			},
		},
	}

	return func(ctx *installbase.StageContext) error {
		return installbase.DeployRoleBinding(operatorLeaderElectionRoleBinding, ctx.Client, ctx.Flags.MeshNamespace)
	}
}

func clusterRoleBindingSpec(ctx *installbase.StageContext) installbase.InstallFunc {
	operatorManagerClusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: managerClusterRoleBinding,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     managerClusterRole,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "default",
				Namespace: ctx.Flags.MeshNamespace,
			},
		},
	}

	operatorProxyClusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: proxyClusterRoleBinding,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     proxyClusterRole,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "default",
				Namespace: ctx.Flags.MeshNamespace,
			},
		},
	}

	operatorMetricsReaderClusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: metricsReaderClusterRoleBinding,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     metricsReaderClusterRole,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "default",
				Namespace: ctx.Flags.MeshNamespace,
			},
		},
	}

	return func(ctx *installbase.StageContext) error {
		clusterRoleBindings := []*rbacv1.ClusterRoleBinding{
			operatorManagerClusterRoleBinding,
			operatorProxyClusterRoleBinding,
			operatorMetricsReaderClusterRoleBinding,
		}

		for _, clusterRoleBinding := range clusterRoleBindings {
			err := installbase.DeployClusterRoleBinding(clusterRoleBinding, ctx.Client)
			if err != nil {
				return errors.Wrapf(err, "Create roleBinding %s", clusterRoleBinding.Name)
			}
		}
		return nil
	}
}
