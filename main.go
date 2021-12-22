package main

import (
	"log"
	"net/http"
)

func main() {
	var config Config
	if err := LoadConfig(&config); err != nil {
		log.Fatalln(err)
	}

	restCall := NewRestCall(http.Client{}, config)
	service := NewService(restCall)

	syspar := PlainSyspar{
		Id:          "61c2eb1f30233ab23447bc3a",
		Variable:    "testing14",
		Value:       "1, 2, 3, 4",
		Description: "abc",
	}

	res, err := service.Set(syspar)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println(res)
}
