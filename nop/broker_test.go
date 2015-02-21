package nop_test

import (
	"github.com/pivotal-cf-experimental/envoy"
	"github.com/pivotal-cf-experimental/envoy/nop"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Broker", func() {
	var broker interface{}

	It("acts as a no-op implementation of the broker API", func() {
		broker = nop.Broker{}

		_, ok := broker.(envoy.Broker)
		Expect(ok).To(BeTrue())
	})
})
