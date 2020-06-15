package v1

import (
	"encoding/json"
	"fmt"
	"github.com/freshapi/protopub/schema"
)

// Config is implementation of protopub image config
type Config struct {
	V          schema.Version `json:"freshapi.protopub.version"`
	ProtoFiles []string       `json:"freshapi.protopub.files"`
}

// Files returns all files
func (c *Config) Files() []string {
	return c.ProtoFiles
}

// SetFiles sets files
func (c *Config) SetFiles(files []string) {
	c.ProtoFiles = files
}

// NewConfig creates v1 configuration
func NewConfig() schema.Config {
	return &Config{
		V: schema.V1,
	}
}

// ParseConfig parses json configuration into v1 config
func ParseConfig(data []byte) (schema.Config, error) {
	var c Config
	err := json.Unmarshal(data, &c)
	if err != nil {
		return nil, fmt.Errorf("error parsing OCI image config: %w", err)
	}
	if c.V != schema.V1 {
		return nil, fmt.Errorf("invalid config version")
	}
	return &c, nil
}
