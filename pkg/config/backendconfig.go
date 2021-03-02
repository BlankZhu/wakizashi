package config

// BackendConfig describes the configuration for data storage backend
type BackendConfig struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	URI      string `yaml:"url"` // connection URI, for mongoDB only
	Timeout  uint64 `yaml:"timeout"`
}
