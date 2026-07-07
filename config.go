package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InactivitySeconds int    `yaml:"inactivity_seconds"`
	StateFile         string `yaml:"state_file"`
	Rules             []Rule `yaml:"rules"`
}

type Rule struct {
	Name              string   `yaml:"name"`
	LogDir            string   `yaml:"log_dir"`
	FileGlob          string   `yaml:"file_glob"`
	Patterns          []string `yaml:"patterns"`
	Actions           []string `yaml:"actions"`
	CooldownSeconds   int      `yaml:"cooldown_seconds"`
	InactivitySeconds int      `yaml:"inactivity_seconds,omitempty"` // per-rule override
	SampleLength      int      `yaml:"sample_length,omitempty"`      // chars to include in alert
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	cfg := &Config{
		InactivitySeconds: 300,
		StateFile:         "state.json",
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if len(cfg.Rules) == 0 {
		return nil, fmt.Errorf("no rules defined in config")
	}

	for i, rule := range cfg.Rules {
		if rule.SampleLength == 0 {
			rule.SampleLength = 32 // default
		}
		if rule.Name == "" {
			return nil, fmt.Errorf("rule %d has no name", i)
		}
		if rule.LogDir == "" {
			return nil, fmt.Errorf("rule %q has no log_dir", rule.Name)
		}
		if len(rule.Patterns) == 0 {
			return nil, fmt.Errorf("rule %q has no patterns", rule.Name)
		}
		if len(rule.Actions) == 0 {
			return nil, fmt.Errorf("rule %q has no actions", rule.Name)
		}
		if rule.CooldownSeconds <= 0 {
			rule.CooldownSeconds = 300
		}
		if rule.InactivitySeconds == 0 {
			rule.InactivitySeconds = cfg.InactivitySeconds
		}
		cfg.Rules[i] = rule
	}

	return cfg, nil
}
