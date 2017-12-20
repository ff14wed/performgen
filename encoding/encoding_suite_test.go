package encoding_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestEncoding(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Encoding Suite")
}
