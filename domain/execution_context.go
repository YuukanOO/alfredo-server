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
