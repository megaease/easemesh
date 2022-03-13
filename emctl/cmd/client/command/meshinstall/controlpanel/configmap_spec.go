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
	"fmt"

	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func configMapSpec(ctx *installbase.StageContext) installbase.InstallFunc {
	config := installbase.EasegressConfig{
		// Injected from env EG_NAME
		// Name:                    "" ,

		ClusterName: installbase.ControlPlaneStatefulSetName,
		ClusterRole: installbase.EasegressPrimaryClusterRole,
		Cluster: installbase.ClusterOptions{
			ListenPeerURLs:   []string{fmt.Sprintf("http://0.0.0.0:%d", ctx.Flags.EgPeerPort)},
			ListenClientURLs: []string{fmt.Sprintf("http://0.0.0.0:%d", ctx.Flags.EgClientPort)},

			// Injected from command line.
			// AdvertiseClientURLs: nil,

			// Injected from command line.
			// InitialAdvertisePeerURLs: nil,

			// Injected from command line.
			// InitialCluster: nil,
		},
		APIAddr: fmt.Sprintf("0.0.0.0:%d", ctx.Flags.EgAdminPort),
		HomeDir: installbase.ControlPlaneHomeDir,
		DataDir: installbase.ControlPlaneDataDir,
	}

	yamlBuff, _ := yaml.Marshal(config)
	data := map[string]string{
		installbase.ControlPlaneConfigMapKey: string(yamlBuff),
	}

	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      installbase.ControlPlaneConfigMapName,
			Namespace: ctx.Flags.MeshNamespace,
		},
		Data: data,
	}

	return func(ctx *installbase.StageContext) error {
		err := installbase.DeployConfigMap(configMap, ctx.Client, ctx.Flags.MeshNamespace)
		if err != nil {
			return err
		}
		return nil
	}
}
