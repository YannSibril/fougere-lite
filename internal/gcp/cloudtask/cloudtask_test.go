// ©Copyright 2022-2023 Metrio
package cloudtask

import (
	"context"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/api/option"
	"metrio.net/fougere-lite/internal/utils"
)

// Helper method to create client
func getMockedClient(url string) *Client {
	client, err := NewClient(context.Background(), option.WithoutAuthentication(), option.WithEndpoint(url))
	if err != nil {
		Fail(err.Error())
	}
	return client
}

var _ = Describe("Storage client", func() {
	var queueConfig TaskQueue

	BeforeEach(func() {
		queueConfig = TaskQueue{
			Name:                    "patate-23423k",
			Region:                  "northamerica-northeast1",
			ProjectId:               "projet-123",
			MinBackoff:              "1s",
			MaxBackoff:              "10s",
			MaxConcurrentDispatches: 1000,
			MaxDispatchesPerSecond:  500.0,
			ClientName:              "banane",
		}
	})
	Describe("create queue", func() {
		It("successfully creates the queue", func() {
			mockServerCalls := make(chan utils.MockServerCall, 1)
			mockServerCalls <- utils.MockServerCall{
				UrlMatchFunc: func(url string) bool {
					return strings.HasPrefix(url, "/v2/projet-123")
				},
				Method: "post",
			}
			mockServer := utils.NewMockServer(mockServerCalls)
			defer mockServer.Close()

			client := getMockedClient(mockServer.URL)

			err := client.create(queueConfig)
			Expect(err).ToNot(HaveOccurred())
		})
	})
	Describe("update queue", func() {
		It("successfully updates the queue", func() {
			mockServerCalls := make(chan utils.MockServerCall, 1)
			mockServerCalls <- utils.MockServerCall{
				UrlMatchFunc: func(url string) bool {
					return strings.HasPrefix(url, "/v2/patate-23423k?")
				},
				Method: "patch",
			}
			mockServer := utils.NewMockServer(mockServerCalls)
			defer mockServer.Close()

			client := getMockedClient(mockServer.URL)

			err := client.update(queueConfig)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
