package sidecarinjector_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSidecarInjector(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SidecarInjector Suite")
}
