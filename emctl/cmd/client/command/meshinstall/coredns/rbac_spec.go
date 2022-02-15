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

package coredns

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

func clusterRoleSpec(ctx *installbase.StageContext) installbase.InstallFunc {
	coreDNSClusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: coreDNSClusterRole,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"endpoints", "services", "pods", "namespaces"},
				Verbs:     []string{roleVerbGet, roleVerbList, roleVerbWatch},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"nodes"},
				Verbs:     []string{roleVerbGet},
			},
			{
				APIGroups: []string{"discovery.k8s.io"},
				Resources: []string{"endpointslices"},
				Verbs:     []string{roleVerbGet, roleVerbList, roleVerbWatch},
			},
		},
	}

	return func(ctx *installbase.StageContext) error {
		err := installbase.DeployClusterRole(coreDNSClusterRole, ctx.Client)
		if err != nil {
			return errors.Wrapf(err, "deploty ClusterRole %s failed", coreDNSClusterRole.Name)
		}

		return nil
	}
}
