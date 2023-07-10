// ©Copyright 2022-2023 Metrio
package cloudtask

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var validBucketConfig = []byte(`
cloudTasks:
  metrio-test:
    region: us-central1
    projectId: some-project`)

var invalidConfig = []byte(`
cloudTasks:
  some-queue:
    region:
      - should_not_be_an_array`)

var _ = Describe("config", func() {
	BeforeEach(func() {
		viper.Reset()
		viper.SetConfigType("yaml")
	})
	Describe("GetTaskConfig", func() {
		It("should successfully parse a storage queue config", func() {
			err := viper.ReadConfig(bytes.NewBuffer(validBucketConfig))
			Expect(err).ToNot(HaveOccurred())
			storageConfig, err := GetTaskConfig(viper.GetViper(), "metrio-client")
			Expect(err).To(BeNil())
			Expect(len(storageConfig.CloudTasks)).To(Equal(1))
			queue := storageConfig.CloudTasks["metrio-test"]
			Expect(queue.Region).To(Equal("us-central1"))
			Expect(queue.ProjectId).To(Equal("some-project"))
		})
		It("returns an error if cannot parse the config", func() {
			err := viper.ReadConfig(bytes.NewBuffer(invalidConfig))
			Expect(err).ToNot(HaveOccurred())
			_, err = GetTaskConfig(viper.GetViper(), "metrio-client")
			Expect(err).NotTo(BeNil())
		})
	})
	Context("validates task queues", func() {
		It("should not detect error", func() {
			config := &Config{
				CloudTasks: map[string]TaskQueue{
					"foooo": {
						Region:                  "us-central1",
						ProjectId:               "mock-project",
						Name:                    "foooo",
						MinBackoff:              "1s",
						MaxBackoff:              "10s",
						MaxConcurrentDispatches: 1000,
						MaxDispatchesPerSecond:  500.0,
					},
				},
			}
			err := ValidateConfig(config)
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should detect empty name task queues", func() {
			config := &Config{
				CloudTasks: map[string]TaskQueue{
					"foooo": {
						Region:    "us-central1",
						ProjectId: "mock-project",
					},
				},
			}
			err := ValidateConfig(config)
			Expect(err).Should(MatchError(ContainSubstring("validate failed on the required rule")))
		})
		It("should detect an empty region", func() {
			config := &Config{
				CloudTasks: map[string]TaskQueue{
					"foooo": {
						ProjectId: "mock-project",
						Name:      "foooo",
					},
				},
			}
			err := ValidateConfig(config)
			Expect(err).Should(MatchError(ContainSubstring("validate failed on the required rule")))
		})
		It("should detect a missing project id", func() {
			config := &Config{
				CloudTasks: map[string]TaskQueue{
					"foooo": {
						Name: "foooo",
					},
				},
			}
			err := ValidateConfig(config)
			Expect(err).Should(MatchError(ContainSubstring("validate failed on the required rule")))
		})
	})
})
