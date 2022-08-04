package help

import (
	"context"
	"net/http/httptest"
	"testing"

	"sb.im/gosd/app/api"
	"sb.im/gosd/app/client"
	"sb.im/gosd/app/cmd"
	"sb.im/gosd/app/config"
)

func StartSingleServer(t *testing.T, ctx context.Context) *client.Client {
	cfg := config.Parse()
	cfg.StorageURL = "file://" + t.TempDir()

	server := httptest.NewServer(cmd.NewHandler(ctx, cfg))
	return client.NewClient(server.URL+api.ApiPrefix, cfg.ApiKey)
}
