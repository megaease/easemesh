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
	DefaultKubernetesConfigDir  = path.Join(homedir.HomeDir(), DefaultKubeDir)
	DefaultKubernetesConfigPath = path.Join(DefaultKubernetesConfigDir, DefaultKubernetesConfig)
)

type EasegressConfig struct {
	Name                    string   `yaml:"name" jsonschema:"required"`
	ClusterName             string   `yaml:"cluster-name" jsonschema:"required"`
	ClusterRole             string   `yaml:"cluster-role" jsonschema:"required"`
	ClusterListenClientUrls []string `yaml:"cluster-listen-client-urls" jsonschema:"required"`
	ClusterListenPeerUrls   []string `yaml:"cluster-listen-peer-urls" jsonschema:"required"`
	ClusterJoinUrls         []string `yaml:"cluster-join-urls" jsonschema:"required"`
	ApiAddr                 string   `yaml:"api-addr" jsonschema:"required"`
	DataDir                 string   `yaml:"data-dir" jsonschema:"required"`
	WalDir                  string   `yaml:"wal-dir" wal-dir:"required"`
	CpuProfileFile          string   `yaml:"cpu-profile-file" jsonschema:"required"`
	MemoryProfileFile       string   `yaml:"memory-profile-file" jsonschema:"required"`
	LogDir                  string   `yaml:"log-dir" jsonschema:"required"`
	MemberDir               string   `yaml:"member-dir" jsonschema:"required"`
	StdLogLevel             string   `yaml:"std-log-level" jsonschema:"required"`
}

type MeshControllerConfig struct {
	Name              string `json:"name" jsonschema:"required"`
	Kind              string `json:"kind" jsonschema:"required"`
	RegistryType      string `json:"registryType" jsonschema:"required"`
	HeartbeatInterval string `json:"heartbeatInterval" jsonschema:"required"`
	IngressPort       int32  `json:"ingressPort" jsonschema:"omitempty"`
}

type MeshOperatorConfig struct {
	ImageRegistryURL     string `yaml:"image-registry-url" jsonschema:"required"`
	ClusterName          string `yaml:"cluster-name" jsonschema:"required"`
	ClusterJoinURLs      string `yaml:"cluster-join-urls" jsonschema:"required"`
	MetricsAddr          string `yaml:"metrics-bind-address" jsonschema:"required"`
	EnableLeaderElection bool   `yaml:"leader-elect" jsonschema:"required"`
	ProbeAddr            string `yaml:"health-probe-bind-address" jsonschema:"required"`
}

type EasegressReaderParams struct {
	ClusterJoinUrls       string            `yaml:"cluster-join-urls" jsonschema:"required"`
	ClusterRequestTimeout string            `yaml:"cluster-request-timeout" jsonschema:"required"`
	ClusterRole           string            `yaml:"cluster-role" jsonschema:"required"`
	ClusterName           string            `yaml:"cluster-name" jsonschema:"required"`
	Name                  string            `yaml:"name" jsonschema:"required"`
	Labels                map[string]string `yaml:"labels" jsonschema:"required"`
}
type StageContext struct {
	Cmd                 *cobra.Command
	Client              *kubernetes.Clientset
	Flags               *flags.Install
	APIExtensionsClient *apiextensions.Clientset
	ClearFuncs          []func(*StageContext) error
}

type InstallFunc func(*cobra.Command, *kubernetes.Clientset, *flags.Install) error

func (i InstallFunc) Deploy(cmd *cobra.Command, c *kubernetes.Clientset, flags *flags.Install) error {
	return i(cmd, c, flags)
}
