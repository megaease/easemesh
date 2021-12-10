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
	"fmt"
	"time"

	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
)

// Deploy deploy resources of shadow service controller
func Deploy(ctx *installbase.StageContext) error {
	err := installbase.BatchDeployResources(ctx, []installbase.InstallFunc{
		clusterRoleSpec(ctx),
		clusterRoleBindingSpec(ctx),
		deploymentSpec(ctx),
		shadowServiceKindSpec(ctx),
	})
	if err != nil {
		return err
	}

	return checkShadowServiceStatus(ctx.Client, ctx.Flags)
}

// PreCheck check prerequisite for installing shadow service controller
func PreCheck(context *installbase.StageContext) error {
	return nil
}

// Clear will clear all installed resource about shadow service controller
func Clear(context *installbase.StageContext) error {
	deleteShadowServiceKindSpec(context)
	appsV1Resources := [][]string{
		{"deployments", installbase.DefaultShadowServiceControllerName},
	}
	installbase.DeleteResources(context.Client, appsV1Resources, context.Flags.MeshNamespace, installbase.DeleteAppsV1Resource)
	return nil
}

// DescribePhase leverage human-readable text to describe different phase
// in the process of the shadow service controller
func DescribePhase(context *installbase.StageContext, phase installbase.InstallPhase) string {
	switch phase {
	case installbase.BeginPhase:
		return fmt.Sprintf("Begin to install shadow service controller in the namespace:%s", context.Flags.MeshNamespace)
	case installbase.EndPhase:
		return fmt.Sprintf("\nShadow service controller deployed successfully, deployment:%s\n%s", installbase.DefaultShadowServiceControllerName,
			installbase.FormatPodStatus(context.Client, context.Flags.MeshNamespace,
				installbase.AdaptListPodFunc(shadowServiceLabel())))
	}
	return ""
}

func checkShadowServiceStatus(client kubernetes.Interface, installFlags *flags.Install) error {
	i := 0
	for {
		time.Sleep(time.Millisecond * 100)
		i++
		if i > 600 {
			return errors.Errorf("easeMesh shadow service controller deploy failed, shadow service controller (EG deployment) not ready")
		}
		ready, err := installbase.CheckDeploymentResourceStatus(client, installFlags.MeshNamespace,
			installbase.DefaultShadowServiceControllerName,
			installbase.DeploymentReadyPredict)
		if ready {
			return nil
		}
		if err != nil {
			return err
		}
	}
}
