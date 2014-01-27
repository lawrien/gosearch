package gosearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type httpVerb string

const (
	GET    httpVerb = "GET"
	PUT             = "PUT"
	POST            = "POST"
	DELETE          = "DELETE"
	HEAD            = "HEAD"
)

const DEFAULT_POOL = 200

var DefaultConnectionPool *HttpClientPool

func init() {
	DefaultConnectionPool = NewPool(DEFAULT_POOL)
}

type Response struct {
	Status int
	Body   []byte
}

func (self *Response) IsSuccess() bool {
	return self.Status >= 200 && self.Status < 400
}

func (self *Response) Convert(v interface{}) error {
	return json.Unmarshal(self.Body, v)
}

type HttpClientPool struct {
	clients chan *http.Client
}

func NewPool(size int) *HttpClientPool {
	pool := new(HttpClientPool)
	pool.clients = make(chan *http.Client, size)

	tr := &http.Transport{DisableKeepAlives: false, MaxIdleConnsPerHost: size * 2}
	for i := 0; i < size; i++ {
		c := &http.Client{Transport: tr}
		pool.clients <- c
	}

	return pool
}

func (self *HttpClientPool) getClient() *http.Client {
	client := <-self.clients
	return client
}

func (self *HttpClientPool) returnClient(client *http.Client) {
	self.clients <- client
}

func (self *HttpClientPool) Do(method httpVerb, url string, in interface{}) (*Response, error) {
	content, merr := json.Marshal(in)
	fmt.Printf("Do Content => %s\n", string(content))
	if merr != nil {
		return nil, merr
	}

	if req, err := http.NewRequest(string(method), url, bytes.NewReader(content)); err != nil {
		return nil, err
	} else {
		client := self.getClient()
		defer self.returnClient(client)

		var resp *http.Response
		if resp, err = client.Do(req); err != nil {
			return nil, err
		}

		defer resp.Body.Close()
		if bytes, err := ioutil.ReadAll(resp.Body); err != nil {
			return nil, err
		} else {
			return &Response{Status: resp.StatusCode, Body: bytes}, nil
		}
	}

}
