package env

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/YuukanOO/alfredo/domain"
	"gopkg.in/mgo.v2"
)

const (
	// Version is the current alfredo version.
	Version = "0.1.0"
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
	Listen         string
	UseHTTPS       bool
	Remote         string
	AllowedOrigins []string
	ShellCommand   []string
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
	Adapters []*domain.Adapter

	adaptersMap map[string]*domain.Adapter
}

var current Env

// Current retrieve the current environment.
func Current() Env {
	return current
}

// GetAdapter retrieves an adapter given its ID.
func (env Env) GetAdapter(id string) *domain.Adapter {
	return env.adaptersMap[id]
}

// LoadFromFile load a configuration from a toml file.
func LoadFromFile(path string) error {
	_, err := toml.DecodeFile(path, &current)

	if err != nil {
		return err
	}

	session, err := mgo.Dial(strings.Join(current.Database.Urls, ","))

	current.Database.session = session

	// Init some variables here
	current.adaptersMap = map[string]*domain.Adapter{}

	return err
}

// transformWidget will append necessary arrow function to easily integrates the widget in
// a client environment.
func transformWidget(widget string) string {
	return fmt.Sprintf("function(device, command) { return %s; }", widget)
}

// LoadAdapters tries to load the smart adapters from a JSON file.
func LoadAdapters(path string) error {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, &current.Adapters); err != nil {
		return err
	}

	for _, v := range current.Adapters {
		if err = v.ParseCommands(); err != nil {
			return err
		}

		if err := v.ParseWidgets(transformWidget); err != nil {
			return err
		}

		// Add the adapter to the inner map (for easy and fast retrieval)
		current.adaptersMap[v.ID] = v
	}

	return nil
}

// Cleanup necessary stuff such as handles.
func Cleanup() {
	current.Database.session.Close()
}

// GetSession retrieve a session to the database. Don't forget to close it when you're done.
func (db *DatabaseConfig) GetSession() (*mgo.Session, *mgo.Database) {
	cloned := db.session.Clone()

	return cloned, cloned.DB("")
}
