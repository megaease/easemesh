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

package operator

import (
	"fmt"
	"strconv"

	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func configMapSpec(ctx *installbase.StageContext) installbase.InstallFunc {
	cfg := installbase.MeshOperatorConfig{
		ImageRegistryURL:          ctx.Flags.ImageRegistryURL,
		ClusterName:               installbase.DefaultMeshControlPlaneName,
		ClusterJoinURLs:           []string{"http://" + flags.DefaultMeshControlPlaneHeadfulServiceName + "." + ctx.Flags.MeshNamespace + ":" + strconv.Itoa(ctx.Flags.EgPeerPort)},
		MetricsAddr:               "127.0.0.1:8080",
		EnableLeaderElection:      false,
		ProbeAddr:                 ":8081",
		WebhookPort:               installbase.DefaultMeshOperatorMutatingWebhookPort,
		CertDir:                   installbase.DefaultMeshOperatorCertDir,
		CertName:                  installbase.DefaultMeshOperatorCertFileName,
		KeyName:                   installbase.DefaultMeshOperatorKeyFileName,
		SidecarImageName:          installbase.DefaultSidecarImageName,
		AgentInitializerImageName: installbase.DefaultEaseagentInitializerImageName,
		Log4jConfigName:           installbase.DefaultLog4jConfigName,
	}

	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      meshOperatorConfigMap,
			Namespace: ctx.Flags.MeshNamespace,
		},
	}
	operatorConfig, err := yaml.Marshal(cfg)
	if err == nil {
		// error will left for high order function to jude
		data := map[string]string{}
		data["operator-config.yaml"] = string(operatorConfig)
		configMap.Data = data
	}

	return func(ctx *installbase.StageContext) error {
		if err != nil {
			return errors.Wrap(err, "ConfigMap build")
		}
		err = installbase.DeployConfigMap(configMap, ctx.Client, ctx.Flags.MeshNamespace)
		if err != nil {
			return fmt.Errorf("create configMap failed: %v ", err)
		}
		return err
	}
}
