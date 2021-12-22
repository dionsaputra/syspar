package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Url    string            `json:"url"`
	Header map[string]string `json:"header"`
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
