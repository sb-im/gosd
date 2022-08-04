package integration_test

import (
	"context"
	"testing"

	"sb.im/gosd/tests/help"

	"github.com/stretchr/testify/assert"
)

func TestSmoke(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := help.StartSingleServer(t, ctx)
	assert.NoError(t, c.ServerStatus(), "Http Error")
}
