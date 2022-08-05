package help

import (
	"context"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"sb.im/gosd/app/api"
	"sb.im/gosd/app/client"
	"sb.im/gosd/app/cmd"
	"sb.im/gosd/app/config"
)

const (
	databaseSuffix = "_test"
)

func Config() *config.Config {
	cfgP := config.Parse()
	cfg := *cfgP

	cfg.Debug = false

	u, err := url.Parse(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}

	// Set Test Database
	u.Path = strings.TrimLeft(u.Path, "/") + databaseSuffix
	cfg.DatabaseURL = u.String()
	return &cfg
}

func StartSingleServer(t *testing.T, ctx context.Context) *client.Client {
	cfg := Config()
	cfg.StorageURL = "file://" + t.TempDir()
	server := httptest.NewServer(cmd.NewHandler(ctx, cfg))
	return client.NewClient(server.URL+api.ApiPrefix, cfg.ApiKey)
}

func Setup() error {
	rawCfg := config.Parse()
	rawOrm, err := cmd.DatabaseOrm(rawCfg)
	if err != nil {
		return err
	}

	u, _ := url.Parse(rawCfg.DatabaseURL)
	dbName := strings.TrimLeft(u.Path, "/") + databaseSuffix

	rawOrm.Exec("DROP DATABASE " + dbName)
	rawOrm.Exec("CREATE DATABASE " + dbName)

	orm, err := cmd.DatabaseOrm(Config())
	if err != nil {
		return err
	}

	cmd.DatabaseMigrate(orm)
	return cmd.DatabaseSeed(orm)
}
