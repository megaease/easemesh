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

	certv1beta1 "k8s.io/api/certificates/v1beta1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// deployCSR deploys a new CertificateSigningRequest with approval action.
func deployCSR(ctx *installbase.StageContext, csrPem, keyPem []byte) (CertPem []byte, err error) {
	csr := &certv1beta1.CertificateSigningRequest{
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
			Request: []byte(csrPem),
			Usages: []certv1beta1.KeyUsage{
				certv1beta1.UsageDigitalSignature,
				certv1beta1.UsageKeyEncipherment,
				certv1beta1.UsageServerAuth,
			},
		},
	}

	_, err = ctx.Client.CertificatesV1beta1().CertificateSigningRequests().Create(context.TODO(),
		csr, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	for i := 0; ; i++ {
		csr, err = ctx.Client.CertificatesV1beta1().CertificateSigningRequests().Get(context.TODO(),
			installbase.DefaultMeshOperatorCSRName, metav1.GetOptions{})
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

		csr, err = ctx.Client.CertificatesV1beta1().CertificateSigningRequests().UpdateApproval(context.TODO(),
			csr, metav1.UpdateOptions{})
		if errors.IsConflict(err) && i < 10 {
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

func approveCSR(csr *certv1beta1.CertificateSigningRequest) (*certv1beta1.CertificateSigningRequest, error) {
	var alreadyHasCondition bool
	for _, c := range csr.Status.Conditions {
		if c.Type == certv1beta1.CertificateDenied {
			return nil, fmt.Errorf("certificate signing request %q is already %s", csr.Name, c.Type)
		}
		if c.Type == certv1beta1.CertificateApproved {
			alreadyHasCondition = true
		}
	}
	if alreadyHasCondition {
		return csr, nil
	}

	csr.Status.Conditions = append(csr.Status.Conditions, certv1beta1.CertificateSigningRequestCondition{
		Type:           certv1beta1.RequestConditionType(certv1beta1.CertificateApproved),
		Status:         v1.ConditionTrue,
		Reason:         "EmctlApprove",
		Message:        "This CSR was approved by emctl certificate approve.",
		LastUpdateTime: metav1.Now(),
	})

	return csr, nil
}
