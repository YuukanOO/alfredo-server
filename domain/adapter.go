package domain

// Adapter represents an available smart adapter.
type Adapter struct {
	Name     string
	Category string
	Config   map[string]interface{}
	Commands map[string]string
}
