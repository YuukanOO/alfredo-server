package env

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/YuukanOO/alfredo/domain"
	"gopkg.in/mgo.v2"
)

// VERSION is the current alfredo version.
const VERSION = "0.1.0"

// ServerConfig  represents the current environment settings
// relative to the HTTP server.
type ServerConfig struct {
	Listen         string
	Remote         string
	AllowedOrigins []string
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
	Adapters map[string]domain.Adapter
}

var current Env

// Current retrieve the current environment.
func Current() Env {
	return current
}

// LoadFromFile load a configuration from a toml file.
func LoadFromFile(path string) error {
	_, err := toml.DecodeFile(path, &current)

	if err != nil {
		return err
	}

	session, err := mgo.Dial(strings.Join(current.Database.Urls, ","))

	current.Database.session = session

	return err
}

// LoadAdapters tries to load the smart adapters from a JSON files.
func LoadAdapters(path string) error {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	return json.Unmarshal(data, &current.Adapters)
}

// Cleanup necessary stuff such as handles.
func Cleanup() {
	current.Database.session.Close()
}

// GetSession retrieve a session to the database. Don't forget to close it when you're done.
func (db *DatabaseConfig) GetSession() *mgo.Session {
	return db.session.Clone()
}
