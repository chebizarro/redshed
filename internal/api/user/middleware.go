package user

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	identityKey = "id"
	Auth *jwt.GinJWTMiddleware
	err  error
)

// Login is login struct
type login struct {
	Username string `form:"username" validate:"required"`
	Password string `form:"password" validate:"required"`
}

// Validator .
func (login *login) validator(c *gin.Context) (*User, string, bool) {

	db := c.MustGet("db").(*gorm.DB)
	msg := "login successful"

	// Get model if exist
	var user User
	if err := db.Where("id = ?", login.Username).First(&user).Error; err != nil {
		msg = "The account does not exist"
		return nil, msg, false
	}

	if user.Password != login.Password {
		msg = "Incorrect password"
		return nil, msg, false
	}

	return &user, msg, true
}


func init() {
	// the jwt middleware
	fmt.Println("init jwt")
	Auth, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,

		PayloadFunc: func(data interface{}) jwt.MapClaims {
			fmt.Println("PayloadFunc")
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.Username,
				}
			}
			return jwt.MapClaims{
				identityKey: data,
			}
		},

		IdentityHandler: func(c *gin.Context) interface{} {
			fmt.Println("transfer：IdentityHandler")
			claims := jwt.ExtractClaims(c)
			// return &User{
			// 	Username: claims["id"].(string),
			// }
			fmt.Println(claims[identityKey])
			return claims[identityKey]
		},

		Authenticator: func(c *gin.Context) (interface{}, error) {
			login := &login{}
			if err := c.ShouldBind(login); err != nil {
				return "", err
			}
			user, msg, result := login.validator(c)

			if result {
				session := sessions.Default(c)
				session.Set("role", user.Role)
				session.Save()
				return &user, nil
			}

			return nil, errors.New(msg)
		},

		Authorizator: func(data interface{}, c *gin.Context) bool {
			session := sessions.Default(c)
			session.Set("userInfo", data)
			session.Save()
			return true
		},

		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"token":   token,
				"expire":  expire.Format(time.RFC3339),
				"message": "login success!",
			})
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

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,

		// Optionally return the token as a cookie
		SendCookie: true,
	})
}
