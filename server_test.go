package gosearch

import (
	"encoding/json"
	"fmt"
	. "launchpad.net/gocheck"
	"strings"
	"testing"
)

// Start of setup
type ServerSuite struct{}

func TestServer(t *testing.T) {
	Suite(&ServerSuite{})
	TestingT(t)
}

func (s *ServerSuite) TestConnect(c *C) {
	server := ConnectURL("http://localhost:9200")

	status, err := server.Status()
	c.Assert(err, IsNil)
	c.Assert(status.Status, Equals, 200)
	fmt.Printf("Status => %s\n", status)
}

func (s *ServerSuite) BenchmarkConnect(c *C) {
	server := ConnectURL("http://localhost:9200")

	c.ResetTimer()
	for i := 0; i < c.N; i++ {
		status, err := server.Status()
		c.Assert(err, IsNil)
		c.Assert(status.Status, Equals, 200)
	}
}

func (s *ServerSuite) TestHasIndex(c *C) {
	server := ConnectURL("http://localhost:9200")

	err := server.HasIndex("index_does_not_exist")
	c.Assert(err, NotNil)
}

func (s *ServerSuite) TestCreateIndex(c *C) {
	server := ConnectURL("http://localhost:9200")

	err := server.CreateIndex("test_index")
	c.Assert(err, IsNil)
	c.Assert(server.HasIndex("test_index"), IsNil)
	server.DeleteIndex("test_index")
}

func (s *ServerSuite) TestCreateIndexWithMapping(c *C) {
	server := ConnectURL("http://localhost:9200")

	settings := `{
      "mappings" : {
        "type1" : {
          "properties" : {
            "field1" : { "type" : "string", "index" : "not_analyzed" }
          }
        }
      }
    }`

	reader := strings.NewReader(settings)

	err := server.CreateIndexWithSettings("test_mapping_index", reader)
	c.Assert(err, IsNil)
	c.Assert(server.HasIndex("test_mapping_index"), IsNil)
	server.DeleteIndex("test_mapping_index")

}

func (s *ServerSuite) TestCreateDocument(c *C) {
	doc := `{ "name":"George", "age":25 }`
	server := ConnectURL("http://localhost:9200")

	server.CreateIndex("test_index")
	c.Assert(server.HasIndex("test_index"), IsNil)

	err := server.PutDocument("test_index", "person", "1", doc)
	c.Assert(err, IsNil)

	server.DeleteIndex("test_index")

}

func (s *ServerSuite) TestGetDocumentSource(c *C) {

	doc := `{ "name":"George", "age":25 }`
	server := ConnectURL("http://localhost:9200")

	server.CreateIndex("test_index")
	c.Assert(server.HasIndex("test_index"), IsNil)

	fmt.Printf("Putting document\n")
	err := server.PutDocument("test_index", "person", "1", doc)
	c.Assert(err, IsNil)

	if doc, err := server.GetDocument("test_index", "person", "1"); err != nil {
		c.Errorf("Error", err)
	} else {
		c.Assert(doc, NotNil)
		c.Assert(doc.Exists, Equals, true)
		bytes, _ := json.Marshal(doc)
		fmt.Printf("Document => %s\n", string(bytes))
	}
	server.DeleteIndex("test_index")
}

func (s *ServerSuite) BenchmarkGetDocumentSource(c *C) {
	doc := `{ "name":"George", "age":25 }`
	server := ConnectURL("http://localhost:9200")

	server.CreateIndex("benchmark_index")
	server.PutDocument("benchmark_index", "person", "1", doc)

	c.ResetTimer()
	for i := 0; i < c.N; i++ {
		if doc, err := server.GetDocument("benchmark_index", "person", "1"); err != nil {
			c.Errorf("Unable to load document")
		} else {
			c.Assert(doc, NotNil)
			c.Assert(doc.Source["name"], Equals, "George")
		}
	}
	c.StopTimer()
	server.DeleteIndex("benchmark_index")
}

func (s *ServerSuite) TestGetDocumentFields(c *C) {
	doc := `{ "name":"George", "age":25 }`
	server := ConnectURL("http://localhost:9200")

	server.CreateIndex("test_index")
	c.Assert(server.HasIndex("test_index"), IsNil)

	c.Assert(server.PutDocument("test_index", "person", "1", doc), IsNil)

	if doc, err := server.GetDocumentFields("test_index", "person", "1", "name"); err != nil {
		c.Errorf("Unable to load document")
	} else {
		c.Assert(doc, NotNil)
		c.Assert(doc.Exists, Equals, true)
		bytes, _ := json.Marshal(doc)
		fmt.Printf("Document => %s\n", string(bytes))
	}
	server.DeleteIndex("test_index")
}

func (s *ServerSuite) TestSearch(c *C) {
	doc := `{ "name":"George", "age":25 }`
	server := ConnectURL("http://localhost:9200")

	if !server.HasIndex("test_search") {
		server.CreateIndex("test_search")
		for i := 1; i < 20; i++ {
			server.PutDocument("test_search", "person", fmt.Sprintf("%d", i), doc)
		}
	}

	search := server.NewSearch()
	search.Index = "test_search"
	search.Limit = 5
	search.Query = JSON{"match": JSON{"name": "george"}}

	if r, err := search.Run(); err != nil {
		c.Errorf("Search failed")
	} else {
		bytes, _ := json.Marshal(r)
		fmt.Printf("Search => %s\n", string(bytes))
	}

	// server.DeleteIndex("test_search")
}
