package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RestCall interface {
	//Get() (*http.Response, error)
	Post(bodyReq []byte) ([]byte, error)
	Put(id string, bodyReq []byte) ([]byte, error)
}

type restCall struct {
	client http.Client
	config Config
}

func (r restCall) Get() (*http.Response, error) {
	panic("implement me")
}

func (r restCall) Post(bodyReq []byte) ([]byte, error) {
	return r.httpCall(http.MethodPost, "", bodyReq)
}

func (r restCall) Put(id string, bodyReq []byte) ([]byte, error) {
	return r.httpCall(http.MethodPut, id, bodyReq)
}

func (r restCall) httpCall(method string, id string, bodyReq []byte) ([]byte, error) {
	req, err := r.buildRequest(method, id, bodyReq)
	if err != nil {
		return nil, err
	}

	res, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("httpCall error client do request, method: %s, id: %s, req: %+v, %+v", method, id, req, err)
	}
	defer res.Body.Close()

	return r.parseResponse(res)
}

func (r restCall) buildRequest(method string, id string, bodyReq []byte) (*http.Request, error) {
	url := r.config.Url
	if id != "" {
		url += fmt.Sprintf("/%s", id)
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(bodyReq))
	if err != nil {
		return nil, fmt.Errorf("buildRequest error new request, method: %s, id: %s, %+v", method, id, err)
	}

	for key, val := range r.config.Header {
		req.Header.Add(key, val)
	}

	return req, nil
}

func (r restCall) parseResponse(response *http.Response) ([]byte, error) {
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, fmt.Errorf("parseResponse invalid response code, response: %+v", response)
	}

	bodyRes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("parseResponse error read body, response: %+v, %+v", response, err)
	}

	return bodyRes, nil
}

func NewRestCall(client http.Client, config Config) RestCall {
	return restCall{client, config}
}
