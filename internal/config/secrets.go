package config

import (
	"fmt"
	"os"
	"strings"
)

const (
	// noenv resolves to nothing. This means no environment variables
	EnvVarsTypeNoEnv = "noenv"
	// plain resolves to the declared values.
	EnvVarsTypePlain = "plain"
	// file resolves to values read from a file.
	EnvVarsTypeFile = "file"
	// injected resolves to values read from the environment (os.LookupEnv).
	EnvVarsTypeInjected = "injected"
)

// Resolve turns the declared secrets block into actual values, per the
// `type` semantics in docs/schema.md#variablesyaml. Returns (nil, nil) for
// type: noenv.
func (e *Secrets) Resolve() (map[string]string, error) {
	switch e.Type {
	case EnvVarsTypeNoEnv:
		return nil, nil
	case EnvVarsTypePlain:
		return e.Values, nil
	case EnvVarsTypeFile:
		return resolveFromFile(e.Filepath, e.Values)
	case EnvVarsTypeInjected:
		return resolveFromEnv(e.Values)
	default:
		return nil, fmt.Errorf("secrets: unknown type %q (accepted: %s | %s | %s | %s)", e.Type, EnvVarsTypeNoEnv, EnvVarsTypePlain, EnvVarsTypeFile, EnvVarsTypeInjected)
	}
}

// resolveFromFile reads a minimal KEY=VALUE file (like .env) — no need for
// a parsing dependency for something this small.
func resolveFromFile(path string, secrets map[string]string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("secrets: reading %s: %w", path, err)
	}

	values := map[string]string{}
	for line := range strings.SplitSeq(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		values[strings.TrimSpace(key)] = strings.Trim(strings.TrimSpace(val), `"'`)
	}

	resolved := make(map[string]string, len(secrets))
	for name, key := range secrets {
		v, ok := values[key]
		if !ok {
			return nil, fmt.Errorf("secrets: key %q not found in %s (referenced by %q)", key, path, name)
		}
		resolved[name] = v
	}
	return resolved, nil
}

func resolveFromEnv(secrets map[string]string) (map[string]string, error) {
	resolved := make(map[string]string, len(secrets))
	for name, key := range secrets {
		v, ok := os.LookupEnv(key)
		if !ok {
			return nil, fmt.Errorf("secrets: %q is not set in the environment (referenced by %q)", key, name)
		}
		resolved[name] = v
	}
	return resolved, nil
}
