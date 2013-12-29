package gosearch

import (
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

	b := server.HasIndex("index_does_not_exist")
	c.Assert(b, Equals, false)
}

func (s *ServerSuite) TestCreateIndex(c *C) {
	server := ConnectURL("http://localhost:9200")

	b := server.CreateIndex("test_index")
	c.Assert(b, Equals, true)
	c.Assert(server.HasIndex("test_index"), Equals, true)
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

	b := server.CreateIndexWithSettings("test_mapping_index", reader)
	c.Assert(b, Equals, true)
	c.Assert(server.HasIndex("test_mapping_index"), Equals, true)
	// server.DeleteIndex("test_index")

}
