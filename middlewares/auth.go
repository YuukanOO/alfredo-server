package middlewares

import (
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/YuukanOO/alfredo/domain"
	"github.com/YuukanOO/alfredo/env"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const authHeaderValue = "Authorization"
const bearerHeaderValue = "Bearer "

const controllerKey = "controller"

// GetController retrieves the controller attached to this context
func GetController(c *gin.Context) *domain.Controller {
	return c.MustGet(controllerKey).(*domain.Controller)
}

// Auth restrict access to valid tokens.
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get(authHeaderValue)

		if authHeader == "" || !strings.Contains(authHeader, bearerHeaderValue) {
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			tokenStr := authHeader[len(bearerHeaderValue):]

			token, err := jwt.Parse(tokenStr, func(tok *jwt.Token) (interface{}, error) {
				return []byte(env.Current().Security.Secret), nil
			})

			if !token.Valid || err != nil {
				c.AbortWithStatus(http.StatusForbidden)
			} else {
				claims, _ := token.Claims.(jwt.MapClaims)
				idStr := claims["sub"].(string)

				c.Set(controllerKey, &domain.Controller{
					ID: bson.ObjectIdHex(idStr),
				})
				c.Next()
			}
		}
	}
}
