package v1

import (
	"encoding/json"
	"fmt"
	"github.com/freshapi/protopub/schema"
)

type Config struct {
	V          schema.Version `json:"freshapi.protopub.version"`
	ProtoFiles []string       `json:"freshapi.protopub.files"`
}

func (c *Config) Files() []string {
	return c.ProtoFiles
}

func (c *Config) SetFiles(files []string) {
	c.ProtoFiles = files
}

func NewConfig() schema.Config {
	return &Config{
		V: schema.V1,
	}
}

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
