package gosearch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Status struct {
	Status  int
	Name    string
	Version struct {
		Number   string
		Snapshot bool `json:"snapshot_build"`
	}
}

type Server struct {
	url string
}

func Connect(host string, port int) *Server {
	return ConnectURL(fmt.Sprintf("http://%s:%d", host, port))
}

func ConnectURL(url string) *Server {
	es := &Server{url: url}
	return es
}

func (self *Server) do(method, cmd string, body io.Reader) (*http.Response, error) {
	var url string
	if cmd == "" {
		url = self.url
	} else {
		url = fmt.Sprintf("%s/%s", self.url, cmd)
	}

	fmt.Printf("Command => %s:%s\n", method, url)
	if req, err := http.NewRequest(method, url, body); err != nil {
		return nil, err
	} else {
		return http.DefaultClient.Do(req)
	}

}

func (self *Server) Get(cmd string) (*http.Response, error) {
	return self.do("GET", cmd, nil)
}

func (self *Server) Put(cmd string, body io.Reader) (*http.Response, error) {
	return self.do("PUT", cmd, body)
}

func (self *Server) Post(cmd string, body io.Reader) (*http.Response, error) {
	return self.do("POST", cmd, body)
}

func (self *Server) Delete(cmd string) (*http.Response, error) {
	return self.do("DELETE", cmd, nil)
}

func (self *Server) Head(cmd string) (*http.Response, error) {
	return self.do("HEAD", cmd, nil)
}

func (self *Server) Status() (*Status, error) {
	status := &Status{}

	if resp, err := self.Get(""); err != nil {
		return nil, err
	} else if resp.StatusCode > 399 {
		return nil, fmt.Errorf("Response status %s", resp.Status)
	} else {
		defer resp.Body.Close()
		err := json.NewDecoder(resp.Body).Decode(status)
		return status, err
	}
}

func (self *Server) HasIndex(index string) bool {

	if resp, err := self.Head(index); err != nil {
		// FIXME Log error
		return false
	} else {
		return resp.StatusCode == 200
	}
}

func (self *Server) CreateIndex(index string) bool {

	if resp, err := self.Put(index, nil); err != nil {
		// FIXME Log error
		return false
	} else {
		return resp.StatusCode == 200
	}
}

func (self *Server) CreateIndexWithSettings(index string, settings io.Reader) bool {
	if resp, err := self.Put(index, settings); err != nil {
		// FIXME Log error
		return false
	} else {
		return resp.StatusCode == 200
	}
}

func (self *Server) DeleteIndex(index string) bool {

	if resp, err := self.Delete(index); err != nil {
		// FIXME Log error
		return false
	} else {
		return resp.StatusCode == 200
	}
}
