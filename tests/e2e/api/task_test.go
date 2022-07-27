package api_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"time"

	"sb.im/gosd/app/api"
	"sb.im/gosd/app/client"
	"sb.im/gosd/app/cmd"
	"sb.im/gosd/app/config"
	"sb.im/gosd/app/model"

	"sb.im/gosd/tests/help"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Task", func() {
	var s *httptest.Server
	var c *client.Client

	nodeID := "1"

	task := model.Task{
		Name:   "E2E Test",
		NodeID: nodeID,
	}

	ctx, cancel := context.WithCancel(context.Background())
	handler := cmd.NewHandler(ctx)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())
		s = httptest.NewServer(handler)

		// TODO: need to ApiKey
		c = client.NewClient(s.URL+api.ApiPrefix, "")

		go help.StartNcp(ctx, config.Parse().MqttURL, nodeID)

		// Wait mqttd server startup && sub topic on broker
		time.Sleep(100 * time.Millisecond)
	})

	AfterEach(func() {
		s.Close()
		cancel()
	})

	Context("Test Context", func() {
		It("create a new task", func() {
			fmt.Println(os.Getenv("LUA_FILE"))
			if node, err := c.NodeShow(nodeID); err != nil {
				c.NodeCreate(&model.Node{
					ID:   nodeID,
					Name: "Test Node",
				})
			} else {
				fmt.Println(node)
			}
			err := c.TaskCreate(&task)
			Expect(err).NotTo(HaveOccurred())
		})

		It("run this new task", func() {
			err := c.JobCreate(&task)
			Expect(err).NotTo(HaveOccurred())
		})

	})
})
