// ©Copyright 2022-2023 Metrio
package cloudtask_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCloudtasks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cloudtasks Suite")
}
