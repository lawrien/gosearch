package gosearch

import (
	"encoding/json"
)

type Document struct {
	Index   string                 `json:"_index"`
	Type    string                 `json:"_type"`
	Id      string                 `json:"_id"`
	Version float64                `json:"_version,omitempty"`
	Exists  bool                   `json:"exists,omitempty"`
	Source  map[string]interface{} `json:"_source,omitempty"`
}

func (self *Document) UnmarshalJSON(data []byte) error {

	var objmap map[string]interface{}

	if err := json.Unmarshal(data, &objmap); err != nil {
		return err
	}

	if _, ok := objmap["fields"]; ok {
		objmap["_source"] = objmap["fields"]
		delete(objmap, "fields")
	}

	if prop, ok := objmap["_type"]; ok {
		self.Type = prop.(string)
	}

	if prop, ok := objmap["_id"]; ok {
		self.Id = prop.(string)
	}
	if prop, ok := objmap["_version"]; ok {
		self.Version = prop.(float64)
	}
	if prop, ok := objmap["exists"]; ok {
		self.Exists = prop.(bool)
	}

	if prop, ok := objmap["_source"]; ok {
		self.Source = prop.(map[string]interface{})
	}

	return nil
}
