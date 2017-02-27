package webapp

import (
	"gopkg.in/gin-gonic/gin.v1"
)

var settings *Settings

// Serve starts the web server and listen to incoming connections.
func Serve(confPath string) error {
	settings, err := loadSettings(confPath)

	if err != nil {
		return err
	}

	r := gin.Default()

	return r.Run(settings.Server.listenAddr())
}
