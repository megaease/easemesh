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
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"
	meshtesting "github.com/megaease/easemeshctl/cmd/client/testing"

	"github.com/spf13/cobra"
	appsV1 "k8s.io/api/apps/v1"
	certv1beta1 "k8s.io/api/certificates/v1beta1"
	v1 "k8s.io/api/core/v1"
	extensionfake "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestDeploy(t *testing.T) {
	client := fake.NewSimpleClientset()
	exptensionClient := extensionfake.NewSimpleClientset()

	install := &flags.Install{}
	cmd := &cobra.Command{}
	install.AttachCmd(cmd)
	ctx := meshtesting.PrepareInstallContext(cmd, client, exptensionClient, install)
	Deploy(ctx)

	client.PrependReactor("create", "secrets", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, nil
	})
	client.PrependReactor("create", "certificatesigningrequests", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, nil
	})

	for _, f := range []func(*installbase.StageContext) installbase.InstallFunc{
		secretSpec, configMapSpec, roleSpec, clusterRoleSpec, roleBindingSpec, clusterRoleBindingSpec,
		operatorDeploymentSpec, serviceSpec, mutatingWebhookSpec,
	} {
		f(ctx).Deploy(ctx)
	}

	client.PrependReactor("get", "secrets", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &v1.Secret{}, nil
	})

	mutatingWebhookSpec(ctx).Deploy(ctx)
	secretSpec(ctx).Deploy(ctx)

	client.PrependReactor("get", "secrets", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &v1.Secret{
			Data: map[string][]byte{installbase.DefaultMeshOperatorCertFileName: []byte(helloWorld)},
		}, nil
	})
	mutatingWebhookSpec(ctx).Deploy(ctx)

	client.PrependReactor("get", "secrets", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, k8serr.NewNotFound(schema.GroupResource{
			Resource: "Namespace",
			Group:    "v1",
		}, "na")
	})
	client.PrependReactor("get", "certificatesigningrequests", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &certv1beta1.CertificateSigningRequest{
			TypeMeta: metav1.TypeMeta{
				Kind:       "CertificateSigningRequest",
				APIVersion: "certificates.k8s.io/v1beta1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      installbase.DefaultMeshOperatorCSRName,
				Namespace: ctx.Flags.MeshNamespace,
			},
			Spec: certv1beta1.CertificateSigningRequestSpec{
				Groups: []string{"system:authenticated"},
				// NOTE: []byte will be automatically encoded as a base64-encoded string.
				// Reference: https://golang.org/pkg/encoding/json/#Marshal
				Request: []byte(helloWorld),
				Usages: []certv1beta1.KeyUsage{
					certv1beta1.UsageDigitalSignature,
					certv1beta1.UsageKeyEncipherment,
					certv1beta1.UsageServerAuth,
				},
			},
		}, nil
	})
	secretSpec(ctx).Deploy(ctx)

	DescribePhase(ctx, installbase.BeginPhase)
	DescribePhase(ctx, installbase.EndPhase)

	client.PrependReactor("get", "*", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		var replicas int32 = 1
		switch action.GetResource().Resource {
		case "deployments":
			return true, &appsV1.Deployment{
				Spec: appsV1.DeploymentSpec{
					Replicas: &replicas,
				},
				Status: appsV1.DeploymentStatus{
					ReadyReplicas: replicas,
				},
			}, nil
		case "statefulsets":
			return true, &appsV1.StatefulSet{
				Spec: appsV1.StatefulSetSpec{
					Replicas: &replicas,
				},
				Status: appsV1.StatefulSetStatus{
					ReadyReplicas: replicas,
				},
			}, nil
		}
		return true, nil, nil
	})

	checkOperatorStatus(ctx.Client, ctx.Flags)

	PreCheck(ctx)

}

var helloWorld = "aGVsbG8gd29ybGQK"
