package performgen_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPerformGen(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PerformGen Suite")
}
