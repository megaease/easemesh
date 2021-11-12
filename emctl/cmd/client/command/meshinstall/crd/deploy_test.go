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

package crd

import (
	"testing"

	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"
	meshtesting "github.com/megaease/easemeshctl/cmd/client/testing"

	"github.com/spf13/cobra"
	extensionfake "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	"k8s.io/client-go/kubernetes/fake"
)

func prepareContext() *installbase.StageContext {

	client := fake.NewSimpleClientset()
	exptensionClient := extensionfake.NewSimpleClientset()

	install := &flags.Install{}
	cmd := &cobra.Command{}
	install.AttachCmd(cmd)
	return meshtesting.PrepareInstallContext(cmd, client, exptensionClient, install)
}
func TestDeploy(t *testing.T) {
	Deploy(prepareContext())
}

func TestDescribePhase(t *testing.T) {
	ctx := prepareContext()
	DescribePhase(ctx, installbase.BeginPhase)
	DescribePhase(ctx, installbase.EndPhase)
	DescribePhase(ctx, installbase.ErrorPhase)
	PreCheck(ctx)
}

var helloWorld = "aGVsbG8gd29ybGQK"
