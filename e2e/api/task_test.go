package api_test

import (
	"context"
	"net/http/httptest"

	"sb.im/gosd/app/cmd"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Task", func() {
	var server *httptest.Server
	ctx, cancel := context.WithCancel(context.Background())
	handler := cmd.NewHandler(ctx)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())
		server = httptest.NewServer(handler)
	})

	AfterEach(func() {
		server.Close()
		cancel()
	})

	Context("Test Context", func() {
		It("should be a novel", func() {
			Expect(1).To(Equal(1))
		})
	})
})
