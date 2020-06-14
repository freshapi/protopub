package schema

import (
	"encoding/json"
)

type Config interface {
	Files() []string
	SetFiles([]string)
}

func RenderConfig(config Config) ([]byte, error) {
	return json.Marshal(config)
}
