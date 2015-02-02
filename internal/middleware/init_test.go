package middleware_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEnvoyMiddlewareSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Envoy Middleware Suite")
}
