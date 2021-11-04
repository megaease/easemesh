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

package controlpanel

import (
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func meshControlPlaneLabel() map[string]string {
	selector := map[string]string{}
	selector["mesh-controlpanel-app"] = "easegress-mesh-controlpanel"
	return selector
}

func serviceSpec(ctx *installbase.StageContext) installbase.InstallFunc {
	labels := meshControlPlaneLabel()

	headlessService := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      installbase.DefaultMeshControlPlaneHeadlessServiceName,
			Namespace: ctx.Flags.MeshNamespace,
		},
	}

	headlessService.Spec.ClusterIP = "None"
	headlessService.Spec.Selector = labels
	headlessService.Spec.Ports = []v1.ServicePort{
		{
			Name:       installbase.DefaultMeshAdminPortName,
			Port:       int32(ctx.Flags.EgAdminPort),
			TargetPort: intstr.IntOrString{IntVal: 2381},
		},
		{
			Name:       installbase.DefaultMeshPeerPortName,
			Port:       int32(ctx.Flags.EgPeerPort),
			TargetPort: intstr.IntOrString{IntVal: 2380},
		},
		{
			Name:       installbase.DefaultMeshClientPortName,
			Port:       int32(ctx.Flags.EgClientPort),
			TargetPort: intstr.IntOrString{IntVal: 2379},
		},
	}

	headfulService := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ctx.Flags.EgServiceName,
			Namespace: ctx.Flags.MeshNamespace,
		},
	}

	headfulService.Spec.Selector = labels
	headfulService.Spec.Ports = []v1.ServicePort{
		{
			Name:       installbase.DefaultMeshAdminPortName,
			Port:       int32(ctx.Flags.EgAdminPort),
			TargetPort: intstr.IntOrString{IntVal: 2381},
		},
		{
			Name:       installbase.DefaultMeshPeerPortName,
			Port:       int32(ctx.Flags.EgPeerPort),
			TargetPort: intstr.IntOrString{IntVal: 2380},
		},
		{
			Name:       installbase.DefaultMeshClientPortName,
			Port:       int32(ctx.Flags.EgClientPort),
			TargetPort: intstr.IntOrString{IntVal: 2379},
		},
	}

	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      installbase.DefaultMeshControlPlanePlubicServiceName,
			Namespace: ctx.Flags.MeshNamespace,
		},
	}
	service.Spec.Ports = []v1.ServicePort{
		{
			Name:       installbase.DefaultMeshAdminPortName,
			Port:       int32(ctx.Flags.EgAdminPort),
			TargetPort: intstr.IntOrString{IntVal: 2381},
		},
		{
			Name:       installbase.DefaultMeshPeerPortName,
			Port:       int32(ctx.Flags.EgPeerPort),
			TargetPort: intstr.IntOrString{IntVal: 2380},
		},
		{
			Name:       installbase.DefaultMeshClientPortName,
			Port:       int32(ctx.Flags.EgClientPort),
			TargetPort: intstr.IntOrString{IntVal: 2379},
		},
	}

	// FIXME: for test we leverage nodeport for expose controlpanel service
	// for production, we will give users options to switch to Loadbalance or ingress
	service.Spec.Type = v1.ServiceTypeNodePort
	service.Spec.Selector = labels

	return func(ctx *installbase.StageContext) error {
		err := installbase.DeployService(headlessService, ctx.Client, ctx.Flags.MeshNamespace)
		if err != nil {
			return errors.Wrap(err, "deploy easemesh controlpanel inner service failed")
		}
		err = installbase.DeployService(service, ctx.Client, ctx.Flags.MeshNamespace)
		if err != nil {
			return errors.Wrap(err, "deploy easemesh controlpanel public service failed")
		}

		err = installbase.DeployService(headfulService, ctx.Client, ctx.Flags.MeshNamespace)
		if err != nil {
			return errors.Wrap(err, "deploy easemesh controlpanel headful service failed")
		}
		return nil
	}
}
