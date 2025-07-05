package health_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHealthSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Health Suite")
} 