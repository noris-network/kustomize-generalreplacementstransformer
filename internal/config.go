package grt

import (
	"fmt"
	"os"

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type Resource struct {
	Kind string `yaml:"kind"`
	Name string `yaml:"name"`
}

type Replacement struct {
	Resource   Resource   `yaml:"resource"`
	Type       string     `yaml:"type"`
	Delimiters *[2]string `yaml:"delimiters"`
}

type SelectValue struct {
	Name     string `yaml:"name"`
	Default  string `yaml:"default"`
	Splat    bool   `yaml:"splat"`
	Resource struct {
		Kind      string `yaml:"kind"`
		Name      string `yaml:"name"`
		FieldPath string `yaml:"fieldPath"`
	} `yaml:"resource"`
}

type Config struct {
	APIVersion string     `yaml:"apiVersion"`
	Kind       string     `yaml:"kind"`
	Delimiters *[2]string `yaml:"delimiters"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Values       map[string]any `yaml:"values"`
	ValuesFile   string         `yaml:"valuesFile"`
	SelectValues []SelectValue  `yaml:"selectValues"`
	Replacements []Replacement  `yaml:"replacements"`
}

func (t *Transformer) configure(config []byte) error {

	cfg := Config{}
	err := yaml.Unmarshal(config, &cfg)
	if err != nil {
		return fmt.Errorf("unmarshal config: %v", err)
	}

	if cfg.Values == nil {
		cfg.Values = map[string]any{}
	}

	if cfg.ValuesFile != "" {
		data, err := os.ReadFile(cfg.ValuesFile)
		if err != nil {
			return fmt.Errorf("readFile: %v", err)
		}
		fileValues := map[string]any{}
		err = yaml.Unmarshal(data, &fileValues)
		if err != nil {
			return fmt.Errorf("unmarshal fileValues: %v", err)
		}
		for k, v := range fileValues {
			if cfg.Values[k] == nil {
				cfg.Values[k] = v
			}
		}
	}

	t.config = cfg

	return nil
}
