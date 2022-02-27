package v3

import (
	"errors"
	"net/http"
	"time"

	"sb.im/gosd/app/model"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
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

func InitAuthMiddleware(r *gin.RouterGroup, h *Handler) error {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "gosd_zone",
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
		Authenticator: h.handlerAuthJWTlogin,
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
		LoginResponse:   responseAuthJWTLogin,
		RefreshResponse: responseAuthJWTRefresh,
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
		return err
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	if err = authMiddleware.MiddlewareInit(); err != nil {
		return err
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
	if h.singleUserMode() {
		r.Use(middlewareAuthSingleUser())
	}

	// 2. Auth: header X-Api-Key
	if h.cfg.ApiKey != "" {
		r.Use(middlewareAuthApiKey("X-Api-Key", h.cfg.ApiKey))
	}

	// 3. Auth: header Basic Authentication
	if h.cfg.BasicAuth {
		r.Use(middlewareAuthBasic(h.orm))
	}

	// 4. Auth: header JWT token
	r.Use(func(c *gin.Context) {
		if isAuth(c) {
			c.Next()
		} else {
			authMid(c)
		}
	})

	r.POST("/switch_team/:id", func(c *gin.Context) {
		current, err := h.handlerAuthJWTSwitchTeam(c)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if token, expire, err := authMiddleware.TokenGenerator(current); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, &responseJWT{
				Token:  token,
				Expire: expire,
			})
		}
	})
	return nil
}

func (h *Handler) handlerAuthJWTlogin(c *gin.Context) (interface{}, error) {
	var login bindLogin
	if err := c.ShouldBind(&login); err != nil {
		return nil, jwt.ErrMissingLoginValues
	}

	var user model.User
	if err := h.orm.Where("username = ?", login.Username).First(&user).Error; err != nil {
		return nil, err
	}
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

type responseJWT struct {
	Expire time.Time `json:"expire" example:"2022-02-22T02:20:22.002222+08:00"`
	Token  string    `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDU5NDcyNDAsIm9yaWdfaWF0IjoxNjQ1OTQzNjQwLCJzaWQiOjcsInRpZCI6MSwidWlkIjoxfQ.fnkuA08Be8Q3HGjOFdmND5Kc8aqWABXaoUravKX0bqg"`
}

// @Summary User Login
// @Schemes Auth
// @Description user login
// @Tags auth
// @Accept multipart/form-data
// @Produce json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Success 200 {object} responseJWT
// @Failure 401
// @Router /login [POST]
func responseAuthJWTLogin(c *gin.Context, code int, message string, expire time.Time) {
	c.JSON(code, &responseJWT{
		Expire: expire,
		Token:  message,
	})
}

// @Summary Token Refresh
// @Schemes Auth
// @Description Refresh a token expired time
// @Tags auth
// @Security JWTSecret
// @Accept multipart/form-data
// @Produce json
// @Success 200 {object} responseJWT
// @Failure 401
// @Router /refresh_token [GET]
func responseAuthJWTRefresh(c *gin.Context, code int, message string, expire time.Time) {
	c.JSON(code, &responseJWT{
		Expire: expire,
		Token:  message,
	})
}

// @Summary Switch Current Team
// @Schemes Auth
// @Description Switch Current Team
// @Tags auth
// @Security JWTSecret
// @Accept multipart/form-data
// @Produce json
// @Param id path uint true "Team ID"
// @Success 200 {object} responseJWT
// @Failure 404
// @Router /switch_team/{id} [POST]
func (h *Handler) handlerAuthJWTSwitchTeam(c *gin.Context) (current *Current, err error) {
	current = h.getCurrent(c)
	teamId := mustStringToUint(c.Param("id"))

	var count int64
	if err = h.orm.Find(&model.UserTeam{}, "team_id = ? AND user_id = ?", teamId, current.UserID).Count(&count).Error; err != nil {
		return
	}

	if count == 0 {
		err = errors.New("Not Found This Team")
		return
	}

	current.TeamID = teamId
	return
}
