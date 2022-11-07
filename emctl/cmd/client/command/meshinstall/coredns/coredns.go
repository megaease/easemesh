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

package coredns

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"
	"github.com/megaease/easemeshctl/cmd/client/command/meshinstall/installation"
	"github.com/megaease/easemeshctl/cmd/common"
	"sigs.k8s.io/yaml"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	coreDNSSpecFile = "coredns-old-spec.yaml"

	warnMessage = fmt.Sprintf(`WARN: The process of installation for coredns can't be reverted in $ emctl reset
you could use generated file to revert coredns spec by: $ kubectl apply -f %s
`, coreDNSSpecFile)
)

// CoreDNSCmd returns coredns command.
func CoreDNSCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "coredns",
		Short:   "Deploy EaseMesh dedicated CoreDNS",
		Example: "emctl install coredns --clean-when-failed",
	}
	flags := &flags.CoreDNS{}
	flags.AttachCmd(cmd)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		fmt.Println(warnMessage)

		var err error
		kubeClient, clientConfig, err := installbase.NewKubernetesClient()
		if err != nil {
			common.ExitWithErrorf("%s failed: %v", cmd.Short, err)
		}

		apiExtensionClient, err := installbase.NewKubernetesAPIExtensionsClient()
		if err != nil {
			common.ExitWithErrorf("%s failed: %v", cmd.Short, err)
		}

		ctx := &installbase.StageContext{
			Cmd:                 cmd,
			CoreDNSFlags:        flags,
			ClientConfig:        clientConfig,
			Client:              kubeClient,
			APIExtensionsClient: apiExtensionClient,
		}

		err = storeOldCoreDNS(ctx)
		if err != nil {
			fmt.Printf("WARN: store old coredns spec failed: %v\n\n", err)
		}

		stages := []installation.InstallStage{
			installation.Wrap(PreCheck, Deploy, Clear, DescribePhase),
		}

		install := installation.New(stages...)

		err = install.DoInstallStage(ctx)
		if err != nil {
			if flags.CleanWhenFailed {
				install.ClearResource(ctx)
			}
			common.ExitWithErrorf("install coredns failed: %s", err)
		}
	}

	return cmd
}

func storeOldCoreDNS(ctx *installbase.StageContext) error {
	deploy, err := ctx.Client.AppsV1().Deployments(coreDNSNamespace).Get(context.Background(), coreDNSDeployment, metav1.GetOptions{})
	if err != nil {
		return err
	}

	deploy.Kind = "Deployment"
	deploy.APIVersion = "apps/v1"

	buff, err := yaml.Marshal(deploy)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(coreDNSSpecFile, buff, 0o644)
	if err != nil {
		return err
	}

	return nil
}
