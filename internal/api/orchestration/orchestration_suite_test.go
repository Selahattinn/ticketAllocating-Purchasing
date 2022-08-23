package orchestration_test

import (
	"errors"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	errService = errors.New("service err")
)

func TestAPIOrchestration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API Orchestration Suite")
}
