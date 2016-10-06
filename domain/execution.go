package domain

// ExecutionContext represents a device execution context
type ExecutionContext struct {
	Config map[string]interface{}
	Params map[string]interface{}
}

func newExecutionContext(
	config map[string]interface{},
	params map[string]interface{}) *ExecutionContext {
	return &ExecutionContext{
		Config: config,
		Params: params,
	}
}

// ExecutionResult holds information about the running of an adapter command.
type ExecutionResult struct {
	Success bool   `json:"success"`
	Out     string `json:"out"`
	Err     string `json:"err"`
}
