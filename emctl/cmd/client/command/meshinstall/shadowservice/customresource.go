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
	"github.com/megaease/easemeshctl/cmd/client/command/meshclient"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"
	"github.com/megaease/easemeshctl/cmd/client/resource"
	"gopkg.in/yaml.v2"
)

const shadowServiceKind = `kind: CustomResourceKind
apiVersion: mesh.megaease.com/v1alpla1
metadata:
  name: ShadowService
spec:
  jsonSchema:
    type: object
    properties:
      name:
        type: string
      namespace:
        type: string
      serviceName:
        type: string
      mysql:
        type: object
        properties:
          uris:
            type: string
          userName:
            type: string
          password:
            type: string
      kafka:
        type: object
        properties:
          uris:
            type: string
      redis:
        type: object
        properties:
          uris:
            type: string
          userName:
            type: string
          password:
            type: string
      rabbitMq:
        type: object
        properties:
          uris:
            type: string
          userName:
            type: string
          password:
            type: string
      elasticSearch:
        type: object
        properties:
          uris:
            type: string
          userName:
            type: string
          password:
            type: string`

func shadowServiceKindSpec(ctx *installbase.StageContext) installbase.InstallFunc {
	return func(ctx *installbase.StageContext) error {
		entrypoints, err := installbase.GetMeshControlPlaneEndpoints(ctx.Client, ctx.Flags.MeshNamespace,
			installbase.ControlPlanePlubicServiceName,
			installbase.ControlPlaneStatefulSetAdminPortName)
		if err != nil {
			return err
		}

		var kind resource.CustomResourceKind
		err = yaml.Unmarshal([]byte(shadowServiceKind), &kind)
		if err != nil {
			return err
		}
		client := meshclient.New(entrypoints[0])
		return client.V1Alpha1().CustomResourceKind().Create(ctx.Cmd.Context(), &kind)
	}
}

func deleteShadowServiceKindSpec(ctx *installbase.StageContext) error {
	entrypoints, err := installbase.GetMeshControlPlaneEndpoints(ctx.Client, ctx.Flags.MeshNamespace,
		installbase.ControlPlanePlubicServiceName,
		installbase.ControlPlaneStatefulSetAdminPortName)
	if err != nil {
		return err
	}

	client := meshclient.New(entrypoints[0])
	return client.V1Alpha1().CustomResourceKind().Delete(ctx.Cmd.Context(), "ShadowService")
}
