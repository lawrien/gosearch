package gosearch

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
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

const DEFAULT_IDLE_TIMEOUT = time.Second
const MAX_CONNECTIONS = 100
const CONN_TIMEOUT = time.Second

func timeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, nil
	}
}

var tr = &http.Transport{
	DisableKeepAlives:   false,
	MaxIdleConnsPerHost: MAX_CONNECTIONS * 2,
	Dial:                timeoutDialer(time.Second, CONN_TIMEOUT),
}

var client = &http.Client{Transport: tr}
var throttle = make(chan int, MAX_CONNECTIONS)

func (self *Server) do(method, cmd string, body io.Reader, fn func(*http.Response) error) error {
	var url string
	if cmd == "" {
		url = self.url
	} else {
		url = fmt.Sprintf("%s/%s", self.url, cmd)
	}

	// fmt.Printf("Command => %s:%s\n", method, url)
	if req, err := http.NewRequest(method, url, body); err != nil {
		return err
	} else {
		throttle <- 1
		defer func() { <-throttle }()
		if resp, err := client.Do(req); err != nil {
			return err
		} else {
			defer resp.Body.Close()
			return fn(resp)
		}
	}
}

func (self *Server) Get(cmd string, fn func(*http.Response) error) error {
	return self.do("GET", cmd, nil, fn)
}

func (self *Server) Put(cmd string, body io.Reader, fn func(*http.Response) error) error {
	return self.do("PUT", cmd, body, fn)
}

func (self *Server) Post(cmd string, body io.Reader, fn func(*http.Response) error) error {
	return self.do("POST", cmd, body, fn)
}

func (self *Server) Delete(cmd string, fn func(*http.Response) error) error {
	return self.do("DELETE", cmd, nil, fn)
}

func (self *Server) Head(cmd string, fn func(*http.Response) error) error {
	return self.do("HEAD", cmd, nil, fn)
}

func (self *Server) Status() (*Status, error) {
	status := &Status{}

	return status, self.Get("", func(resp *http.Response) error {
		if resp.StatusCode != 200 {
			return fmt.Errorf("Unexpected status code %d", resp.StatusCode)
		} else {
			return json.NewDecoder(resp.Body).Decode(status)
		}
	})
}

func (self *Server) HasIndex(index string) bool {

	return self.Head(index, func(resp *http.Response) error {
		if resp.StatusCode != 200 {
			return fmt.Errorf("Unexpected status code %d", resp.StatusCode)
		} else {
			return nil
		}
	}) == nil
}

func (self *Server) CreateIndex(index string) bool {

	return self.Put(index, nil, func(resp *http.Response) error {
		if resp.StatusCode != 200 {
			return fmt.Errorf("Unexpected status code %d", resp.StatusCode)
		} else {
			return nil
		}
	}) == nil
}

func (self *Server) CreateIndexWithSettings(index string, settings io.Reader) bool {
	return self.Put(index, settings, func(resp *http.Response) error {
		if resp.StatusCode != 200 {
			return fmt.Errorf("Unexpected status code %d", resp.StatusCode)
		} else {
			return nil
		}
	}) == nil
}

func (self *Server) DeleteIndex(index string) bool {

	return self.Delete(index, func(resp *http.Response) error {
		if resp.StatusCode != 200 {
			return fmt.Errorf("Unexpected status code %d", resp.StatusCode)
		} else {
			return nil
		}
	}) == nil
}

type Document struct {
	Type    string
	Id      string
	Version float64
	Exists  bool
	Source  interface{}
}

func (self *Document) UnmarshalJSON(data []byte) error {
	var objmap map[string]*json.RawMessage

	if err := json.Unmarshal(data, &objmap); err != nil {
		return err
	}

	if _, ok := objmap["_type"]; ok {
		if err := json.Unmarshal(*objmap["_type"], &self.Type); err != nil {
			return err
		}
	}

	if _, ok := objmap["_id"]; ok {
		if err := json.Unmarshal(*objmap["_id"], &self.Id); err != nil {
			return err
		}
	}

	if _, ok := objmap["_version"]; ok {
		if err := json.Unmarshal(*objmap["_version"], &self.Version); err != nil {
			return err
		}
	}

	if _, ok := objmap["exists"]; ok {
		if err := json.Unmarshal(*objmap["exists"], &self.Exists); err != nil {
			return err
		}
	}

	if _, ok := objmap["_source"]; ok {
		if err := json.Unmarshal(*objmap["_source"], &self.Source); err != nil {
			return err
		}
	} else if _, ok := objmap["fields"]; ok {
		if err := json.Unmarshal(*objmap["fields"], &self.Source); err != nil {
			return err
		}
	}
	return nil
}

func (self *Server) PutDocument(index string, doctype string, id string, doc io.Reader) bool {
	cmd := fmt.Sprintf("%s/%s/%s", index, doctype, id)

	err := self.Put(cmd, doc, func(resp *http.Response) error {
		if resp.StatusCode != 200 && resp.StatusCode != 201 {
			return fmt.Errorf("Unexpected status code %d", resp.StatusCode)
		} else {
			return nil
		}
	})

	return err == nil
}

func (self *Server) GetDocument(index string, doctype string, id string) *Document {
	return self.GetDocumentFields(index, doctype, id, "")
}

func (self *Server) GetDocumentFields(index string, doctype string, id string, fields string) *Document {
	var cmd string
	var doc *Document

	if fields == "" {
		cmd = fmt.Sprintf("%s/%s/%s", index, doctype, id)
	} else {
		cmd = fmt.Sprintf("%s/%s/%s?fields=%s", index, doctype, id, fields)
	}

	self.Get(cmd, func(resp *http.Response) error {
		switch resp.StatusCode {
		case 200:
			doc = &Document{}
			return json.NewDecoder(resp.Body).Decode(doc)
		case 404:
			doc = &Document{Type: doctype, Id: id, Exists: false}
			return nil
		default:
			return fmt.Errorf("Unexpected status code %d", resp.StatusCode)
		}
	})

	// fmt.Printf("Returning doc %s\n", *doc)
	return doc
}

func (self *Server) Search() *Search {
	s := &Search{Server: self}
	return s
}
