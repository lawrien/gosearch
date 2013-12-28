package gosearch

import (
	"encoding/json"
	"fmt"
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

func (self *Server) Status() (*Status, error) {
	status := &Status{}

	if resp, err := http.Get(self.url); err != nil {
		return nil, err
	} else if resp.StatusCode > 399 {
		return nil, fmt.Errorf("Response status %s", resp.Status)
	} else {
		defer resp.Body.Close()
		err := json.NewDecoder(resp.Body).Decode(status)
		return status, err
	}
}
