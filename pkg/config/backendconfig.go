package config

import (
	"BlankZhu/wakizashi/pkg/constant"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// BackendConfig describes the configuration for data storage backend
type BackendConfig struct {
	Type      string       `yaml:"type"`         // backend type
	Timeout   uint         `yaml:"timeout"`      // connection timeout
	Database  string       `yaml:"database"`     // database to store traffic records
	Table     string       `yaml:"table"`        // table(collection in Mongo or measurement in influx) to store traffic records
	InfluxCfg InfluxConfig `yaml:"influxConfig"` // influxdb connection config section
	RedisCfg  RedisConfig  `yaml:"redisConfig"`  // redis connection config section
	MongoCfg  MongoConfig  `yaml:"mongoConfig"`  // mongo connection config section
}

func NewBackendConfig() *BackendConfig {
	ret := &BackendConfig{}
	ret.Database = constant.WakizashiDefaultDatabase
	ret.Table = constant.WakizashiDefaultTable
	return ret
}

func (bc *BackendConfig) LoadConfigFromYAML(path string) error {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, bc)
	if err != nil {
		return err
	}
	return nil
}

func (bc BackendConfig) ToString() string {
	ret := fmt.Sprintf("%+v", bc)
	return ret
}
