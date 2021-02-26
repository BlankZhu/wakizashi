package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// ProbeConfig describe the configuration for traffic collecting probe
type ProbeConfig struct {
	CenterAddr  string   `yaml:"centerAddr"`            // center's address
	LogLev      int      `yaml:"logLev"`                // log level
	DumpDir     string   `yaml:"dumpDir"`               // directory for temp dumping
	NetworkDevs []string `yaml:"networkDevs"`           // network devices' name where the probe work
	AutoClear   bool     `yaml:"autoClear,omitempty"`   // decide if remove the caputre file or not automatically
	CapInterval int      `yaml:"capInterval,omitempty"` // interval of rotating dump file, in second; if non-positive, use 1
	UploadRetry int      `yaml:"uploadRetry,omitempty"` // count of retry to upload traffic status to center
}

// LoadConfigFromYAML load config from given path
func (pc *ProbeConfig) LoadConfigFromYAML(path string) error {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, pc)
	if err != nil {
		return err
	}
	if pc.CapInterval <= 0 {
		pc.CapInterval = 1
	}
	if pc.UploadRetry <= 0 {
		pc.UploadRetry = 0
	}
	return nil
}

// ToString return a string representing the config
func (pc ProbeConfig) ToString() string {
	ret := fmt.Sprintf("%+v", pc)
	return ret
}

// CreateDumpDir create dump directory if not exists
func (pc ProbeConfig) CreateDumpDir() error {
	err := os.Mkdir(pc.DumpDir, os.ModeDir)
	if err != nil || os.IsExist(err) {
		return nil
	}
	return err
}
