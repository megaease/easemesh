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
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"

	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	"github.com/spf13/cobra"
	certv1beta1 "k8s.io/api/certificates/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func secretSpec(installFlags *flags.Install) installbase.InstallFunc {
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      installbase.DefaultMeshOperatorSecretName,
			Namespace: installFlags.MeshNamespace,
		},
	}

	return func(cmd *cobra.Command, client *kubernetes.Clientset, installFlags *flags.Install) error {
		certBase64, keyBase64, err := generateCertKeyInBase64(installFlags)
		if err != nil {
			return err
		}

		csr := &certv1beta1.CertificateSigningRequest{
			ObjectMeta: metav1.ObjectMeta{
				Name:      installbase.DefaultMeshOperatorCSRName,
				Namespace: installFlags.MeshNamespace,
			},
			Spec: certv1beta1.CertificateSigningRequestSpec{
				Groups:  []string{"system:authenticated"},
				Request: certBase64,
				Usages: []certv1beta1.KeyUsage{
					certv1beta1.UsageDigitalSignature,
					certv1beta1.UsageKeyEncipherment,
					certv1beta1.UsageServerAuth,
				},
			},
		}

		_, err = client.CertificatesV1beta1().CertificateSigningRequests().Create(context.TODO(), csr, metav1.CreateOptions{})
		if err != nil {
			return err
		}

		_, err = client.CertificatesV1beta1().CertificateSigningRequests().UpdateApproval(context.TODO(), csr, metav1.UpdateOptions{})
		if err != nil {
			return err
		}

		secret.Data[installbase.DefaultMeshOperatorCertFileName] = certBase64
		secret.Data[installbase.DefaultMeshOperatorKeyFileName] = keyBase64

		err = installbase.DeploySecret(secret, client, installFlags.MeshNamespace)
		if err != nil {
			return fmt.Errorf("create secret failed: %v ", err)
		}

		return err
	}
}

func generateCertKeyInBase64(installFlags *flags.Install) ([]byte, []byte, error) {
	cert, key, err := generateCertKey(installFlags)
	if err != nil {
		return nil, nil, err
	}

	certBase64 := base64.StdEncoding.EncodeToString(cert)
	keyBase64 := base64.StdEncoding.EncodeToString(key)

	return []byte(certBase64), []byte(keyBase64), nil
}

func generateCertKey(installFlags *flags.Install) ([]byte, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	now := time.Now()
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"MegaEase"},
		},
		DNSNames: []string{
			installbase.DefaultMeshOperatorServiceName,
			fmt.Sprintf("%s.%s", installbase.DefaultMeshOperatorServiceName, installFlags.MeshNamespace),
			fmt.Sprintf("%s.%s.svc", installbase.DefaultMeshOperatorServiceName, installFlags.MeshNamespace),
		},
		NotBefore:             now.UTC(),
		NotAfter:              now.Add(10 * 365 * 24 * time.Hour).UTC(),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template,
		&privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, err
	}

	certBuffer := &bytes.Buffer{}
	err = pem.Encode(certBuffer, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
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

	return certBuffer.Bytes(), keyBuffer.Bytes(), nil
}
