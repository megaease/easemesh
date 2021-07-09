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
	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func namespaceSpec(installFlags *flags.Install) installbase.InstallFunc {
	ns := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{
		Name:   installFlags.MeshNameSpace,
		Labels: map[string]string{},
	}}
	return func(cmd *cobra.Command, client *kubernetes.Clientset, installFlags *flags.Install) error {
		err := installbase.CreateNameSpace(ns, client)
		if err != nil && !errors.IsAlreadyExists(err) {
			return err
		}
		return nil
	}
}
