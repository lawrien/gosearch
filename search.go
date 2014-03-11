package gosearch

import (
	// "bytes"
	"encoding/json"
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

type JSON map[string]interface{}

type Search struct {
	Server *Server
	Index  string
	Type   string
	Offset int
	Limit  int
	Fields string
	Query  JSON
	Filter JSON
}

func (self *Search) MarshalJSON() ([]byte, error) {
	var js = JSON{}
	if self.Offset != 0 {
		js["from"] = self.Offset
	}
	if self.Limit != 0 {
		js["size"] = self.Limit
	}
	if self.Fields != "" {
		js["fields"] = self.Fields
	}

	if len(self.Filter) > 0 {
		if len(self.Query) > 0 {
			js["filtered"] = JSON{"query": self.Query, "filter": self.Filter}
		} else {
			js["filtered"] = JSON{"filter": self.Filter}
		}
	} else if len(self.Query) > 0 {
		js["query"] = self.Query
	}

	return json.Marshal(js)
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
