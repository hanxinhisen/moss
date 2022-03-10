package config

import "github.com/hanxinhisen/moss/internal/userservice/options"

type Config struct {
	*options.Options
}

func CreateConfigFromOptions(opts *options.Options) (*Config, error) {
	return &Config{opts}, nil
}
