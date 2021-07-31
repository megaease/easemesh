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

package meshingress

import (
	"fmt"
	"time"

	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
)

// Deploy deploy resources of mesh ingress controller
func Deploy(ctx *installbase.StageContext) error {
	err := installbase.BatchDeployResources(ctx, []installbase.InstallFunc{
		configMapSpec(ctx),
		serviceSpec(ctx),
		deploymentSpec(ctx),
	})
	if err != nil {
		return err
	}

	return checkMeshIngressStatus(ctx.Client, ctx.Flags)
}

// PreCheck check prerequisite for installing mesh ingress controller
func PreCheck(context *installbase.StageContext) error {
	return nil
}

// Clear will clear all installed resource about mesh ingress panel
func Clear(context *installbase.StageContext) error {
	appsV1Resources := [][]string{
		{"deployments", installbase.DefaultMeshIngressControllerName},
	}
	coreV1Resources := [][]string{
		{"services", installbase.DefaultMeshIngressService},
		{"configmap", installbase.DefaultMeshIngressConfig},
	}

	installbase.DeleteResources(context.Client, appsV1Resources, context.Flags.MeshNamespace, installbase.DeleteAppsV1Resource)
	installbase.DeleteResources(context.Client, coreV1Resources, context.Flags.MeshNamespace, installbase.DeleteCoreV1Resource)
	return nil
}

// DescribePhase leverage human-readable text to describe different phase
// in the process of the mesh ingress controller
func DescribePhase(context *installbase.StageContext, phase installbase.InstallPhase) string {
	switch phase {
	case installbase.BeginPhase:
		return fmt.Sprintf("Begin to install mesh ingress controller in the namespace:%s", context.Flags.MeshNamespace)
	case installbase.EndPhase:
		return fmt.Sprintf("\nMesh ingress controller deployed successfully, deployment:%s\n%s", installbase.DefaultMeshIngressControllerName,
			installbase.FormatPodStatus(context.Client, context.Flags.MeshNamespace,
				installbase.AdaptListPodFunc(meshIngressLabel())))
	}
	return ""
}

func checkMeshIngressStatus(client *kubernetes.Clientset, installFlags *flags.Install) error {
	i := 0
	for {
		time.Sleep(time.Millisecond * 100)
		i++
		if i > 600 {
			return errors.Errorf("easeMesh meshingress controller deploy failed, mesh ingress controller (EG deployment) not ready")
		}
		ready, err := installbase.CheckDeploymentResourceStatus(client, installFlags.MeshNamespace,
			installbase.DefaultMeshIngressControllerName,
			installbase.DeploymentReadyPredict)
		if ready {
			return nil
		}
		if err != nil {
			return err
		}
	}
}
