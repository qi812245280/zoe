package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Debug     bool   `yaml:"debug" binding:"required"`
	Name      string `yaml:"name" binding:"required"`
	LogFormat string `yaml:"logformat" binding:"required"`
	Listen    string `yaml:"listen" binding:"required"`
	Database  struct {
		Engine           string `yaml:"engine" binding:"required"`
		ConnectionString string `yaml:"connectionstring" binding:"required"`
	} `yaml:"database"`
}

var C *Config

func LoadConfig(path string) error {
	var config Config
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		return err
	}
	C = &config

	return nil
}
