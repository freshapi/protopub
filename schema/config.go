package schema

import (
	"encoding/json"
)

// Config is an interface representing protopub image configuration
type Config interface {
	Files() []string
	SetFiles([]string)
}

// RenderConfig renders Config
func RenderConfig(config Config) ([]byte, error) {
	return json.Marshal(config)
}
