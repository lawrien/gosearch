package gosearch

import (
	"fmt"
	"io"
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
	cmd := fmt.Sprintf("%s", self.url)

	if resp, err := DefaultConnectionPool.Do(GET, cmd, nil); err != nil {
		return nil, err
	} else if resp.Status != 200 {
		return nil, fmt.Errorf("%d: Cannot get status", resp.Status)
	} else {
		s := &Status{}
		err = resp.Convert(s)
		return s, err
	}
}

func (self *Server) HasIndex(index string) error {
	cmd := fmt.Sprintf("%s/%s", self.url, index)

	if resp, err := DefaultConnectionPool.Do(HEAD, cmd, nil); err != nil {
		return err
	} else if resp.Status != 200 {
		return fmt.Errorf("%d: Index %s not found.", resp.Status, index)
	} else {
		return nil
	}
}

func (self *Server) CreateIndex(index string) error {
	return self.CreateIndexWithSettings(index, nil)
}

func (self *Server) CreateIndexWithSettings(index string, settings io.Reader) error {
	cmd := fmt.Sprintf("%s/%s", self.url, index)

	if resp, err := DefaultConnectionPool.Do(PUT, cmd, settings); err != nil {
		return err
	} else if resp.Status != 200 {
		return fmt.Errorf("%d: Unable to create index %s.", resp.Status, index)
	} else {
		return nil
	}
}

func (self *Server) DeleteIndex(index string) error {
	cmd := fmt.Sprintf("%s/%s", self.url, index)

	if resp, err := DefaultConnectionPool.Do(DELETE, cmd, nil); err != nil {
		return err
	} else if resp.Status != 200 {
		return fmt.Errorf("%d: Unable to delete index %s.", resp.Status, index)
	} else {
		return nil
	}
}

func (self *Server) PutDocument(index string, doctype string, id string, doc interface{}) error {
	cmd := fmt.Sprintf("%s/%s/%s/%s", self.url, index, doctype, id)

	// convert doc into json

	if resp, err := DefaultConnectionPool.Do(PUT, cmd, doc); err != nil {
		return err
	} else if resp.Status != 200 && resp.Status != 201 {
		return fmt.Errorf("%d: Unable to put document %s to index %s.", resp.Status, id, index)
	} else {
		return nil
	}
}

func (self *Server) GetDocument(index string, doctype string, id string) (*Document, error) {
	return self.GetDocumentFields(index, doctype, id, "")
}

func (self *Server) GetDocumentFields(index string, doctype string, id string, fields string) (*Document, error) {
	var cmd string

	if fields == "" {
		cmd = fmt.Sprintf("%s/%s/%s/%s", self.url, index, doctype, id)
	} else {
		cmd = fmt.Sprintf("%s/%s/%s/%s?fields=%s", self.url, index, doctype, id, fields)
	}

	if resp, err := DefaultConnectionPool.Do(GET, cmd, nil); err != nil {
		return nil, err
	} else if resp.Status == 200 {
		doc := &Document{}
		err = resp.Convert(doc)
		return doc, err
	} else if resp.Status == 404 {
		doc := &Document{Type: doctype, Id: id, Exists: false}
		return doc, nil
	} else {
		return nil, fmt.Errorf("%d: Unable to get document %s from index %s.", resp.Status, id, index)
	}
}

func (self *Server) NewSearch() *Search {
	s := &Search{Server: self}
	return s
}
