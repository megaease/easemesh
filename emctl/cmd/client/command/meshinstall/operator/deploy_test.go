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
	"testing"

	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	meshtesting "github.com/megaease/easemeshctl/cmd/client/testing"
	"github.com/spf13/cobra"
	extensionfake "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	"k8s.io/client-go/kubernetes/fake"
)

func TestDeploy(t *testing.T) {
	client := fake.NewSimpleClientset()
	exptensionClient := extensionfake.NewSimpleClientset()

	install := &flags.Install{}
	cmd := &cobra.Command{}
	install.AttachCmd(cmd)
	ctx := meshtesting.PrepareInstallContext(cmd, client, exptensionClient, install)
	Deploy(ctx)

	secretSpec(ctx).Deploy(ctx)
	configMapSpec(ctx).Deploy(ctx)
	roleSpec(ctx).Deploy(ctx)
	clusterRoleSpec(ctx).Deploy(ctx)
	roleBindingSpec(ctx).Deploy(ctx)
	clusterRoleBindingSpec(ctx).Deploy(ctx)

	operatorDeploymentSpec(ctx).Deploy(ctx)

	serviceSpec(ctx).Deploy(ctx)
	mutatingWebhookSpec(ctx).Deploy(ctx)
}
