package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"syspar/file"
	"syspar/rest"
)

type Config struct {
	RestConfig rest.RestConfig `json:"restConfig"`
	FileConfig file.FileConfig `json:"fileConfig"`
}

func LoadConfig(config *Config) error {
	data, err := ioutil.ReadFile("syspar.config.json")
	if err != nil {
		log.Println("[LoadConfig] error read config file ", err)
		return err
	}

	if err := json.Unmarshal(data, config); err != nil {
		log.Println("[LoadConfig] error unmarshal data ", err)
		return err
	}

	return nil
}
