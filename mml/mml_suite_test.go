package mml_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMml(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mml Suite")
}
