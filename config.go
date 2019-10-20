package main

import (
	"GOFirst/app"
	"GOFirst/server"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Logger *LoggerConfig  `yaml:"logger"`
	Server *server.Config `yaml:"server"`
	App    *app.Config    `yaml:"app"`
}

func ReadConfig(name string) (*Config, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	conf := &Config{}
	if err := yaml.Unmarshal(data, conf); err != nil {
		return nil, err
	}
	return conf, nil
}
