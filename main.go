package main

import (
	"encoding/json"
	"log"
	"net/http"
	"syspar/file"
	"syspar/rest"
)

func main() {
	var config Config
	if err := LoadConfig(&config); err != nil {
		log.Fatalln(err)
	}

	fileService := file.NewService(config.FileConfig)
	restService := rest.NewService(http.Client{}, config.RestConfig)

	jsonSyspar, err := fileService.Read()
	if err != nil {
		log.Fatalln(err.Error())
	}

	byteVal, err := json.Marshal(jsonSyspar.Value)
	if err != nil {
		log.Fatalln(err.Error())
	}

	res, err := restService.Set(rest.StringSyspar{
		Id:          jsonSyspar.Id,
		Variable:    jsonSyspar.Variable,
		Value:       string(byteVal),
		Description: jsonSyspar.Description,
	})
	if err != nil {
		log.Fatalln(err.Error())
	}

	jsonSyspar.Id = res.Id

	log.Fatalln(fileService.Write(*jsonSyspar))
}
