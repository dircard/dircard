package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	defaultFileSize    = 10
	defaultLineStart   = 0
	defaultLineCount   = 10
	defaultSearchDepth = 5
)

type ConfigOption func(*Config)

type Config struct {
	CandidateOrder []string `yaml:"candidate_order"`
	FileSizeKB     int      `yaml:"file_size_kb"`
	LineStart      *int     `yaml:"lineStart"`
	LineCount      *int     `yaml:"lineCount"`
	Depth          *int     `yaml:"depth"`
}

func intPtr(v int) *int {
	return &v
}

func configFilePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "dircard", "config.yaml"), nil
}

func Load() *Config {
	cfg := &Config{}
	path, err := configFilePath()
	if err != nil {
		return ApplyDefaultValues(cfg)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return ApplyDefaultValues(cfg)
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return ApplyDefaultValues(cfg)
	}
	return ApplyDefaultValues(cfg)
}

func ApplyDefaultValues(c *Config) *Config {
	if c.FileSizeKB <= 0 {
		c.FileSizeKB = defaultFileSize
	}
	if c.LineStart == nil || *c.LineStart < 0 {
		c.LineStart = intPtr(defaultLineStart)
	}
	if c.LineCount == nil || *c.LineCount < 0 {
		c.LineCount = intPtr(defaultLineCount)
	}
	if c.Depth == nil || *c.Depth < 0 {
		c.Depth = intPtr(defaultSearchDepth)
	}
	return c
}

func WithCandidateOrder(order []string) ConfigOption {
	return func(c *Config) {
		c.CandidateOrder = order
	}
}

func WithFileSizeKB(size int) ConfigOption {
	return func(c *Config) {
		c.FileSizeKB = size
	}
}

func WithLineStart(start int) ConfigOption {
	return func(c *Config) {
		c.LineStart = &start
	}
}

func WithLineCount(count int) ConfigOption {
	return func(c *Config) {
		c.LineCount = &count
	}
}

func WithDepth(depth int) ConfigOption {
	return func(c *Config) {
		c.Depth = &depth
	}
}

func ApplyOptions(opts ...ConfigOption) error {
	cfg := Load()
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg.Save()
}

func (c *Config) Save() error {
	path, err := configFilePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
