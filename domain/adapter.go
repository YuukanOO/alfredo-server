package domain

// Adapter represents an available smart adapter.
type Adapter struct {
	Name     string
	Category string
	Commands map[string]string
}
