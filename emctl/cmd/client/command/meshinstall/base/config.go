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

package installbase

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	"github.com/spf13/cobra"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	// DefaultKubernetesConfigDir is the default kubernetes config directory.
	DefaultKubernetesConfigDir = path.Join(homedir.HomeDir(), DefaultKubeDir)

	// DefaultKubernetesConfigPath is the default kubernetes config path.
	DefaultKubernetesConfigPath = path.Join(DefaultKubernetesConfigDir, DefaultKubernetesFilename)
)

type (
	// EasegressConfig is the config of Easegress.
	EasegressConfig struct {
		// meta
		Name    string            `yaml:"name"`
		Labels  map[string]string `yaml:"labels"`
		APIAddr string            `yaml:"api-addr"`
		Debug   bool              `yaml:"debug,omitempty"`

		// cluster options
		ClusterName string         `yaml:"cluster-name"`
		ClusterRole string         `yaml:"cluster-role"`
		Cluster     ClusterOptions `yaml:"cluster"`

		// Path.
		HomeDir   string `yaml:"home-dir,omitempty"`
		DataDir   string `yaml:"data-dir,omitempty"`
		WALDir    string `yaml:"wal-dir,omitempty"`
		LogDir    string `yaml:"log-dir,omitempty"`
		MemberDir string `yaml:"member-dir,omitempty"`

		// Profile.
		CPUProfileFile    string `yaml:"cpu-profile-file,omitempty"`
		MemoryProfileFile string `yaml:"memory-profile-file,omitempty"`
	}

	// ClusterOptions is the start-up options of Easegress cluster.
	ClusterOptions struct {
		// Primary members define following URLs to form a cluster.
		ListenPeerURLs           []string          `yaml:"listen-peer-urls"`
		ListenClientURLs         []string          `yaml:"listen-client-urls"`
		AdvertiseClientURLs      []string          `yaml:"advertise-client-urls"`
		InitialAdvertisePeerURLs []string          `yaml:"initial-advertise-peer-urls"`
		InitialCluster           map[string]string `yaml:"initial-cluster"`
		StateFlag                string            `yaml:"state-flag"`

		// Secondary members define URLs to connect to cluster formed by primary members.
		PrimaryListenPeerURLs []string `yaml:"primary-listen-peer-urls"`
		MaxCallSendMsgSize    int      `yaml:"max-call-send-msg-size"`
	}

	// MeshControllerConfig is the config of EaseMesh Controller.
	MeshControllerConfig struct {
		Name              string `yaml:"name" jsonschema:"required"`
		Kind              string `yaml:"kind" jsonschema:"required"`
		RegistryType      string `yaml:"registryType" jsonschema:"required"`
		HeartbeatInterval string `yaml:"heartbeatInterval" jsonschema:"required"`
		IngressPort       int32  `yaml:"ingressPort" jsonschema:"omitempty"`
		APIPort           int    `yaml:"apiPort" jsonschema:"required"`
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
		Client              kubernetes.Interface
		ClientConfig        clientcmd.ClientConfig
		Flags               *flags.Install
		CoreDNSFlags        *flags.CoreDNS
		APIExtensionsClient apiextensions.Interface
		ClearFuncs          []func(*StageContext) error
	}

	// InstallFunc is the type of function for installation.
	InstallFunc func(ctx *StageContext) error
)

// Deploy executes the install function.
func (fn InstallFunc) Deploy(ctx *StageContext) error {
	return fn(ctx)
}

// ControlPlanePodName returns the pod name of control plane.
func ControlPlanePodName(index int) string {
	return fmt.Sprintf("%s-%d", ControlPlaneStatefulSetName, index)
}

// ControlPlanePodAdvertiseClientURL returns the advertise URL of pod of control plane.
func ControlPlanePodAdvertiseClientURL(podName string, ctx *StageContext) string {
	clientPort := ctx.Flags.EgClientPort
	namespace := ctx.Flags.MeshNamespace

	return fmt.Sprintf("http://%s.%s.%s:%d", podName,
		ControlPlaneHeadlessServiceName, namespace, clientPort)
}

// ControlPlanePodAdvertisePeerURL returns the advertise URL of pod of control plane.
func ControlPlanePodAdvertisePeerURL(podName string, ctx *StageContext) string {
	peerPort := ctx.Flags.EgPeerPort
	namespace := ctx.Flags.MeshNamespace

	return fmt.Sprintf("http://%s.%s.%s:%d", podName,
		ControlPlaneHeadlessServiceName, namespace, peerPort)
}

// ControlPlaneInitialCluster returns initial cluster of control plane.
func ControlPlaneInitialCluster(ctx *StageContext) map[string]string {
	replicas := ctx.Flags.EasegressControlPlaneReplicas

	initCluster := map[string]string{}
	for i := 0; i < replicas; i++ {
		podName := ControlPlanePodName(i)
		initCluster[podName] = ControlPlanePodAdvertisePeerURL(podName, ctx)
	}

	return initCluster
}

// ControlPlaneInitialClusterStr returns initial cluster in string of control plane.
func ControlPlaneInitialClusterStr(ctx *StageContext) string {
	initCluster := ControlPlaneInitialCluster(ctx)
	initClusterSlice := []string{}
	for k, v := range initCluster {
		initClusterSlice = append(initClusterSlice, fmt.Sprintf("%s=%s", k, v))
	}

	sort.Strings(initClusterSlice)

	return strings.Join(initClusterSlice, ",")
}

// ControlPlanePeerURLs returns peer URLs of control plane.
func ControlPlanePeerURLs(ctx *StageContext) []string {
	initCluster := ControlPlaneInitialCluster(ctx)
	peerURLs := []string{}
	for _, peerURL := range initCluster {
		peerURLs = append(peerURLs, peerURL)
	}

	return peerURLs
}

// ControlPlanePeerURLsStr returns peer URLs in string of control plane.
func ControlPlanePeerURLsStr(ctx *StageContext) string {
	peerURLs := ControlPlanePeerURLs(ctx)
	return strings.Join(peerURLs, ",")
}
