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

package operator

import (
	"fmt"
	"strconv"

	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func configMapSpec(installFlags *flags.Install) installbase.InstallFunc {

	cfg := installbase.MeshOperatorConfig{
		ImageRegistryURL:     installFlags.ImageRegistryURL,
		ClusterName:          installbase.DefaultMeshControlPlaneName,
		ClusterJoinURLs:      "http://" + flags.DefaultMeshControlPlaneHeadfulServiceName + "." + installFlags.MeshNameSpace + ":" + strconv.Itoa(installFlags.EgPeerPort),
		MetricsAddr:          "127.0.0.1:8080",
		EnableLeaderElection: false,
		ProbeAddr:            ":8081",
	}

	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      meshOperatorConfigMap,
			Namespace: installFlags.MeshNameSpace,
		},
	}
	operatorConfig, err := yaml.Marshal(cfg)
	if err == nil {
		// error will left for high order function to jude
		data := map[string]string{}
		data["operator-config.yaml"] = string(operatorConfig)
		configMap.Data = data
	}

	return func(cmd *cobra.Command, client *kubernetes.Clientset, installFlags *flags.Install) error {
		if err != nil {
			return errors.Wrap(err, "ConfigMap build error")
		}
		err = installbase.DeployConfigMap(configMap, client, installFlags.MeshNameSpace)
		if err != nil {
			return fmt.Errorf("create configMap failed: %v ", err)
		}
		return err
	}
}
