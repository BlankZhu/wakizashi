package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// CenterConfig describe the configuration for traffic convergent center
type CenterConfig struct {
	LogLev        int           `yaml:"logLev"`        // log level
	Port          uint16        `yaml:"port"`          // port to listen for grpc
	HealthPort    uint16        `yaml:"healthPort"`    // port for health probe
	RecovDir      string        `yaml:"recoverDir"`    // directory to store the recovery info
	BackendConfig BackendConfig `yaml:"backendConfig"` // configuration for specific data storage backend
}

// LoadConfigFromYAML load config from given path
func (cc *CenterConfig) LoadConfigFromYAML(path string) error {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, cc)
	if err != nil {
		return err
	}
	return nil
}

// ToString return a string representing the config
func (cc CenterConfig) ToString() string {
	ret := fmt.Sprintf("%+v", cc)
	return ret
}

// CreateRecoveryDir create dump directory if not exists
func (cc CenterConfig) CreateRecoveryDir() error {
	err := os.Mkdir(cc.RecovDir, os.ModeDir)
	if err != nil || os.IsExist(err) {
		return nil
	}
	return err
}
