package v3

import (
	"fmt"
	"time"

	"sb.im/gosd/app/model"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	identityGinKey  = "jwt_current"
	identityTeamKey = "tid"
	identityUserKey = "uid"
	identitySessKey = "sid"
)

var (
	defaultCurrent = &Current{
		TeamID: 1,
		UserID: 1,
		SessID: 1,
	}
)

type bindLogin struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func (h Handler) InitAuth(r *gin.RouterGroup) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte(h.cfg.Secret),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityGinKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*Current); ok {
				return jwt.MapClaims{
					identityTeamKey: v.TeamID,
					identityUserKey: v.UserID,
					identitySessKey: v.SessID,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &Current{
				TeamID: uint(claims[identityTeamKey].(float64)),
				UserID: uint(claims[identityUserKey].(float64)),
				SessID: uint(claims[identitySessKey].(float64)),
			}
		},
		Authenticator: h.login,
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// TODO: other auth method
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	r.POST("/login", authMiddleware.LoginHandler)

	//r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
	//	claims := jwt.ExtractClaims(c)
	//	log.Printf("NoRoute claims: %#v\n", claims)
	//	c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	//})

	//auth := r.Group("/auth")
	// Refresh time can be longer than token timeout
	r.GET("/refresh_token", authMiddleware.RefreshHandler)

	// Set Middleware to Handler
	//r.Use(authMiddleware.MiddlewareFunc())
	authMid := authMiddleware.MiddlewareFunc()

	isAuth := func(c *gin.Context) bool {
		_, ok := c.Get(identityGinKey)
		return ok
	}

	// 1. Auth: SingleUserMode
	r.Use(func(c *gin.Context) {
		if h.singleUserMode() {
			c.Set(identityGinKey, defaultCurrent)
		}
		c.Next()
	})

	// 2. Auth: header X-Api-Key
	r.Use(func(c *gin.Context) {
		if v := c.Request.Header.Values("X-Api-Key"); len(v) > 0 && v[0] == h.cfg.ApiKey {
			c.Set(identityGinKey, defaultCurrent)
		}
		c.Next()
	})

	// 3. Auth: header Basic Authentication
	r.Use(func(c *gin.Context) {
		if h.cfg.BasicAuth {
			if v := c.Request.Header.Values("Authorization"); len(v) > 0 {
				if username, password, err := basicAuthDecode(v[0]); err == nil {
					var user model.User
					if h.orm.Where("username = ?", username).First(&user).Error == nil {
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
		}
		c.Next()
	})

	// 4. Auth: header JWT token
	r.Use(func(c *gin.Context) {
		if isAuth(c) {
			c.Next()
		} else {
			authMid(c)
		}
	})
}

// @Summary User Login
// @Schemes Auth
// @Description user login
// @Tags auth
// @Accept multipart/form-data
// @Produce json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Success 200
// @Router /login [POST]
func (h Handler) login(c *gin.Context) (interface{}, error) {
	var login bindLogin
	if err := c.ShouldBind(&login); err != nil {
		return nil, jwt.ErrMissingLoginValues
	}

	var user model.User
	h.orm.Where("username = ?", login.Username).First(&user)
	fmt.Println(user)
	if err := user.VerifyPassword(login.Password); err == nil {

		// Create Session
		session := &model.Session{
			TeamID: user.TeamID,
			UserID: user.ID,
			IP:     c.ClientIP(),
		}
		h.orm.Save(session)

		current := Current{
			TeamID: user.TeamID,
			UserID: user.ID,
			SessID: session.ID,
		}
		return &current, nil
	}

	return nil, jwt.ErrFailedAuthentication
}

// @Summary Token Refresh
// @Schemes Auth
// @Description Refresh a token expired time
// @Tags auth
// @Security JWTSecret
// @Accept multipart/form-data
// @Produce json
// @Success 200
// @Router /refresh_token [GET]
func (h *Handler) Refresh(c *gin.Context) {}
