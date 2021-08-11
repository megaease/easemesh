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
