// Package config loads and resolves main.yaml/variables.yaml per
// docs/schema.md. It only knows the shape of the YAML — it doesn't know
// how to run a flow (see internal/pipeline for that).
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Stage struct {
	Type string `yaml:"type"`
}

type Step struct {
	Uses string         `yaml:"uses"`
	With map[string]any `yaml:"with,omitempty"`
	Via  []string       `yaml:"via,omitempty"`
	On   []int          `yaml:"on,omitempty"`
}

// Secrets mirrors the `secrets` block, whether it lives in main.yaml or
// variables.yaml — same shape either way, see docs/schema.md.
type Secrets struct {
	Type     string            `yaml:"type"`
	Filepath string            `yaml:"filepath,omitempty"`
	Values   map[string]string `yaml:"values,omitempty"`
}

type Main struct {
	Name   string `yaml:"name"`
	Config struct {
		Input  Stage `yaml:"input"`
		Output Stage `yaml:"output"`
	} `yaml:"config"`
	Secrets *Secrets `yaml:"secrets,omitempty"`
	Steps   []Step   `yaml:"steps"`
}

type Variables struct {
	Secrets *Secrets `yaml:"secrets"`
}

func LoadMain(path string) (*Main, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	var m Main
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	if len(m.Steps) == 0 {
		return nil, fmt.Errorf("%s: at least one step is required", path)
	}
	return &m, nil
}

// LoadVariables returns (nil, nil) when the file doesn't exist — it's
// optional, see docs/schema.md.
func LoadVariables(path string) (*Variables, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	var v Variables
	if err := yaml.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	return &v, nil
}

// Load reads both files main.yaml needs; variables.yaml is optional.
func Load(mainPath, varsPath string) (*Main, *Variables, error) {
	m, err := LoadMain(mainPath)
	if err != nil {
		return nil, nil, err
	}
	v, err := LoadVariables(varsPath)
	if err != nil {
		return nil, nil, err
	}
	return m, v, nil
}

// ResolveSecrets implements the precedence from
// docs/schema.md#secrets-resolution: main.yaml wins if present, then
// variables.yaml; it's an error to define it in both, or in neither.
func ResolveSecrets(m *Main, v *Variables) (*Secrets, error) {
	mainHas := m.Secrets != nil
	varsHas := v != nil && v.Secrets != nil

	switch {
	case mainHas && varsHas:
		return nil, fmt.Errorf("secrets is defined in both main.yaml and variables.yaml; pick one")
	case mainHas:
		return m.Secrets, nil
	case varsHas:
		return v.Secrets, nil
	default:
		return nil, fmt.Errorf("secrets is missing: define it in main.yaml or variables.yaml (or set type: noenv)")
	}
}
