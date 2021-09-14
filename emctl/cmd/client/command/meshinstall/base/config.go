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

package installbase

import (
	"path"

	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	"github.com/spf13/cobra"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/homedir"
)

var (
	// DefaultKubernetesConfigDir is the default kubernetes config directory.
	DefaultKubernetesConfigDir = path.Join(homedir.HomeDir(), DefaultKubeDir)

	// DefaultKubernetesConfigPath is the default kubernetes config path.
	DefaultKubernetesConfigPath = path.Join(DefaultKubernetesConfigDir, DefaultKubernetesConfig)
)

type (
	// EasegressConfig is the config of Easegress.
	EasegressConfig struct {
		Name                    string   `yaml:"name" jsonschema:"required"`
		ClusterName             string   `yaml:"cluster-name" jsonschema:"required"`
		ClusterRole             string   `yaml:"cluster-role" jsonschema:"required"`
		ClusterListenClientUrls []string `yaml:"cluster-listen-client-urls" jsonschema:"required"`
		ClusterListenPeerUrls   []string `yaml:"cluster-listen-peer-urls" jsonschema:"required"`
		ClusterJoinUrls         []string `yaml:"cluster-join-urls" jsonschema:"required"`
		APIAddr                 string   `yaml:"api-addr" jsonschema:"required"`
		DataDir                 string   `yaml:"data-dir" jsonschema:"required"`
		WalDir                  string   `yaml:"wal-dir" wal-dir:"required"`
		CPUProfileFile          string   `yaml:"cpu-profile-file" jsonschema:"required"`
		MemoryProfileFile       string   `yaml:"memory-profile-file" jsonschema:"required"`
		LogDir                  string   `yaml:"log-dir" jsonschema:"required"`
		MemberDir               string   `yaml:"member-dir" jsonschema:"required"`
		StdLogLevel             string   `yaml:"std-log-level" jsonschema:"required"`
	}

	// MeshControllerConfig is the config of EaseMesh Controller.
	MeshControllerConfig struct {
		Name              string `json:"name" jsonschema:"required"`
		Kind              string `json:"kind" jsonschema:"required"`
		RegistryType      string `json:"registryType" jsonschema:"required"`
		HeartbeatInterval string `json:"heartbeatInterval" jsonschema:"required"`
		IngressPort       int32  `json:"ingressPort" jsonschema:"omitempty"`
	}

	// MeshOperatorConfig is the config of EaseMesh operator.
	MeshOperatorConfig struct {
		ImageRegistryURL     string   `yaml:"image-registry-url" jsonschema:"required"`
		ClusterName          string   `yaml:"cluster-name" jsonschema:"required"`
		ClusterJoinURLs      []string `yaml:"cluster-join-urls" jsonschema:"required"`
		MetricsAddr          string   `yaml:"metrics-bind-address" jsonschema:"required"`
		EnableLeaderElection bool     `yaml:"leader-elect" jsonschema:"required"`
		ProbeAddr            string   `yaml:"health-probe-bind-address" jsonschema:"required"`
		WebhookPort          uint16   `yaml:"webhook-port" jsonschema:"required"`
		CertDir              string   `yaml:"cert-dir" jsonschema:"required"`
		CertName             string   `yaml:"cert-name" jsonschema:"required"`
		KeyName              string   `yaml:"key-name" jsonschema:"required"`
		// The image name of the injecting sidecar
		SidecarImageName string `yaml:"sidecar-image-name" jsonschema:"required"`

		// The image name of the easeagent initializer
		AgentInitializerImageName string `yaml:"agent-initializer-image-name" jsonschema:"required"`
		// Log4jConfigName default is easeagent-log4j.xml
		Log4jConfigName string `yaml:"log4j-config-name" jsonschema:"required"`
	}

	// EasegressReaderParams is the parameters of Easegress reader role.
	EasegressReaderParams struct {
		ClusterJoinUrls       string            `yaml:"cluster-join-urls" jsonschema:"required"`
		ClusterRequestTimeout string            `yaml:"cluster-request-timeout" jsonschema:"required"`
		ClusterRole           string            `yaml:"cluster-role" jsonschema:"required"`
		ClusterName           string            `yaml:"cluster-name" jsonschema:"required"`
		Name                  string            `yaml:"name" jsonschema:"required"`
		Labels                map[string]string `yaml:"labels" jsonschema:"required"`
	}

	// StageContext is the context for every installation stage.
	StageContext struct {
		Cmd                 *cobra.Command
		Client              *kubernetes.Clientset
		Flags               *flags.Install
		APIExtensionsClient *apiextensions.Clientset
		ClearFuncs          []func(*StageContext) error
	}

	// InstallFunc is the type of function for installation.
	InstallFunc func(ctx *StageContext) error
)

// Deploy executes the install function.
func (fn InstallFunc) Deploy(ctx *StageContext) error {
	return fn(ctx)
}
