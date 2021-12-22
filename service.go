package main

import (
	"encoding/json"
	"fmt"
)

const CODE_SUCCESS = "SUCCESS"

type PlainSyspar struct {
	Id          string `json:"id"`
	Variable    string `json:"variable"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type SysparResponse struct {
	Code string      `json:"code"`
	Data PlainSyspar `json:"data"`
}

type Service interface {
	Set(syspar PlainSyspar) (*PlainSyspar, error)
	//Put(syspar *PlainSyspar) error
}

type service struct {
	restCall RestCall
}

func (s service) Set(syspar PlainSyspar) (*PlainSyspar, error) {
	bodyReq, err := s.toBytes(syspar)
	if err != nil {
		return nil, err
	}

	var bodyRes []byte
	if syspar.Id == "" {
		bodyRes, err = s.restCall.Post(bodyReq)
	} else {
		bodyRes, err = s.restCall.Put(syspar.Id, bodyReq)
	}

	if err != nil {
		return nil, err
	}

	return s.toPlainSyspar(bodyRes)
}

func (s service) toBytes(syspar PlainSyspar) ([]byte, error) {
	bodyReq, err := json.Marshal(syspar)
	if err != nil {
		return nil, fmt.Errorf("service::toBytes error json marshal, syspar: %+v, %+v", syspar, err)
	}
	return bodyReq, nil
}

func (s service) toPlainSyspar(bodyRes []byte) (*PlainSyspar, error) {
	var baseRes SysparResponse
	if err := json.Unmarshal(bodyRes, &baseRes); err != nil {
		return nil, fmt.Errorf("service::toPlainSyspar error json unmarshal, %+v", err)
	}

	if baseRes.Code != CODE_SUCCESS {
		return nil, fmt.Errorf("service::toPlainSyspar invalid SysparResponse code: %s", baseRes.Code)
	}

	return &baseRes.Data, nil
}

func NewService(restCall RestCall) Service {
	return service{restCall}
}
