package webapp

import "github.com/BurntSushi/toml"
import "fmt"

// Settings holds all application settings as given by the .toml file.
type Settings struct {
	Server ServerSettings
}

// ServerSettings holds settings related to the web server.
type ServerSettings struct {
	Host string
	Port int
}

func (s ServerSettings) listenAddr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func loadSettings(path string) (*Settings, error) {
	var settings Settings

	if _, err := toml.DecodeFile(path, &settings); err != nil {
		return nil, err
	}

	return &settings, nil
}
