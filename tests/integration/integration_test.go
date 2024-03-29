package integration_test

import (
	"context"
	"os"
	"testing"

	"sb.im/gosd/tests/help"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	if err := help.Setup(); err != nil {
		panic(err)
	}
	code := m.Run()
	//clean()
	os.Exit(code)
}

func TestSmoke(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := help.StartSingleServer(t, ctx)
	assert.NoError(t, c.ServerStatus(), "Http Error")
}
