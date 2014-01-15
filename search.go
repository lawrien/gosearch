package gosearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
	query  *Query
}

type Query struct {
	Must    []map[string]interface{}
	MustNot []map[string]interface{}
}

func (self *Search) Query() *Query {
	if self.query == nil {
		self.query = &Query{}
	}
	return self.query
}

func (self *Search) Run() *SearchResults {
	var cmd string

	if self.Type != "" {
		cmd = fmt.Sprintf("%s/%s/_search", self.Index, self.Type)
	} else {
		cmd = fmt.Sprintf("%s/_search", self.Index)
	}

	var results *SearchResults

	if searchJson, err := json.Marshal(self); err != nil {
		// FIXME log error
	} else {
		self.Server.Post(cmd, bytes.NewReader(searchJson), func(resp *http.Response) error {
			if resp.StatusCode == 200 {
				results = &SearchResults{}
				return json.NewDecoder(resp.Body).Decode(results)
			} else {
				return fmt.Errorf("Unexpected status code %d", resp.StatusCode)
			}
		})
	}
	return results

}
