package config

type Config struct {
	ProjectName string `yaml:"project_name"`
}

func DefaultConfig() *Config {
	return &Config{
		ProjectName: "Docute",
	}
}
