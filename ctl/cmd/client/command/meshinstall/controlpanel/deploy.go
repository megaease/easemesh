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

package controlpanel

import (
	"fmt"
	"strings"
	"time"

	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"
	"github.com/megaease/easemeshctl/cmd/common"
	"github.com/megaease/easemeshctl/cmd/common/client"
	"gopkg.in/yaml.v2"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/kubernetes"
)

// Deploy will deploy resource of control panel
func Deploy(context *installbase.StageContext) error {

	installFuncs := []installbase.InstallFunc{
		namespaceSpec(context.Flags),
		configMapSpec(context.Flags),
		serviceSpec(context.Flags),
		statefulsetSpec(context.Flags),
	}

	err := installbase.BatchDeployResources(context.Cmd, context.Client, context.Flags, installFuncs)
	if err != nil {
		return errors.Wrap(err, "deploy mesh control panel resource")
	}

	err = checkEasegressControlPlaneStatus(context.Cmd, context.Client, context.Flags)
	if err != nil {
		return errors.Wrap(err, "check mesh control panel status")
	}

	err = provisionEaseMeshControlPanel(context.Cmd, context.Client, context.Flags)
	if err != nil {
		return errors.Wrap(err, "provision mesh control panel")
	}
	return nil
}

// PreCheck will check prerequisite for installing control plane
func PreCheck(context *installbase.StageContext) error {
	var err error

	// 1. check available PersistentVolume
	pvList, err := installbase.ListPersistentVolume(context.Client)
	if err != nil {
		return err
	}

	availablePVCount := 0
	quantity := resource.MustParse(context.Flags.MeshControlPlanePersistVolumeCapacity)
	boundedPVCSuffixes := []string{}
	for i := 0; i < context.Flags.EasegressControlPlaneReplicas; i++ {
		boundedPVCSuffixes = append(boundedPVCSuffixes, fmt.Sprintf("%s-%d", installbase.DefaultMeshControlPlaneName, i))

	}
	for _, pv := range pvList.Items {
		if pv.Status.Phase == v1.VolumeAvailable &&
			pv.Spec.StorageClassName == context.Flags.MeshControlPlaneStorageClassName &&
			pv.Spec.Capacity.Storage().Cmp(quantity) >= 0 &&
			checkPVAccessModes(v1.ReadWriteOnce, &pv) {
			availablePVCount++
		} else if pv.Status.Phase == v1.VolumeBound {
			// If PV already bound to PVC of EaseMesh controlpanel
			// we regarded it as availablePVCount
			for _, pvNameSuffix := range boundedPVCSuffixes {
				if pv.Spec.ClaimRef.Kind == "PersistentVolumeClaim" &&
					pv.Spec.ClaimRef.Namespace == context.Flags.MeshNamespace &&
					strings.HasSuffix(pv.Spec.ClaimRef.Name, pvNameSuffix) {
					availablePVCount++
					break
				}
			}
		}
	}

	if availablePVCount < context.Flags.EasegressControlPlaneReplicas {
		return errors.Errorf(flags.MeshControlPlanePVNotExistedHelpStr)
	}

	return nil

}

// Clear will clear all installed resource about control panel
func Clear(context *installbase.StageContext) error {
	statefulsetResource := [][]string{
		{"statefulsets", installbase.DefaultMeshControlPlaneName},
	}
	coreV1Resources := [][]string{
		{"services", context.Flags.EgServiceName},
		{"services", installbase.DefaultMeshControlPlanePlubicServiceName},
		{"services", installbase.DefaultMeshControlPlaneHeadlessServiceName},
		{"configmaps", installbase.DefaultMeshControlPlaneConfig},
	}

	clearEaseMeshControlPanelProvision(context.Cmd, context.Client, context.Flags)

	installbase.DeleteResources(context.Client, statefulsetResource, context.Flags.MeshNamespace, installbase.DeleteStatefulsetResource)
	installbase.DeleteResources(context.Client, coreV1Resources, context.Flags.MeshNamespace, installbase.DeleteCoreV1Resource)
	return nil
}

// Describe leverage human-readable text to describe different phase
// in the process of the control plane installation
func Describe(context *installbase.StageContext, phase installbase.InstallPhase) string {
	switch phase {
	case installbase.BeginPhase:
		return fmt.Sprintf("Begin to install mesh control plane service in the namespace %s", context.Flags.MeshNamespace)
	case installbase.EndPhase:
		return fmt.Sprintf("\nControl panel statefulset %s\n%s", installbase.DefaultMeshControlPlaneName,
			installbase.FormatPodStatus(context.Client, context.Flags.MeshNamespace,
				installbase.AdaptListPodFunc(meshControlPanelLabel())))
	}
	return ""
}

func checkPVAccessModes(accessModel v1.PersistentVolumeAccessMode, volume *v1.PersistentVolume) bool {
	for _, mode := range volume.Spec.AccessModes {
		if mode == accessModel {
			return true
		}
	}
	return false
}

func checkEasegressControlPlaneStatus(cmd *cobra.Command, kubeClient *kubernetes.Clientset, installFlags *flags.Install) error {

	// Wait a fix time for the Easegress cluster to start
	time.Sleep(time.Second * 10)

	entrypoints, err := installbase.GetMeshControlPanelEntryPoints(kubeClient, installFlags.MeshNamespace,
		installbase.DefaultMeshControlPlanePlubicServiceName,
		installbase.DefaultMeshAdminPortName)
	if err != nil {
		return errors.Wrap(err, "get mesh control plane entrypoint failed")
	}

	timeOutPerTry := installFlags.MeshControlPlaneCheckHealthzMaxTime / len(entrypoints)

	for i := 0; i < len(entrypoints); i++ {
		_, err := client.NewHTTPJSON(
			client.WrapRetryOptions(3, time.Second*time.Duration(timeOutPerTry)/3, func(body []byte, err error) bool {
				if err != nil && strings.Contains(err.Error(), "connection refused") {
					return true
				}

				members, err := unmarshalMember(body)
				if err != nil {
					common.OutputErrorf("parse member body error: %s", err)
					return true
				}

				return len(members) < (installFlags.EaseMeshOperatorReplicas/2 + 1)
			})...).
			Get(entrypoints[i]+installbase.MemberList, nil, time.Second*time.Duration(timeOutPerTry), nil).
			HandleResponse(func(body []byte, statusCode int) (interface{}, error) {
				if statusCode != 200 {
					return nil, errors.Errorf("check control plane member list error, return status code is :%d", statusCode)
				}
				members, err := unmarshalMember(body)

				if err != nil {
					return nil, err
				}

				if len(members) < (installFlags.EasegressControlPlaneReplicas/2 + 1) {
					return nil, errors.Errorf("easemesh control plane is not ready, expect %d of replicas, but %d", installFlags.EasegressControlPlaneReplicas, len(members))
				}
				return nil, nil
			})
		if err != nil {
			common.OutputErrorf("check mesh control plane status failed, ignored check next node, current error is: %s", err)
		} else {
			return nil
		}
	}
	return errors.Errorf("mesh control plane is not ready")
}

func unmarshalMember(body []byte) ([]map[string]interface{}, error) {
	var options []map[string]interface{}
	err := yaml.Unmarshal(body, &options)
	if err != nil {
		return nil, err
	}
	return options, nil
}
