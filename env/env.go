package env

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"os/exec"

	"github.com/BurntSushi/toml"
	"github.com/Sirupsen/logrus"
	"github.com/YuukanOO/alfredo/domain"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	// Version is the current alfredo version.
	Version = "0.1.0"
	// SystemCollection is the name of the mongo collection related to system settings
	SystemCollection = "system"
	// EnvKey is the name used as a selector to persist the env to the database
	EnvKey = "env"
	// AdaptersCollection is the name of the mongo collection
	AdaptersCollection = "adapters"
	// ControllersCollection is the name of the mongo collection
	ControllersCollection = "controllers"
	// RoomsCollection is the name of the mongo collection
	RoomsCollection = "rooms"
	// DevicesCollection is the name of the mongo collection
	DevicesCollection = "devices"
	// TransformScript represents the script used to transform jsx
	TransformScript = "console.log(require('babel-core').transform(process.argv[1], { plugins: ['transform-react-jsx', 'transform-es2015-arrow-functions'],}).code);"
)

// SystemSettings represents a part of the env stored in the database.
type SystemSettings struct {
	Key              string
	Version          string
	AdaptersCheckSum string
}

// AdapterEntries represents the result of loading adapters via the environment.
type AdapterEntries struct {
	Checksum string
	Adapters []*domain.Adapter
}

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

	session, err := mgo.Dial(strings.Join(current.Database.Urls, ","))

	current.Database.session = session

	// Init some variables here
	current.adaptersMap = map[string]*domain.Adapter{}

	return err
}

// transformWidget will transform the jsx into a valid React component ready to be used.
func transformWidget(widget string) (string, error) {
	logrus.WithField("widget", widget).Debug("Transforming")

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

	return fmt.Sprintf("function(device, command, showView) { return %s; }", strings.TrimSpace(stdout.String())), nil
}

// LoadAdapters tries to load the smart adapters from a JSON file.
func LoadAdapters(oldChecksum string, path string) (*AdapterEntries, error) {
	logrus.WithFields(logrus.Fields{
		"path":     path,
		"checksum": oldChecksum,
	}).Info("Loading adapters")

	data, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	loadedAdapters := []*domain.Adapter{}
	curChecksum := domain.ComputeCheckSum(data)
	logChecksum := logrus.WithField("checksum", curChecksum)

	if oldChecksum == curChecksum {
		// TODO: We must compare each adapters widgets defined in another jsx file
		logChecksum.Info("Adapters checksum did not change, continue")
	} else {
		logChecksum.Info("Adapters checksum have changed")

		if err = json.Unmarshal(data, &current.Adapters); err != nil {
			return nil, err
		}

		for _, v := range current.Adapters {
			logAdapter := logrus.WithField("adapter", v.Name)

			logAdapter.Info("Parsing commands")

			if err = v.ParseCommands(); err != nil {
				return nil, err
			}

			logAdapter.Info("Parsing widgets")

			if err := v.ParseWidgets(transformWidget); err != nil {
				return nil, err
			}

			// Add the adapter to the inner map (for easy and fast retrieval)
			current.adaptersMap[v.ID] = v
			loadedAdapters = append(loadedAdapters, v)
		}
	}

	return &AdapterEntries{
		Checksum: curChecksum,
		Adapters: loadedAdapters,
	}, nil
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

// ByKey selects elements by their key property
func ByKey(key string) bson.M {
	return bson.M{
		"key": key,
	}
}
