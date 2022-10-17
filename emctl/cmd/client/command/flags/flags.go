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

package flags

import (
	"time"

	"github.com/megaease/easemeshctl/cmd/client/command/rcfile"
	"github.com/megaease/easemeshctl/cmd/common"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
)

const (
	// DefaultMeshNamespace all mesh infrastructure component should be deployed in this namespace
	DefaultMeshNamespace = "easemesh"

	// DefaultMeshControlPlaneReplicas is default number of the control plane service's replicas
	DefaultMeshControlPlaneReplicas = 3

	// DefaultMeshIngressReplicas is default number of the mesh ingress service's replicas
	DefaultMeshIngressReplicas = 1

	// DefaultMeshOperatorReplicas is default number of the operator's  replicas
	DefaultMeshOperatorReplicas = 1

	// DefaultMeshClientPort is the default port interacted with the mesh control plane service
	DefaultMeshClientPort = 2379

	// DefaultMeshPeerPort is the default port with which control plane service interact each other
	DefaultMeshPeerPort = 2380

	// DefaultMeshAdminPort is the default administrator port of control plane service
	DefaultMeshAdminPort = 2381

	// DefaultMeshControlPlaneHeadfulServiceName is the default headful service name of the EaseMesh control plane
	DefaultMeshControlPlaneHeadfulServiceName = "easemesh-control-plane-service"

	// DefaultMeshControlPlaneCheckHealthzMaxTime is a duration that installation wait for checking the control plane service status
	DefaultMeshControlPlaneCheckHealthzMaxTime = 60

	// DefaultMeshControlPlaneStorageClassName is a storage class name of persistent volume used by control plane service
	DefaultMeshControlPlaneStorageClassName = "easemesh-storage"

	// DefaultMeshControlPlanePersistVolumeCapacity is the default capacity of persistent volume needed by control plane service
	DefaultMeshControlPlanePersistVolumeCapacity = "3Gi" // 3 Gib

	// DefaultMeshRegistryType is default registry type of the EaseMesh
	DefaultMeshRegistryType = "eureka"

	// DefaultHeartbeatInterval is default heartbeat
	DefaultHeartbeatInterval = 5

	// MeshControllerKind is kind of the EaseMesh controller in the Easegress
	MeshControllerKind = "MeshController"

	// DefaultMeshIngressServicePort is default port listened by the Easegress acted as an ingress role
	DefaultMeshIngressServicePort = 19527

	// DefaultWaitControlPlaneSeconds is the default wait control plane ready elapse, in seconds (intall command)
	DefaultWaitControlPlaneSeconds = 3

	// MeshControlPlanePVNameHelpStr is a text described name of persistent volume
	MeshControlPlanePVNameHelpStr = "The name of PersistentVolume for EaseMesh control plane storage"
	// MeshControlPlanePVHostPathHelpStr is a text described local path of persistent volume
	MeshControlPlanePVHostPathHelpStr = "The host path of the PersistentVolume for EaseMesh control plane storage"
	// MeshControlPlanePVCapacityHelpStr is a text described capacity of persistent volume
	MeshControlPlanePVCapacityHelpStr = "The capacity of the PersistentVolume for EaseMesh control plane storage"

	// MeshRegistryTypeHelpStr is a text described registry type of persistent volume
	MeshRegistryTypeHelpStr = "The registry type for application service registry, support eureka, consul, nacos"

	// MeshControlPlaneStartupFailedHelpStr is a text described the failure when control plane service started
	MeshControlPlaneStartupFailedHelpStr = `
		EaseMesh Control Plane deploy failed. Please check the K8S resource under the %s namespace for finding errors as follows:

		$ kubectl get statefulsets.apps -n %s

		$ kubectl get pods -n %s
	`
	// MeshControlPlanePVNotExistedHelpStr is a text described the persistent volume that doesn't exist
	MeshControlPlanePVNotExistedHelpStr = `EaseMesh control plane needs PersistentVolume to store data.
You need to create PersistentVolume in advance and specify its storageClassName as the value of --mesh-storage-class-name.

You can create PersistentVolume by the following definition:

apiVersion: v1
kind: PersistentVolume
metadata:
  labels:
    app: easemesh
  name: easemesh-pv
spec:
  storageClassName: {easemesh-storage}
  accessModes:
  - {ReadWriteOnce}
  capacity:
    storage: {3Gi}
  hostPath:
    path: {/opt/easemesh/}
    type: "DirectoryOrCreate"`

	// DefaultEasegressImage is default name of Easegress docker image
	DefaultEasegressImage = "megaease/easegress:easemesh"
	// DefaultEaseMeshOperatorImage is default name of the operator docker image
	DefaultEaseMeshOperatorImage = "megaease/easemesh-operator:latest"
	// DefaultShadowServiceControllerImage is default name of the shadow service docker image
	DefaultShadowServiceControllerImage = "megaease/easemesh-shadowservice-controller:latest"
	// DefaultImageRegistryURL is default registry url
	DefaultImageRegistryURL = "docker.io"
	// DefaultImagePullPolicy is default image pull policy.
	DefaultImagePullPolicy = v1.PullIfNotPresent
)

type (
	// OperationGlobal is global option for emctl
	OperationGlobal struct {
		MeshNamespace string
		EgServiceName string
	}

	// Install holds configurations for installation of the EaseMesh
	Install struct {
		*OperationGlobal

		ImageRegistryURL string
		ImagePullPolicy  string

		CleanWhenFailed bool

		// Easegress Control Plane params
		EasegressImage                string
		EasegressControlPlaneReplicas int
		EgClientPort                  int
		EgAdminPort                   int
		EgPeerPort                    int

		EgServicePeerPort  int
		EgServiceAdminPort int

		MeshControlPlaneStorageClassName      string
		MeshControlPlanePersistVolumeName     string
		MeshControlPlanePersistVolumeHostPath string
		MeshControlPlanePersistVolumeCapacity string
		MeshControlPlaneCheckHealthzMaxTime   int

		MeshIngressReplicas    int
		MeshIngressServicePort int32

		OnlyAddOn                    bool
		AddOns                       []string
		ShadowServiceControllerImage string

		// EaseMesh Controller  params
		EaseMeshRegistryType string
		HeartbeatInterval    int

		// EaseMesh Operator params
		EaseMeshOperatorImage    string
		EaseMeshOperatorReplicas int

		SpecFile string

		WaitControlPlaneTimeoutInSeconds int
	}

	// CoreDNS holds the options for installing EaseMesh-version CoreDNS.
	CoreDNS struct {
		*OperationGlobal
		Replicas        int
		Image           string
		CleanWhenFailed bool
	}

	// Reset holds the option for the EaseMesh resest sub command
	Reset struct {
		*OperationGlobal
		OnlyAddOn bool
		AddOns    []string
	}

	// AdminGlobal holds the option for all the EaseMesh admin command
	AdminGlobal struct {
		Server  string
		Timeout time.Duration
	}

	// AdminFileInput holds the option for all the EaseMesh admin command
	AdminFileInput struct {
		YamlFile  string
		Recursive bool
	}

	// Apply holds the option for the apply sub command
	Apply struct {
		*AdminGlobal
		*AdminFileInput
	}

	// Delete holds the option for the emctl delete sub command
	Delete struct {
		*AdminGlobal
		*AdminFileInput
	}

	// Get holds the option for the emctl get sub command
	Get struct {
		*AdminGlobal
		OutputFormat string
	}
)

// GetServerAddress return global server address configuration
func GetServerAddress() string {
	rc, err := rcfile.New()
	if err != nil {
		return ""
	}

	err = rc.Unmarshal()
	if err != nil {
		common.OutputErrorf("unmarshal rcfile failed: %v", err)
		return ""
	}
	return rc.Server
}

// AttachCmd attaches options for installation of coredns.
func (c *CoreDNS) AttachCmd(cmd *cobra.Command) {
	c.OperationGlobal = &OperationGlobal{}
	c.OperationGlobal.AttachCmd(cmd)

	cmd.Flags().BoolVar(&c.CleanWhenFailed, "clean-when-failed", true, "Clean resources when installation failed")
	cmd.Flags().IntVar(&c.Replicas, "replicas", 1, "CoreDNS replicas")
	cmd.Flags().StringVar(&c.Image, "image", "megaease/coredns:latest", "CoreDNS image name")
}

// AttachCmd attaches options for installation sub command
func (i *Install) AttachCmd(cmd *cobra.Command) {
	i.OperationGlobal = &OperationGlobal{}
	i.OperationGlobal.AttachCmd(cmd)
	cmd.Flags().IntVar(&i.EgClientPort, "mesh-control-plane-client-port", DefaultMeshClientPort, "Mesh control plane client port for remote accessing")
	cmd.Flags().IntVar(&i.EgAdminPort, "mesh-control-plane-admin-port", DefaultMeshAdminPort, "Port of mesh control plane admin for management")
	cmd.Flags().IntVar(&i.EgPeerPort, "mesh-control-plane-peer-port", DefaultMeshPeerPort, "Port of mesh control plane for consensus each other")
	cmd.Flags().IntVar(&i.MeshControlPlaneCheckHealthzMaxTime,
		"mesh-control-plane-check-healthz-max-time",
		DefaultMeshControlPlaneCheckHealthzMaxTime,
		"Max timeout in second for checking control panel component whether ready or not")

	cmd.Flags().IntVar(&i.EgServicePeerPort, "mesh-control-plane-service-peer-port", DefaultMeshPeerPort, "Port of Easegress cluster peer")
	cmd.Flags().IntVar(&i.EgServiceAdminPort, "mesh-control-plane-service-admin-port", DefaultMeshAdminPort, "Port of Easegress admin address")

	cmd.Flags().StringVar(&i.MeshControlPlaneStorageClassName, "mesh-storage-class-name", DefaultMeshControlPlaneStorageClassName, "Mesh storage class name")
	cmd.Flags().StringVar(&i.MeshControlPlanePersistVolumeCapacity, "mesh-control-plane-pv-capacity", DefaultMeshControlPlanePersistVolumeCapacity,
		MeshControlPlanePVNotExistedHelpStr)

	cmd.Flags().Int32Var(&i.MeshIngressServicePort, "mesh-ingress-service-port", DefaultMeshIngressServicePort, "Port of mesh ingress controller")

	cmd.Flags().StringVar(&i.EaseMeshRegistryType, "registry-type", DefaultMeshRegistryType, MeshRegistryTypeHelpStr)
	cmd.Flags().IntVar(&i.HeartbeatInterval, "heartbeat-interval", DefaultHeartbeatInterval, "Heartbeat interval for mesh service")

	cmd.Flags().StringVar(&i.ImageRegistryURL, "image-registry-url", DefaultImageRegistryURL, "Image registry URL")
	cmd.Flags().StringVar(&i.ImagePullPolicy, "image-pull-policy", string(DefaultImagePullPolicy), "Image pull policy (support Always, IfNotPresent, Never)")
	cmd.Flags().StringVar(&i.EasegressImage, "easegress-image", DefaultEasegressImage, "Easegress image name")
	cmd.Flags().StringVar(&i.EaseMeshOperatorImage, "easemesh-operator-image", DefaultEaseMeshOperatorImage, "Mesh operator image name")

	cmd.Flags().IntVar(&i.EasegressControlPlaneReplicas, "easemesh-control-plane-replicas", DefaultMeshControlPlaneReplicas, "Mesh control plane replicas")
	cmd.Flags().IntVar(&i.MeshIngressReplicas, "easemesh-ingress-replicas", DefaultMeshIngressReplicas, "Mesh ingress controller replicas")
	cmd.Flags().BoolVar(&i.OnlyAddOn, "only-add-on", false, "Only install add-ons")
	cmd.Flags().StringArrayVar(&i.AddOns, "add-ons", []string{}, "Names of add-ons to be installed (support shadowservice)")
	cmd.Flags().StringVar(&i.ShadowServiceControllerImage, "shadowservice-controller-image", DefaultShadowServiceControllerImage, "Shadow service controller image name")
	cmd.Flags().IntVar(&i.EaseMeshOperatorReplicas, "easemesh-operator-replicas", DefaultMeshOperatorReplicas, "Mesh operator controller replicas")
	cmd.Flags().StringVarP(&i.SpecFile, "file", "f", "", "A yaml file specifying the install params")
	cmd.Flags().BoolVar(&i.CleanWhenFailed, "clean-when-failed", true, "Clean resources when installation failed")
	cmd.Flags().IntVar(&i.WaitControlPlaneTimeoutInSeconds, "wait-control-plane-seconds", DefaultWaitControlPlaneSeconds, "Wait control plane ready timeout in seconds")
}

// AttachCmd attaches options for reset sub command
func (r *Reset) AttachCmd(cmd *cobra.Command) {
	r.OperationGlobal = &OperationGlobal{}
	r.OperationGlobal.AttachCmd(cmd)
	cmd.Flags().BoolVar(&r.OnlyAddOn, "only-add-on", false, "Only reset add-ons")
	cmd.Flags().StringArrayVar(&r.AddOns, "add-ons", []string{}, "Names of add-ons to be reset")
}

// AttachCmd attaches options globally
func (o *OperationGlobal) AttachCmd(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.MeshNamespace, "mesh-namespace", DefaultMeshNamespace, "EaseMesh namespace in kubernetes")
	cmd.Flags().StringVar(&o.EgServiceName, "mesh-control-plane-service-name", DefaultMeshControlPlaneHeadfulServiceName, "Mesh control plane service name")
}

// AttachCmd attaches options for base administrator command
func (a *AdminGlobal) AttachCmd(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&a.Server, "server", "s", "", "An address to access the EaseMesh control plane")
	cmd.Flags().DurationVarP(&a.Timeout, "timeout", "t", 30*time.Second, "A duration that limit max time out for requesting the EaseMesh control plane")
}

// AttachCmd attaches file options for base administrator command
func (a *AdminFileInput) AttachCmd(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&a.YamlFile, "file", "f", "", "A location contained the EaseMesh resource files (YAML format) to apply, could be a file, directory, or URL")
	cmd.Flags().BoolVarP(&a.Recursive, "recursive", "r", true, "Whether to recursively iterate all sub-directories and files of the location")
}

// AttachCmd attaches options for apply sub command
func (a *Apply) AttachCmd(cmd *cobra.Command) {
	a.AdminGlobal = &AdminGlobal{}
	a.AdminGlobal.AttachCmd(cmd)

	a.AdminFileInput = &AdminFileInput{}
	a.AdminFileInput.AttachCmd(cmd)
}

// AttachCmd attaches options for delete sub command
func (d *Delete) AttachCmd(cmd *cobra.Command) {
	d.AdminGlobal = &AdminGlobal{}
	d.AdminGlobal.AttachCmd(cmd)

	d.AdminFileInput = &AdminFileInput{}
	d.AdminFileInput.AttachCmd(cmd)
}

// AttachCmd attaches options for get sub command
func (g *Get) AttachCmd(cmd *cobra.Command) {
	g.AdminGlobal = &AdminGlobal{}
	g.AdminGlobal.AttachCmd(cmd)

	cmd.Flags().StringVarP(&g.OutputFormat, "output", "o", "table", "Output format (support table, yaml, json)")
}
