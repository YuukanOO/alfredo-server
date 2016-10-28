package env

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"os/exec"

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
	// TransformScript represents the script used to transform jsx
	TransformScript = "console.log(require('babel-core').transform(process.argv[1], { plugins: ['transform-react-jsx', 'transform-es2015-arrow-functions'],}).code);"
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

// transformWidget will transform the jsx into a valid React component ready to be used.
func transformWidget(widget string) (string, error) {
	// TODO: more robust error handling
	cmd := exec.Command("node", "-e", TransformScript, widget)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		return "", domain.NewErrCommandFailed(cmd, err, stderr.String())
	}

	return fmt.Sprintf("function(device, command, detail) { return %s; }", strings.TrimSpace(stdout.String())), nil
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
