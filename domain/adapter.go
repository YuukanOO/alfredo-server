package domain

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"text/template"
)

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

// Execute an adapter command for the given device and parametrization.
func (adp *Adapter) Execute(shell []string, command string, device *Device, params map[string]interface{}) error {
	// Check if the command has been parsed first
	tmpl := adp.commandsParsed[command]

	// If not found, ensure the commands has been parsed
	if tmpl == nil && adp.Commands[command] != "" {
		if err := adp.ParseCommands(); err != nil {
			return err
		}
		tmpl = adp.commandsParsed[command]
	}

	// If still null, then the command does not exists in this adapter
	if tmpl == nil {
		return errors.New("CommandNotFound")
	}

	// Creates the execution context available in the command template
	ctx := newExecutionContext(device.Config, params)

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, ctx); err != nil {
		return err
	}

	// Executes the command
	// TODO: Robust use of shell[] index out of range...
	cmd := exec.Command(shell[0], append(shell[1:], buf.String())...)

	out, err := cmd.Output()

	if err != nil {
		return err
	}

	fmt.Println(string(out))

	return nil
}
