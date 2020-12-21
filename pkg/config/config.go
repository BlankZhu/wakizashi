// Package config describe the config used by wakizashi
package config

// BaseConfig is an interface for wakizashi's probe and center's config
type BaseConfig interface {
	// LoadConfigFromYAML load config from yaml file
	LoadConfigFromYAML(path string) error
	// ToString return a string representing the config
	ToString() string
}
