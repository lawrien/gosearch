package gosearch

import (
	. "launchpad.net/gocheck"
	"testing"
)

// Start of setup
type serverSuite struct{}

func TestServer(t *testing.T) {
	Suite(&serverSuite{})
	TestingT(t)
}

func (s *serverSuite) TestConnect(c *C) {
	server := ConnectURL("http://localhost:9200")

	status, err := server.Status()
	c.Assert(err, IsNil)
	c.Assert(status.Status, Equals, 200)
}

func (s *serverSuite) BenchmarkConnect(c *C) {
	server := ConnectURL("http://localhost:9200")

	c.ResetTimer()
	for i := 0; i < c.N; i++ {
		status, err := server.Status()
		c.Assert(err, IsNil)
		c.Assert(status.Status, Equals, 200)
	}
}
