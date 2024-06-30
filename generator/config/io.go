package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"path"
)

func ReadConfig(folder string) (*Config, error) {
	c := Config{}

	cfg := path.Join(folder, "config.yaml")
	d, err := os.ReadFile(cfg)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(d, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func SaveConfig(c *Config, folder string) error {
	cfg := path.Join(folder, "config.yaml")

	r, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(cfg, r, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func InitConfig(folder string) error {
	cfg := path.Join(folder, "config.yaml")

	c := DefaultConfig()

	r, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(cfg, r, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
