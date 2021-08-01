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
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"

	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	certv1beta1 "k8s.io/api/certificates/v1beta1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func csrSpec(ctx *installbase.StageContext) installbase.InstallFunc {
	return func(ctx *installbase.StageContext) error {
		csrPem, keyPem, err := generateCsrAndKeyPem(ctx.Flags.MeshNamespace)
		if err != nil {
			return err
		}

		ctx.OperatorCsrPem, ctx.OperatorKeyPem = csrPem, keyPem

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
			return err
		}

		for i := 0; ; i++ {
			csr, err = ctx.Client.CertificatesV1beta1().CertificateSigningRequests().Get(context.TODO(),
				installbase.DefaultMeshOperatorCSRName, metav1.GetOptions{})
			if err != nil {
				return err
			}

			if len(csr.Status.Certificate) != 0 {
				ctx.OperatorCertPem = csr.Status.Certificate
				return nil
			}

			csr, err := approveCSR(csr)
			if err != nil {
				return err
			}

			csr, err = ctx.Client.CertificatesV1beta1().CertificateSigningRequests().UpdateApproval(context.TODO(),
				csr, metav1.UpdateOptions{})
			if errors.IsConflict(err) && i < 10 {
				if csr != nil && len(csr.Status.Certificate) != 0 {
					ctx.OperatorCertPem = csr.Status.Certificate
					return nil
				}
				continue
			}
			if err != nil {
				return err
			}
		}
	}
}

func generateCsrAndKeyPem(namespace string) ([]byte, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	template := x509.CertificateRequest{
		Subject: pkix.Name{
			Organization: []string{"MegaEase"},
		},

		DNSNames: []string{
			installbase.DefaultMeshOperatorServiceName,
			fmt.Sprintf("%s.%s", installbase.DefaultMeshOperatorServiceName, namespace),
			fmt.Sprintf("%s.%s.svc", installbase.DefaultMeshOperatorServiceName, namespace),
		},

		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, privateKey)
	if err != nil {
		return nil, nil, err
	}

	csrBuffer := &bytes.Buffer{}
	err = pem.Encode(csrBuffer, &pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBytes})
	if err != nil {
		return nil, nil, err
	}

	keyBuffer := &bytes.Buffer{}
	err = pem.Encode(keyBuffer, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	if err != nil {
		return nil, nil, err
	}

	return csrBuffer.Bytes(), keyBuffer.Bytes(), nil
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
