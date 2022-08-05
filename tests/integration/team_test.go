package integration_test

import (
	"context"
	"testing"

	"sb.im/gosd/app/model"
	"sb.im/gosd/tests/help"

	"github.com/stretchr/testify/assert"
)

func TestTeamCreate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := help.StartSingleServer(t, ctx)

	team := &model.Team{
		Name: t.Name(),
	}

	assert.NoError(t, c.TeamCreate(team), "Http Error")
	assert.NotEqual(t, team.ID, uint(1), "team id error")
}
