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
	"crypto/x509"
	"encoding/pem"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("generate cert key", func() {
	It("shoud succeed", func() {
		csrPem, keyPem, err := generateCsrAndKeyPem("easemesh")
		Expect(err).NotTo(HaveOccurred())

		csrBlock, _ := pem.Decode(csrPem)
		keyBlock, _ := pem.Decode(keyPem)

		_, err = x509.ParseCertificateRequest(csrBlock.Bytes)
		Expect(err).NotTo(HaveOccurred())

		_, err = x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
		Expect(err).NotTo(HaveOccurred())
	})
})
