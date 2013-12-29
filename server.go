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
	if cmd != "" {
		url = self.url
	} else {
		url = fmt.Sprintf("%s/%s", self.url, cmd)
	}

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

func (self *Server) HasIndex(index string) error {

	if resp, err := self.Get(index); err != nil {
		return err
	} else if resp.StatusCode > 399 {
		return fmt.Errorf("Response status %s", resp.Status)
	} else {
		fmt.Printf("Response => %s\n", resp)
		defer resp.Body.Close()
		return nil
	}
}

func (self *Server) CreateIndex(index string, body io.Reader) error {

	if resp, err := self.Put(index, body); err != nil {
		return err
	} else if resp.StatusCode > 399 {
		return fmt.Errorf("Response status %s", resp.Status)
	} else {
		defer resp.Body.Close()
		return nil
	}
}

func (self *Server) DeleteIndex(index string) error {

	if resp, err := self.Delete(index); err != nil {
		return err
	} else if resp.StatusCode > 399 {
		return fmt.Errorf("Response status %s", resp.Status)
	} else {
		defer resp.Body.Close()
		return nil
	}
}
