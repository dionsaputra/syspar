package file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type (
	FileService interface {
		Read() (*JsonSyspar, error)
		Write(syspar JsonSyspar) error
	}

	fileService struct {
		config FileConfig
	}

	FileConfig struct {
		Editor string `json:"editor"`
	}

	JsonSyspar struct {
		Id          string                 `json:"id"`
		Variable    string                 `json:"variable"`
		Description string                 `json:"description"`
		Value       map[string]interface{} `json:"value"`
	}
)

func NewService(config FileConfig) FileService {
	return fileService{config}
}

func (f fileService) Read() (*JsonSyspar, error) {
	data, err := ioutil.ReadFile(f.config.Editor)
	if err != nil {
		return nil, fmt.Errorf("file::Read error ReadFile, %+v", err)
	}

	var res JsonSyspar
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("file::Read error Unmarshal, %+v", err)
	}

	return &res, nil
}

func (f fileService) Write(syspar JsonSyspar) error {
	data, err := json.MarshalIndent(syspar, "", "\t")
	if err != nil {
		return fmt.Errorf("file::Write error MarshalIndent, syspar: %+v", syspar)
	}
	return ioutil.WriteFile(f.config.Editor, data, 0644)
}
