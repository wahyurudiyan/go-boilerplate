package configz_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConfigz(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Configz Suite")
}
