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
	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func serviceSpec(installFlags *flags.Install) installbase.InstallFunc {
	labels := meshOperatorLabels()

	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      installbase.DefaultMeshOperatorControllerManagerServiceName,
			Namespace: installFlags.MeshNamespace,
		},
	}
	service.Spec.Ports = []v1.ServicePort{
		{
			Name:       "https",
			Port:       int32(8443),
			TargetPort: intstr.IntOrString{StrVal: "https"},
		},
	}
	service.Spec.Selector = labels
	return func(cmd *cobra.Command, kubeClient *kubernetes.Clientset, installFlags *flags.Install) error {
		err := installbase.DeployService(service, kubeClient, installFlags.MeshNamespace)
		if err != nil {
			return errors.Wrapf(err, "Create operator service %s", installFlags.MeshNamespace)
		}
		return err
	}
}
