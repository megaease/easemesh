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
	"encoding/json"
	"fmt"
	"strconv"

	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	yamljsontool "github.com/ghodss/yaml"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func configMapSpec(ctx *installbase.StageContext) installbase.InstallFunc {
	host := "0.0.0.0"

	config := installbase.EasegressConfig{
		Name:                    installbase.DefaultMeshControlPlaneName,
		ClusterName:             installbase.DefaultMeshControlPlaneName,
		ClusterRole:             installbase.WriterClusterRole,
		ClusterListenClientUrls: []string{"http://" + "0.0.0.0:" + strconv.Itoa(ctx.Flags.EgClientPort)},
		ClusterListenPeerUrls:   []string{"http://" + "0.0.0.0:" + strconv.Itoa(ctx.Flags.EgPeerPort)},
		ClusterJoinUrls:         []string{},
		APIAddr:                 host + ":" + strconv.Itoa(ctx.Flags.EgAdminPort),
		DataDir:                 "/opt/eg-data/data",
		WalDir:                  "",
		CPUProfileFile:          "",
		MemoryProfileFile:       "",
		LogDir:                  "/opt/eg-data/log",
		MemberDir:               "/opt/eg-data/member",
		StdLogLevel:             "INFO",
	}

	for i := 0; i < ctx.Flags.EasegressControlPlaneReplicas; i++ {
		config.ClusterJoinUrls = append(config.ClusterJoinUrls,
			fmt.Sprintf("http://%s-%d.%s.%s:%d",
				installbase.DefaultMeshControlPlaneName,
				i,
				installbase.DefaultMeshControlPlaneHeadlessServiceName,
				ctx.Flags.MeshNamespace,
				ctx.Flags.EgPeerPort))
	}

	configData := map[string]string{}
	configBytes, _ := yaml.Marshal(config)
	configData["eg-master.yaml"] = string(configBytes)

	buff, _ := yaml.Marshal(configData)
	configJSON, _ := yamljsontool.YAMLToJSON(buff)

	var params map[string]string
	_ = json.Unmarshal(configJSON, &params)

	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      installbase.DefaultMeshControlPlaneConfig,
			Namespace: ctx.Flags.MeshNamespace,
		},
		Data: params,
	}

	return func(ctx *installbase.StageContext) error {
		err := installbase.DeployConfigMap(configMap, ctx.Client, ctx.Flags.MeshNamespace)
		if err != nil {
			return err
		}
		return nil
	}
}
