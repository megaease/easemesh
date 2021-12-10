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

package shadowservice

import (
	"testing"

	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"
	meshtesting "github.com/megaease/easemeshctl/cmd/client/testing"

	"github.com/spf13/cobra"
	extensionfake "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	"k8s.io/client-go/kubernetes/fake"
)

func prepareContext() (*installbase.StageContext, *fake.Clientset, *extensionfake.Clientset) {
	client := fake.NewSimpleClientset()
	extensionClient := extensionfake.NewSimpleClientset()

	install := &flags.Install{}
	cmd := &cobra.Command{}
	install.AttachCmd(cmd)
	return meshtesting.PrepareInstallContext(cmd, client, extensionClient, install), client, extensionClient
}

func TestDeploy(t *testing.T) {
	ctx, _, _ := prepareContext()

	for _, f := range []func(*installbase.StageContext) installbase.InstallFunc{
		deploymentSpec, shadowServiceKindSpec,
	} {
		f(ctx).Deploy(ctx)
	}

	Deploy(ctx)

}

func TestDescribePhase(t *testing.T) {
	ctx, _, _ := prepareContext()
	DescribePhase(ctx, installbase.BeginPhase)
	DescribePhase(ctx, installbase.EndPhase)
	DescribePhase(ctx, installbase.ErrorPhase)
	PreCheck(ctx)
}
