package env

import (
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/Sirupsen/logrus"
	"github.com/YuukanOO/alfredo/domain"
	"gopkg.in/mgo.v2"
)

const (
	// Version is the current alfredo version.
	Version = "0.1.0"
	// AdaptersCollection is the name of the mongo collection
	AdaptersCollection = "adapters"
	// ControllersCollection is the name of the mongo collection
	ControllersCollection = "controllers"
	// RoomsCollection is the name of the mongo collection
	RoomsCollection = "rooms"
	// DevicesCollection is the name of the mongo collection
	DevicesCollection = "devices"
)

// ServerConfig  represents the current environment settings
// relative to the HTTP server.
type ServerConfig struct {
	Listen           string
	UseHTTPS         bool
	Remote           string
	AllowedOrigins   []string
	ShellCommand     []string
	LogLevel         string
	AllowRegistering bool
	AdaptersPath     string
}

// SecurityConfig contains settings related to security.
type SecurityConfig struct {
	Secret string
}

// DatabaseConfig represents the current env settings relatives
// to the mongodb settings.
type DatabaseConfig struct {
	Urls    []string
	session *mgo.Session
}

// Env represents the running environment.
type Env struct {
	Server   *ServerConfig
	Database *DatabaseConfig
	Security *SecurityConfig
}

var current Env

// Current retrieve the current environment.
func Current() Env {
	return current
}

// LoadFromFile load a configuration from a toml file.
func LoadFromFile(path string) error {
	logrus.WithField("path", path).Info("Loading configuration")
	_, err := toml.DecodeFile(path, &current)

	if err != nil {
		return err
	}

	// Parse the log level and sets it
	logLevel, err := logrus.ParseLevel(current.Server.LogLevel)

	if err != nil {
		return err
	}

	logrus.SetLevel(logLevel)

	// Connect to the database
	logrus.WithField("urls", current.Database.Urls).Info("Connecting to the database")
	session, err := mgo.Dial(strings.Join(current.Database.Urls, ","))

	current.Database.session = session

	// Load adapters and refresh them if needed
	adaptersCollection := session.DB("").C(AdaptersCollection)

	logrus.WithField("path", current.Server.AdaptersPath).Info("Loading adapters")
	adapters, err := domain.LoadAdapters(adaptersCollection.Find, current.Server.AdaptersPath)

	if err != nil {
		return err
	}

	// Adapters loaded, drop the collection and take the loaded ones
	adaptersCollection.DropCollection()

	adaptersToInsert := make([]interface{}, len(adapters))

	for i, adp := range adapters {
		adaptersToInsert[i] = adp
	}

	if err := adaptersCollection.Insert(adaptersToInsert...); err != nil {
		return err
	}

	logrus.Infof("Loaded %d adapter(s)", len(adapters))

	return nil
}

// Cleanup necessary stuff such as handles.
func Cleanup() {
	logrus.Info("Cleaning environment")

	current.Database.session.Close()
}

// GetSession retrieve a session to the database. Don't forget to close it when you're done.
func (db *DatabaseConfig) GetSession() (*mgo.Session, *mgo.Database) {
	cloned := db.session.Clone()

	return cloned, cloned.DB("")
}
