package v3

import (
	"context"

	"sb.im/gosd/app/model"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func middlewareAuthSingleUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(identityGinKey, defaultCurrent)
		c.Next()
	}
}

func middlewareAuthApiKey(name, apikey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if v := c.Request.Header.Values(name); len(v) > 0 && v[0] == apikey {
			c.Set(identityGinKey, defaultCurrent)
		}
		c.Next()
	}
}

func middlewareAuthBasic(orm *gorm.DB) gin.HandlerFunc {
	const header = "Authorization"
	return func(c *gin.Context) {
		if v := c.Request.Header.Values(header); len(v) > 0 {
			if username, password, err := basicAuthDecode(v[0]); err == nil {
				var user model.User
				if orm.Where("username = ?", username).First(&user).Error == nil {
					if user.VerifyPassword(password) == nil {
						c.Set(identityGinKey, &Current{
							TeamID: user.TeamID,
							UserID: user.ID,
							SessID: 1,
						})
					}
				}
			}
		}
		c.Next()
	}
}

func middlewareAuthUrlToken(rdb *redis.Client) gin.HandlerFunc {
	const key = "token"
	return func(c *gin.Context) {
		token := c.Request.URL.Query().Get(key)
		if tid, err := rdb.Get(context.Background(), "token/"+token).Uint64(); err == nil {
			c.Set(identityGinKey, &Current{
				TeamID: uint(tid),
				UserID: 0,
				SessID: 1,
			})
		}
		c.Next()
	}
}
