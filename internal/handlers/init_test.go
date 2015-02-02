package handlers_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEnvoyHandlersSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Envoy Handlers Suite")
}
