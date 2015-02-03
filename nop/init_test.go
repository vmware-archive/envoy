package nop_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestNopSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Nop Suite")
}
