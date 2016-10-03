package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

// CORSOptions holds settings for the CORS middleware
type CORSOptions struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposeHeaders    []string
	AllowCredentials bool
}

// CORS handle OPTIONS request for basic CORS handling
func CORS(s *CORSOptions) gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Writer.Header().Set("Access-Control-Allow-Origin", strings.Join(s.AllowedOrigins, ","))
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(s.AllowedMethods, ","))
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(s.AllowedHeaders, ","))
		c.Writer.Header().Set("Access-Control-Expose-Headers", strings.Join(s.ExposeHeaders, ","))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(s.AllowCredentials))

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
		} else {
			c.Next()
		}
	}
}
