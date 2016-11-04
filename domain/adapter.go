package domain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

// TransformScript represents the script used to transform jsx
const TransformScript = "console.log(require('babel-core').transform(process.argv[1], { plugins: ['transform-react-jsx', 'transform-es2015-arrow-functions'],}).code);"

// Adapter represents an available smart adapter.
type Adapter struct {
	ID              string            `json:"id" bson:"_id"`
	Name            string            `json:"name" bson:"name"`
	Description     string            `json:"description" bson:"description"`
	Dependencies    []string          `json:"dependencies" bson:"dependencies"`
	Category        string            `json:"category" bson:"category"`
	Config          map[string]string `json:"config" bson:"config"`
	Commands        map[string]string `json:"commands" bson:"commands"`
	Widgets         map[string]string `json:"widgets" bson:"widgets"`
	WidgetsCheckSum map[string]string `json:"-" bson:"widgets_checksum"`

	commandsParsed map[string]*template.Template
}

func getWidgetBytes(relativeDir string, entry string) ([]byte, error) {
	if IsComponent(entry) {
		return []byte(entry), nil
	}

	data, err := ioutil.ReadFile(filepath.Join(relativeDir, entry))

	if err != nil {
		return nil, err
	}

	return data, nil
}

// LoadAdapters will load adapters given the path and retransform widgets as needed.
func LoadAdapters(findAdapters QueryFunc, path string) ([]Adapter, error) {
	// Read the adapters file first
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	dir := filepath.Dir(path)

	// And unmarshal the json to retrieve a list of adapters
	var fileAdapters []Adapter

	if err := json.Unmarshal(data, &fileAdapters); err != nil {
		return nil, err
	}

	loadedAdapters := make([]Adapter, len(fileAdapters))

	for i, adapter := range fileAdapters {
		var dbAdapter Adapter

		findAdapters(ByIDStr(adapter.ID)).One(&dbAdapter)

		adapter.WidgetsCheckSum = map[string]string{}

		// Compare checksums and retransforms widgets as needed
		for wk, w := range adapter.Widgets {
			// Compute file adapter checksum
			data, err := getWidgetBytes(dir, w)

			if err != nil {
				return nil, err
			}

			adapter.WidgetsCheckSum[wk] = ComputeCheckSum(data)

			// And then compare, if they are different, transform the widget
			if adapter.WidgetsCheckSum[wk] != dbAdapter.WidgetsCheckSum[wk] {
				widgetStr := string(data)
				widgetTransformed, err := transformWidget(widgetStr)

				if err != nil {
					return nil, newErrTransformWidgetFailed(&adapter, widgetStr, err, widgetTransformed)
				}

				adapter.Widgets[wk] = widgetTransformed
			} else {
				adapter.Widgets[wk] = dbAdapter.Widgets[wk]
			}
		}

		// Check for command dependencies
		if err := adapter.checkDependencies(); err != nil {
			return nil, err
		}

		// Parse adapter commands
		if err := adapter.parseCommands(); err != nil {
			return nil, err
		}

		loadedAdapters[i] = adapter
	}

	return loadedAdapters, nil
}

func (adp *Adapter) validateConfig(config map[string]interface{}) error {
	for ck := range adp.Config {
		// TODO: type checking
		if config[ck] == nil {
			return ErrDeviceConfigInvalid
		}
	}

	return nil
}

func (adp *Adapter) checkDependencies() error {
	for _, dep := range adp.Dependencies {
		if _, err := exec.LookPath(dep); err != nil {
			return newErrDependencyNotResolved(adp, dep, err)
		}
	}

	return nil
}

func (adp *Adapter) parseCommands() error {
	commands := map[string]*template.Template{}

	for k, v := range adp.Commands {
		tmpl, err := template.New(k).Parse(v)

		if err != nil {
			return newErrParseCommandFailed(adp, k, err)
		}

		commands[k] = tmpl
	}

	adp.commandsParsed = commands

	return nil
}

func transformWidget(widget string) (string, error) {
	cmd := exec.Command("node", "-e", TransformScript, widget)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		return stderr.String(), err
	}

	return fmt.Sprintf("function(device, command, showView) { return %s; }", strings.TrimSpace(stdout.String())), nil
}

func (adp *Adapter) getTemplateForCommand(command string) (*template.Template, error) {
	// Check if the command has been parsed first
	tmpl := adp.commandsParsed[command]

	// If not found, ensure the commands has been parsed
	if tmpl == nil && adp.Commands[command] != "" {
		if err := adp.parseCommands(); err != nil {
			return nil, err
		}
		tmpl = adp.commandsParsed[command]
	}

	// If still null, then the command does not exists in this adapter
	if tmpl == nil {
		return nil, ErrAdapterCommandNotFound
	}

	return tmpl, nil
}

// Execute an adapter command for the given device and parametrization.
func (adp *Adapter) Execute(shell []string, command string, device Device, params map[string]interface{}) (*ExecutionResult, error) {
	tmpl, err := adp.getTemplateForCommand(command)

	if err != nil {
		return nil, err
	}

	// Creates the execution context available in the command template
	ctx := newExecutionContext(adp, device, params)

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

	err = cmd.Run()

	result := &ExecutionResult{
		Success: cmd.ProcessState.Success(),
		Err:     strings.TrimSpace(stderr.String()),
		Out:     strings.TrimSpace(stdout.String()),
	}

	return result, err
}
