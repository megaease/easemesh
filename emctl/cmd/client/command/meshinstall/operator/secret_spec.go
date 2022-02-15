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
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"

	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func secretSpec(ctx *installbase.StageContext) installbase.InstallFunc {
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      installbase.DefaultMeshOperatorSecretName,
			Namespace: ctx.Flags.MeshNamespace,
		},
		Data: map[string][]byte{},
	}

	return func(ctx *installbase.StageContext) error {
		_, err := ctx.Client.CoreV1().Secrets(ctx.Flags.MeshNamespace).Get(context.TODO(),
			secret.Name, metav1.GetOptions{})
		if err == nil {
			fmt.Printf("\nsecret %s existed, won't create it again\n\n", secret.Name)
			return nil
		} else if !errors.IsNotFound(err) {
			return fmt.Errorf("deploy secret %s/%s failed: %v",
				ctx.Flags.MeshNamespace, secret.Name, err)
		}

		csrPem, keyPem, err := generateCsrAndKeyPem(ctx.Flags.MeshNamespace)
		if err != nil {
			return fmt.Errorf("generate csr and key failed: %v", err)
		}

		certPem, err := deployCSR(ctx, csrPem, keyPem)
		if err != nil {
			return fmt.Errorf("deploy CertificateSigningRequest failed: %v", err)
		}

		// NOTE: []byte will be automatically encoded as a base64-encoded string.
		// Reference: https://golang.org/pkg/encoding/json/#Marshal
		secret.Data[installbase.DefaultMeshOperatorCertFileName] = certPem
		secret.Data[installbase.DefaultMeshOperatorKeyFileName] = keyPem

		err = installbase.DeploySecret(secret, ctx.Client, ctx.Flags.MeshNamespace)
		if err != nil {
			return fmt.Errorf("deploy secret failed: %v", err)
		}

		return nil
	}
}

func generateCsrAndKeyPem(namespace string) ([]byte, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	template := x509.CertificateRequest{
		Subject: pkix.Name{
			Organization: []string{"system:nodes"},
			CommonName:   "system:node:MegaEase",
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
		Bytes: csrBytes,
	})
	if err != nil {
		return nil, nil, err
	}

	keyBuffer := &bytes.Buffer{}
	err = pem.Encode(keyBuffer, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		return nil, nil, err
	}

	return csrBuffer.Bytes(), keyBuffer.Bytes(), nil
}
