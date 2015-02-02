package domain_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEnvoyDomainSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Envoy Domain Suite")
}
