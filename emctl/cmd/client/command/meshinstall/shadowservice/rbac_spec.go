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

package shadowservice

import (
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	"github.com/pkg/errors"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func clusterRoleSpec(ctx *installbase.StageContext) installbase.InstallFunc {
	clusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{Name: "namespace-lister"},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"namespaces"},
				Verbs:     []string{"list"},
			},
		},
	}

	return func(ctx *installbase.StageContext) error {
		err := installbase.DeployClusterRole(clusterRole, ctx.Client)
		if err != nil {
			return errors.Wrapf(err, "createClusterRole role %s", clusterRole.Name)
		}
		return nil
	}
}

func clusterRoleBindingSpec(ctx *installbase.StageContext) installbase.InstallFunc {
	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "list-namespaces",
			Namespace: ctx.Flags.MeshNamespace,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "namespace-lister",
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
		err := installbase.DeployClusterRoleBinding(clusterRoleBinding, ctx.Client)
		if err != nil {
			return errors.Wrapf(err, "Create roleBinding %s", clusterRoleBinding.Name)
		}
		return nil
	}
}
