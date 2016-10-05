package domain

import "text/template"

// Adapter represents an available smart adapter.
type Adapter struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Category string                 `json:"category"`
	Config   map[string]interface{} `json:"config"`
	Commands map[string]string      `json:"commands"`

	commandsParsed map[string]*template.Template
}

// ParseCommands parses all commands to a valid go text template ready to be use.
func (adp *Adapter) ParseCommands() error {
	commands := map[string]*template.Template{}

	for k, v := range adp.Commands {
		tmpl, err := template.New(k).Parse(v)

		if err != nil {
			return err
		}

		commands[k] = tmpl
	}

	adp.commandsParsed = commands

	return nil
}
