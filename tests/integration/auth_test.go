package integration_test

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"sb.im/gosd/app/api"
	"sb.im/gosd/app/client"
	"sb.im/gosd/app/cmd"
	"sb.im/gosd/tests/help"

	"github.com/stretchr/testify/assert"
)

func TestNoAuth(t *testing.T) {
	t.Setenv("SINGLE_USER", "true")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := help.StartSingleServer(t, ctx)
	assert.NoError(t, c.ServerStatus(), "Http Error")
	time.Sleep(1 * time.Second)
}

func TestAuthApiKey(t *testing.T) {
	t.Setenv("SINGLE_USER", "false")
	t.Setenv("API_KEY", "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := help.Config()
	cfg.StorageURL = "file://" + t.TempDir()

	server := httptest.NewServer(cmd.NewHandler(ctx, cfg))
	cb := client.NewClient(server.URL+api.ApiPrefix, "error_api_key")
	assert.Error(t, cb.ServerStatus(), "Http Error")

	cg := client.NewClient(server.URL+api.ApiPrefix, cfg.ApiKey)
	assert.NoError(t, cg.ServerStatus(), "Http Error")
}
