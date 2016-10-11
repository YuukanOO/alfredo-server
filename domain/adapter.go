package domain

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"text/template"
)

// Adapter represents an available smart adapter.
type Adapter struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Category string                 `json:"category"`
	Config   map[string]interface{} `json:"config"`
	Commands map[string]string      `json:"commands"`
	Widgets  map[string]string      `json:"widgets"`

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

// PrepareWidgets will transform jsx to valid js and cache them to be provided to clients.
func (adp *Adapter) PrepareWidgets(cachePath string) error {

	for k, v := range adp.Widgets {
		widgetStr := ""
		path := path.Join(cachePath, fmt.Sprintf("%s.%s.js", adp.ID, k))

		// Check if we have to read file contents
		if v[:1] == "<" {
			widgetStr = v
		} else {
			data, err := ioutil.ReadFile(v)

			if err != nil {
				return err
			}

			widgetStr = string(data)
		}

		// TODO: more robust error handling

		nodeCmd := exec.Command("node", "alfredo_prepare_widgets.js", widgetStr)

		var stdout bytes.Buffer
		var stderr bytes.Buffer

		nodeCmd.Stdout = &stdout
		nodeCmd.Stderr = &stderr

		err := nodeCmd.Run()

		if err != nil {
			return NewErrCommandFailed(nodeCmd, err, stderr.String())
		}

		if err := ioutil.WriteFile(path, stdout.Bytes(), os.ModeExclusive); err != nil {
			return err
		}
	}

	return nil
}

// Execute an adapter command for the given device and parametrization.
func (adp *Adapter) Execute(shell []string, command string, device *Device, params map[string]interface{}) (*ExecutionResult, error) {
	// Check if the command has been parsed first
	tmpl := adp.commandsParsed[command]

	// If not found, ensure the commands has been parsed
	if tmpl == nil && adp.Commands[command] != "" {
		if err := adp.ParseCommands(); err != nil {
			return nil, err
		}
		tmpl = adp.commandsParsed[command]
	}

	// If still null, then the command does not exists in this adapter
	if tmpl == nil {
		return nil, ErrAdapterCommandNotFound
	}

	// Creates the execution context available in the command template
	ctx := newExecutionContext(device.Config, params)

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, ctx); err != nil {
		return nil, err
	}

	// Executes the command
	// TODO: Robust use of shell[] index out of range...
	cmd := exec.Command(shell[0], append(shell[1:], buf.String())...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	result := &ExecutionResult{
		Success: cmd.ProcessState.Success(),
		Err:     stderr.String(),
		Out:     stdout.String(),
	}

	if err != nil {
		return result, err
	}

	return result, nil
}
