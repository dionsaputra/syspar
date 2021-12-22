package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const CodeSuccess = "SUCCESS"

type (
	RestService interface {
		Set(syspar StringSyspar) (*StringSyspar, error)
	}

	rest struct {
		client http.Client
		config RestConfig
	}

	RestConfig struct {
		Url    string
		Header map[string]string
	}

	StringSyspar struct {
		Id          string `json:"id"`
		Variable    string `json:"variable"`
		Description string `json:"description"`
		Value       string `json:"value"`
	}

	SysparResponse struct {
		Code string       `json:"code"`
		Data StringSyspar `json:"data"`
	}
)

func NewService(client http.Client, config RestConfig) RestService {
	return rest{client, config}
}

func (r rest) Set(syspar StringSyspar) (*StringSyspar, error) {
	bodyReq, err := r.toBytes(syspar)
	if err != nil {
		return nil, err
	}

	bodyRes, err := r.httpCall(r.methodSet(syspar), syspar.Id, bodyReq)
	if err != nil {
		return nil, err
	}

	return r.toPlainSyspar(bodyRes)
}

func (r rest) toBytes(syspar StringSyspar) ([]byte, error) {
	bodyReq, err := json.Marshal(syspar)
	if err != nil {
		return nil, fmt.Errorf("rest::toBytes error json marshal, syspar: %+v, %+v", syspar, err)
	}
	return bodyReq, nil
}

func (r rest) methodSet(syspar StringSyspar) string {
	if syspar.Id == "" {
		return http.MethodPost
	}
	return http.MethodPut
}

func (r rest) Get() (*http.Response, error) {
	panic("implement me")
}

func (r rest) httpCall(method string, id string, bodyReq []byte) ([]byte, error) {
	req, err := r.buildRequest(method, id, bodyReq)
	if err != nil {
		return nil, err
	}

	res, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("rest::httpCall error client do request, method: %s, id: %s, req: %+v, %+v", method, id, req, err)
	}
	defer res.Body.Close()

	return r.parseResponse(res)
}

func (r rest) buildRequest(method string, id string, bodyReq []byte) (*http.Request, error) {
	url := r.config.Url
	if id != "" {
		url += fmt.Sprintf("/%s", id)
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(bodyReq))
	if err != nil {
		return nil, fmt.Errorf("rest::buildRequest error new request, method: %s, id: %s, %+v", method, id, err)
	}

	for key, val := range r.config.Header {
		req.Header.Add(key, val)
	}

	return req, nil
}

func (r rest) parseResponse(response *http.Response) ([]byte, error) {
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, fmt.Errorf("rest::parseResponse invalid response code, response: %+v", response)
	}

	bodyRes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("rest::parseResponse error read body, response: %+v, %+v", response, err)
	}

	return bodyRes, nil
}

func (r rest) toPlainSyspar(bodyRes []byte) (*StringSyspar, error) {
	var baseRes SysparResponse
	if err := json.Unmarshal(bodyRes, &baseRes); err != nil {
		return nil, fmt.Errorf("rest::toPlainSyspar error json unmarshal, %+v", err)
	}

	if baseRes.Code != CodeSuccess {
		return nil, fmt.Errorf("rest::toPlainSyspar invalid SysparResponse code: %r", baseRes.Code)
	}

	return &baseRes.Data, nil
}
