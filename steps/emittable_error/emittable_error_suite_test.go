package emittable_error_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestEmittableError(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EmittableError Suite")
}
