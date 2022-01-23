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
	"fmt"
	"time"

	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
)

const (
	coreDNSNamespace   = "kube-system"
	coreDNSDeployment  = "coredns"
	coreDNSConfigMap   = "coredns"
	coreDNSClusterRole = "system:coredns"
)

// Deploy deploy resources of coreDNS.
func Deploy(ctx *installbase.StageContext) error {
	err := installbase.BatchDeployResources(ctx,
		[]installbase.InstallFunc{
			configMapSpec(ctx),
			clusterRoleSpec(ctx),

			coreDNSDeploymentSpec(ctx),
		})
	if err != nil {
		return err
	}

	return checkCoreDNSStatus(ctx.Client, ctx.Flags)
}

// PreCheck check prerequisite for installing mesh coreDNS
func PreCheck(context *installbase.StageContext) error {
	// Do nothing
	return nil
}

// Clear clears all k8s resources about coreDNS
func Clear(context *installbase.StageContext) error {
	appsV1Resources := [][]string{
		{"deployments", coreDNSDeployment},
	}

	coreV1Resources := [][]string{
		{"configmaps", coreDNSConfigMap},
	}

	rbacV1Resources := [][]string{
		{"clusterroles", coreDNSClusterRole},
	}

	installbase.DeleteResources(context.Client, appsV1Resources,
		coreDNSNamespace, installbase.DeleteAppsV1Resource)
	installbase.DeleteResources(context.Client, coreV1Resources,
		coreDNSNamespace, installbase.DeleteCoreV1Resource)
	installbase.DeleteResources(context.Client, rbacV1Resources,
		coreDNSNamespace, installbase.DeleteRbacV1Resources)

	return nil
}

// DescribePhase leverage human-readable text to describe different phase
// in the process of the mesh coreDNS
func DescribePhase(context *installbase.StageContext, phase installbase.InstallPhase) string {
	switch phase {
	case installbase.BeginPhase:
		return fmt.Sprintf("Begin to install mesh coreDNS in the namespace: %s", coreDNSNamespace)
	case installbase.EndPhase:
		return fmt.Sprintf("\nMesh coreDNS deployed successfully, deployment: %s\n%s", coreDNSNamespace,
			installbase.FormatPodStatus(context.Client, coreDNSNamespace,
				installbase.AdaptListPodFunc(coreDNSLabels())))
	}
	return ""
}

func checkCoreDNSStatus(client kubernetes.Interface, installFlags *flags.Install) error {
	i := 0
	for {
		time.Sleep(time.Millisecond * 100)
		i++
		ready, err := installbase.CheckDeploymentResourceStatus(client,
			coreDNSNamespace,
			coreDNSDeployment,
			installbase.DeploymentReadyPredict)
		if err != nil {
			return err
		}

		if ready {
			return nil
		}

		// Not ready, retry
		if i > 600 {
			return errors.Errorf("easemesh coreDNS deploy failed, mesh coreDNS (EG deployment) not ready")
		}
	}
}
