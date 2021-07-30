package operator

import (
	"crypto/tls"

	"github.com/megaease/easemeshctl/cmd/client/command/flags"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("generate cert key", func() {
	It("shoud succeed", func() {
		installFlags := &flags.Install{
			OperationGlobal: &flags.OperationGlobal{
				MeshNamespace: "easemesh",
			},
		}

		cert, key, err := generateCertKey(installFlags)
		Expect(err).NotTo(HaveOccurred())

		_, err = tls.X509KeyPair(cert, key)
		Expect(err).NotTo(HaveOccurred())
	})
})
