package domain

// Widget represents an adapter widget
type Widget struct {
	ID       string
	Adapter  string
	reactStr string
}

func newWidget(adapter string, id string, reactStr string) *Widget {
	return &Widget{
		ID:       id,
		Adapter:  adapter,
		reactStr: reactStr,
	}
}
