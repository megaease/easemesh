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
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func serviceSpec(ctx *installbase.StageContext) installbase.InstallFunc {
	service := &v1.Service{}
	service.Name = installbase.DefaultMeshIngressService

	service.Spec.Ports = []v1.ServicePort{
		{
			Port:       ctx.Flags.MeshIngressServicePort,
			Protocol:   v1.ProtocolTCP,
			TargetPort: intstr.IntOrString{IntVal: ctx.Flags.MeshIngressServicePort},
		},
	}
	service.Spec.Selector = meshIngressLabel()
	service.Spec.Type = v1.ServiceTypeNodePort
	return func(ctx *installbase.StageContext) error {
		err := installbase.DeployService(service, ctx.Client, ctx.Flags.MeshNamespace)
		return err
	}
}
