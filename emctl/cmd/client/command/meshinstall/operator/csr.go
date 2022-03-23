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

	"github.com/pkg/errors"
	certv1 "k8s.io/api/certificates/v1"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"
)

// deployCSR deploys a new CertificateSigningRequest with approval action.
func deployCSR(ctx *installbase.StageContext, csrPem, keyPem []byte) (CertPem []byte, err error) {
	csr := &certv1.CertificateSigningRequest{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CertificateSigningRequest",
			APIVersion: "certificates.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      installbase.OperatorCSRName,
			Namespace: ctx.Flags.MeshNamespace,
		},
		Spec: certv1.CertificateSigningRequestSpec{
			// Reference: https://kubernetes.io/docs/reference/access-authn-authz/certificate-signing-requests/#kubernetes-signers
			SignerName: "kubernetes.io/kubelet-serving",
			Groups:     []string{"system:authenticated"},
			// NOTE: []byte will be automatically encoded as a base64-encoded string.
			// Reference: https://golang.org/pkg/encoding/json/#Marshal
			Request: []byte(csrPem),
			Usages: []certv1.KeyUsage{
				certv1.UsageKeyEncipherment,
				certv1.UsageDigitalSignature,
				certv1.UsageServerAuth,
			},
		},
	}

	_, err = ctx.Client.CertificatesV1().CertificateSigningRequests().Create(context.TODO(),
		csr, metav1.CreateOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "create CertificateSigningRequest failed")
	}

	for i := 0; ; i++ {
		csr, err = ctx.Client.CertificatesV1().CertificateSigningRequests().Get(context.TODO(),
			installbase.OperatorCSRName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		if len(csr.Status.Certificate) != 0 {
			return csr.Status.Certificate, nil
		}

		csr, err := approveCSR(csr)
		if err != nil {
			return nil, err
		}

		csr, err = ctx.Client.CertificatesV1().CertificateSigningRequests().UpdateApproval(context.TODO(),
			csr.Name, csr, metav1.UpdateOptions{})
		if k8serrors.IsConflict(err) && i < 10 {
			if csr != nil && len(csr.Status.Certificate) != 0 {
				return csr.Status.Certificate, nil
			}
			continue
		}
		if err != nil {
			return nil, err
		}
	}
}

func approveCSR(csr *certv1.CertificateSigningRequest) (*certv1.CertificateSigningRequest, error) {
	var alreadyHasCondition bool
	for _, c := range csr.Status.Conditions {
		if c.Type == certv1.CertificateDenied {
			return nil, errors.Errorf("certificate signing request %q is already %s", csr.Name, c.Type)
		}
		if c.Type == certv1.CertificateApproved {
			alreadyHasCondition = true
		}
	}
	if alreadyHasCondition {
		return csr, nil
	}

	csr.Status.Conditions = append(csr.Status.Conditions, certv1.CertificateSigningRequestCondition{
		Type:           certv1.RequestConditionType(certv1.CertificateApproved),
		Status:         v1.ConditionTrue,
		Reason:         "EmctlApprove",
		Message:        "This CSR was approved by emctl certificate approve.",
		LastUpdateTime: metav1.Now(),
	})

	return csr, nil
}
