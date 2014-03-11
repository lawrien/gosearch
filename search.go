package gosearch

import (
	// "bytes"
	// "encoding/json"
	"fmt"
	// "net/http"
)

type SearchResults struct {
	Time     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Hits     struct {
		Total    int        `json:"total"`
		MaxScore *float64   `json:"max_score"`
		Results  []Document `json:"hits"`
	} `json:"hits"`
}

type Search struct {
	Server *Server `json:"-"`
	Index  string  `json:"-"`
	Type   string  `json:"-"`
	Offset int     `json:"from,omitempty"`
	Limit  int     `json:"size,omitempty"`
	Fields string  `json:"fields,omitempty"`
	query  Query
	filter Filter
}

type Query struct {
	Must    []map[string]interface{} `json:",omitempty"`
	MustNot []map[string]interface{} `json:",omitempty"`
	Should  []map[string]interface{} `json:",omitempty"`
}

func (self *Search) Query() *Query {
	return &self.query
}

type Filter struct {
	Must    []map[string]interface{} `json:",omitempty"`
	MustNot []map[string]interface{} `json:",omitempty"`
	Should  []map[string]interface{} `json:",omitempty"`
}

func (self *Search) Filter() *Filter {
	return &self.filter
}

func (self *Search) Run() (*SearchResults, error) {
	var cmd string

	if self.Type != "" {
		cmd = fmt.Sprintf("%s/%s/%s/_search", self.Server.url, self.Index, self.Type)
	} else {
		cmd = fmt.Sprintf("%s/%s/_search", self.Server.url, self.Index)
	}

	if resp, err := DefaultConnectionPool.Do(GET, cmd, self); err != nil {
		return nil, err
	} else if resp.Status == 200 {
		results := &SearchResults{}
		err = resp.Convert(results)
		return results, err
	} else {
		return nil, fmt.Errorf("%d: Unable to get results from indexes %s.", resp.Status, self.Index)
	}

}
