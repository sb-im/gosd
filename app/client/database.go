package client

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (c *Client) DatabaseMigrate() {
	_, err := c.do(http.MethodPost, c.endpoint+"/database/migrate", nil)
	if err != nil {
		log.Error(err)
	}
	log.Warn("=== Database Migrate Done ===")
}

func (c *Client) DatabaseSeed() {
	_, err := c.do(http.MethodPost, c.endpoint+"/database/seed", nil)
	if err != nil {
		log.Error(err)
	}
	log.Warn("=== Database Seed Done ===")
}
