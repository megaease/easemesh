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
	"testing"

	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"
	meshtesting "github.com/megaease/easemeshctl/cmd/client/testing"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	extensionfake "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
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
	ctx, client, _ := prepareContext()
	ctx.Flags.WaitControlPlaneTimeOutInSeconds = 1

	client.PrependReactor("create", "secrets", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, nil
	})
	client.PrependReactor("create", "certificatesigningrequests", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, nil
	})

	for _, f := range []func(*installbase.StageContext) installbase.InstallFunc{
		configMapSpec, serviceSpec, serviceSpec, statefulsetSpec, namespaceSpec,
	} {
		f(ctx).Deploy(ctx)
	}

	client.PrependReactor("get", "secrets", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &v1.Secret{}, nil
	})

	client.PrependReactor("get", "services", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &v1.Service{
			Spec: v1.ServiceSpec{
				Type: v1.ServiceTypeNodePort,
				Ports: []v1.ServicePort{
					{
						NodePort: 1,
					},
				},
			},
		}, nil
	})
	client.PrependReactor("list", "nodes", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &v1.NodeList{
			Items: []v1.Node{{Status: v1.NodeStatus{Addresses: []v1.NodeAddress{{Type: v1.NodeInternalIP, Address: "127.0.0.2"}}}}},
		}, nil
	})

	ctx.Flags.WaitControlPlaneTimeOutInSeconds = 1
	ctx.Flags.MeshControlPlaneCheckHealthzMaxTime = 1
	Deploy(ctx)

	provisionEaseMeshControlPlane(ctx)

	checkEasegressControlPlaneStatus(ctx)

	clearEaseMeshControlPlaneProvision(ctx.Cmd, ctx.Client, ctx.Flags)

}

func TestCheckPV(T *testing.T) {
	checkPVAccessModes(v1.ReadWriteOnce, &v1.PersistentVolume{})
	checkPVAccessModes(v1.ReadWriteOnce, &v1.PersistentVolume{Spec: v1.PersistentVolumeSpec{AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce}}})
}

func TestUnmarshal(t *testing.T) {
	unmarshalMember([]byte{})
	unmarshalMember([]byte("test"))
}

func TestDescribePhase(t *testing.T) {
	ctx, _, _ := prepareContext()
	DescribePhase(ctx, installbase.BeginPhase)
	DescribePhase(ctx, installbase.EndPhase)
	DescribePhase(ctx, installbase.ErrorPhase)
	PreCheck(ctx)
}

var helloWorld = "aGVsbG8gd29ybGQK"
