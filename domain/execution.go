package domain

import "bytes"

// ExecutionContext represents a device execution context
type ExecutionContext struct {
	Config  map[string]interface{}
	Params  map[string]interface{}
	adapter *Adapter
}

func newExecutionContext(
	adapter *Adapter,
	config map[string]interface{},
	params map[string]interface{}) *ExecutionContext {
	return &ExecutionContext{
		Config:  config,
		Params:  params,
		adapter: adapter,
	}
}

// Command append the given adapter cmd name using the same execution context.
func (ctx *ExecutionContext) Command(cmd string) (string, error) {
	tmpl, err := ctx.adapter.getTemplateForCommand(cmd)

	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, ctx); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ExecutionResult holds information about the running of an adapter command.
type ExecutionResult struct {
	Success bool   `json:"success"`
	Out     string `json:"out"`
	Err     string `json:"err"`
}
