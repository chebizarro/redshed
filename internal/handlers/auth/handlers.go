package auth

import (
	"context"
	"encoding/json"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"

	"net/http"
	"time"

	"github.com/chebizarro/redshed/internal/orm"

	"github.com/dgrijalva/jwt-go"

	"github.com/chebizarro/redshed/internal/logger"
	"github.com/chebizarro/redshed/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

// User is a retrieved and authenticated user.
type User struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
	Hd            string `json:"hd"`
}

type LoginCode struct {
	code string `fom:"code" json:"code" binding:"required"`
}

// Claims JWT claims
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// Begin login with the auth provider
func Begin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// You have to add value context with provider name to get provider name in GetProviderName method
		c.Request = addProviderToContext(c, c.Param(string(utils.ProjectContextKeys.ProviderCtxKey)))
		// try to get the user without re-authenticating
		if gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request); err != nil {
			gothic.BeginAuthHandler(c.Writer, c.Request)
		} else {
			logger.Debugf("user: %#v", gothUser)
		}
	}
}

// Callback callback to complete auth provider flow
func Callback(cfg *utils.ServerConfig, orm *orm.ORM) gin.HandlerFunc {
	return func(c *gin.Context) {
		// You have to add value context with provider name to get provider name in GetProviderName method
		c.Request = addProviderToContext(c, c.Param(string(utils.ProjectContextKeys.ProviderCtxKey)))
		user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		u, err := orm.FindUserByJWT(user.Email, user.Provider, user.UserID)
		// logger.Debugf("gothUser: %#v", user)
		if err != nil {
			if u, err = orm.UpsertUserProfile(&user); err != nil {
				logger.Errorf("[Auth.CallBack.UserLoggedIn.UpsertUserProfile.Error]: %v", err)
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		}

		logger.Debug("[Auth.CallBack.UserLoggedIn]: ", u.ID)

		jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod(cfg.JWT.Algorithm), Claims{
			Email: user.Email,
			StandardClaims: jwt.StandardClaims{
				Id:        user.UserID,
				Issuer:    user.Provider,
				IssuedAt:  time.Now().UTC().Unix(),
				NotBefore: time.Now().UTC().Unix(),
				ExpiresAt: user.ExpiresAt.UTC().Unix(),
			},
		})
		token, err := jwtToken.SignedString([]byte(cfg.JWT.Secret))
		if err != nil {
			logger.Error("[Auth.Callback.JWT] error: ", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		logger.Debug("token: ", token)
		json := gin.H{
			"type":          "Bearer",
			"token":         token,
			"refresh_token": user.RefreshToken,
		}
		c.JSON(http.StatusOK, json)
	}
}

// Logout logs out of the auth provider
func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = addProviderToContext(c, c.Param(string(utils.ProjectContextKeys.ProviderCtxKey)))
		gothic.Logout(c.Writer, c.Request)
		c.Writer.Header().Set("Location", "/")
		c.Writer.WriteHeader(http.StatusTemporaryRedirect)
	}
}

// Callback callback to complete auth provider flow
func Exchange(cfg *utils.ServerConfig, orm *orm.ORM) gin.HandlerFunc {
	return func(c *gin.Context) {
		// You have to add value context with provider name to get provider name in GetProviderName method
		c.Request = addProviderToContext(c, c.Param(string(utils.ProjectContextKeys.ProviderCtxKey)))

		conf := &oauth2.Config{
			ClientID:     "104197714087-qalk2efcpqd4k82q489rndjn590tp5gs.apps.googleusercontent.com",
			ClientSecret: "_4C0GDEJxQjiz9hYCfvxv7wC",
			RedirectURL:  "postmessage",
			Endpoint:     google.Endpoint,
		}

		x, _ := ioutil.ReadAll(c.Request.Body)
		code := string(x)

		tok, err := conf.Exchange(context.Background(), code)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		client := conf.Client(context.Background(), tok)
		email, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		defer email.Body.Close()

		data, err := ioutil.ReadAll(email.Body)
		if err != nil {
			logger.Errorf("[Gin-OAuth] Could not read Body: %s", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		var user User
		err = json.Unmarshal(data, &user)
		if err != nil {
			logger.Errorf("[Gin-OAuth] Unmarshal userinfo failed: %s", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		_, err = orm.FindUserByJWT(user.Email, "google", user.Email)

		if err != nil {
			logger.Errorf("[Auth.CallBack.UserLoggedIn.UpsertUserProfile.Error]: %v", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			/*
				if u, err = orm.UpsertUserProfile(&user); err != nil {
				}*/
		}
		// logger.Debug("[Auth.CallBack.UserLoggedIn.USER]: ", u)
		//logger.Debug("[Auth.CallBack.UserLoggedIn]: ", u.ID)
		jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod(cfg.JWT.Algorithm), Claims{
			Email: user.Email,
			StandardClaims: jwt.StandardClaims{
				Id:        user.Hd,
				Issuer:    "google",
				IssuedAt:  time.Now().UTC().Unix(),
				NotBefore: time.Now().UTC().Unix(),
			},
		})
		token, err := jwtToken.SignedString([]byte(cfg.JWT.Secret))
		if err != nil {
			logger.Error("[Auth.Callback.JWT] error: ", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		logger.Debug("token: ", token)
		json := gin.H{
			"type":          "Bearer",
			"token":         token,
			"refresh_token": "####",
		}
		c.JSON(http.StatusOK, json)
	}
}
