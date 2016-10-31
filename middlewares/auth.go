package middlewares

import (
	"net/http"
	"strings"

	"github.com/YuukanOO/alfredo/domain"
	"github.com/YuukanOO/alfredo/env"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

const (
	authHeaderValue   = "Authorization"
	bearerHeaderValue = "Bearer "
	controllerKey     = "controller"
)

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

				// db := GetDB(c)
				// var controller domain.Controller

				// Check if the controller exists in the db. The token is signed so it must exists but maybe in the future
				// you will be able to blacklist tokens and so I check it right now.
				// if err := db.C(env.ControllersCollection).FindId(bson.ObjectIdHex(idStr)).One(&controller); err != nil || controller.Token != tokenStr {
				// 	if err != nil {
				// 		AbortWithError(c, http.StatusForbidden, err)
				// 	} else {
				// 		c.AbortWithStatus(http.StatusForbidden)
				// 	}
				// } else {
				// 	c.Set(controllerKey, &controller)
				// 	c.Next()
				// }

				c.Set(controllerKey, &domain.Controller{
					ID: bson.ObjectIdHex(idStr),
				})
				c.Next()
			}
		}
	}
}
