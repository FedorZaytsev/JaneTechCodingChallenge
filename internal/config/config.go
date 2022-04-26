package config

import (
	"io/ioutil"

	"jane_tech/internal/logger"

	yaml "gopkg.in/yaml.v2"
)

type DbConfig struct {
}

type ServerConfig struct {
	ServeAddress  string
	ImpressionUrl string
}

//Config is needed for configure programm.
type Config struct {
	Title    string
	Log      logger.LogConfig
	Server   ServerConfig
	Database struct {
		Type   string
		Config DbConfig
	}
}

func Configure(fileName string) (*Config, error) {
	var cnf Config
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return &Config{}, err
	}
	err = yaml.Unmarshal(data, &cnf)
	if err != nil {
		return &Config{}, err
	}
	return &cnf, nil
}
