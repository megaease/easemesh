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

package operator

import (
	"context"
	"fmt"

	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	admissionregv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func mutatingWebhookSpec(ctx *installbase.StageContext) installbase.InstallFunc {
	mutatingPath := installbase.DefaultMeshOperatorMutatingWebhookPath
	mutatingPort := int32(installbase.DefaultMeshOperatorMutatingWebhookPort)
	mutatingScope := admissionregv1.NamespacedScope
	mutatingSideEffects := admissionregv1.SideEffectClassNoneOnDryRun

	mutatingWebhookConfig := func(caBundle []byte) *admissionregv1.MutatingWebhookConfiguration {
		return &admissionregv1.MutatingWebhookConfiguration{
			ObjectMeta: metav1.ObjectMeta{
				Name:      installbase.DefaultMeshOperatorMutatingWebhookName,
				Namespace: ctx.Flags.MeshNamespace,
			},
			Webhooks: []admissionregv1.MutatingWebhook{
				{
					Name: "mesh-injector.megaease.com",
					NamespaceSelector: &metav1.LabelSelector{
						MatchExpressions: []metav1.LabelSelectorRequirement{
							{
								Key:      "kubernetes.io/metadata.name",
								Operator: metav1.LabelSelectorOpNotIn,
								Values: []string{
									ctx.Flags.MeshNamespace,
									"kube-system",
									"kube-public",
								},
							},
							{
								Key:      "mesh.megaease.com/mesh-service",
								Operator: metav1.LabelSelectorOpExists,
							},
						},
					},
					ClientConfig: admissionregv1.WebhookClientConfig{
						Service: &admissionregv1.ServiceReference{
							Name:      installbase.DefaultMeshOperatorServiceName,
							Namespace: ctx.Flags.MeshNamespace,
							Path:      &mutatingPath,
							Port:      &mutatingPort,
						},
						CABundle: caBundle,
					},
					Rules: []admissionregv1.RuleWithOperations{
						{
							Operations: []admissionregv1.OperationType{
								admissionregv1.Create,
								admissionregv1.Update,
							},
							Rule: admissionregv1.Rule{
								APIGroups:   []string{"", "apps"},
								APIVersions: []string{"v1"},
								Resources: []string{
									"pods",
									"replicasets",
									"deployments",
									"statefulsets",
									"daemonsets",
								},
								Scope: &mutatingScope,
							},
						},
					},
					SideEffects:             &mutatingSideEffects,
					AdmissionReviewVersions: []string{"v1"},
				},
			},
		}
	}

	return func(ctx *installbase.StageContext) error {
		secret, err := ctx.Client.CoreV1().Secrets(ctx.Flags.MeshNamespace).Get(context.TODO(), installbase.DefaultMeshOperatorSecretName, metav1.GetOptions{})
		if err != nil {
			return err
		}

		certBase64, exists := secret.Data[installbase.DefaultMeshOperatorCertFileName]
		if !exists {
			return fmt.Errorf("key %v in secret %s not found",
				installbase.DefaultMeshOperatorCertFileName,
				installbase.DefaultMeshOperatorSecretName)
		}

		config := mutatingWebhookConfig(certBase64)

		err = installbase.DeployMutatingWebhookConfig(config, ctx.Client, ctx.Flags.MeshNamespace)
		if err != nil {
			return fmt.Errorf("create configMap failed: %v ", err)
		}
		return err
	}
}
