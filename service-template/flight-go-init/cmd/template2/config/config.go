package config

import (
	"bytes"
	"errors"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/config"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

type Commission struct {
	Type  CommissionType `yaml:"type"`
	Value float64        `yaml:"value"`
}

type CommissionType string

type (
	Config struct {
		config.BaseConfig       `yaml:",inline"`
		{{CONFIG_PROPERTY_STRUCT}} {{CONFIG_PROPERTY_STRUCT}} `yaml:"{{PROVIDER}}_config_properties"`
	}

	{{CONFIG_PROPERTY_STRUCT}} struct {
		Commission Commission `yaml:"commission"`
		PromoCode  string     `yaml:"promo_code"`
	}
)

func New() (*Config, error) {
	return NewFromPath(config.GetConfigPath())
}

func NewFromPath(configPath string) (*Config, error) {
	var cfg Config

	// Resolve absolute path
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, errors.New("failed to resolve config path: " + err.Error())
	}

	file, err := os.ReadFile(absPath)
	if err != nil {
		return nil, errors.New("failed to read config file '" + absPath + "': " + err.Error())
	}

	file = []byte(os.ExpandEnv(string(file)))
	decoder := yaml.NewDecoder(bytes.NewReader(file))
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, errors.New("failed to decode config file: " + err.Error())
	}
	// Validate mandatory fields
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}
