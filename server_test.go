package gosearch

import (
	"fmt"
	"testing"
)

func TestConnect(t *testing.T) {
	server := ConnectURL("http://localhost:9200")

	if status, err := server.Status(); err != nil {
		fmt.Printf("Failed with %s\n", err)
	} else {
		fmt.Printf("Succeeded %s\n", status)
	}
}

func BenchmarkConnect(b *testing.B) {
	server := ConnectURL("http://localhost:9200")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if status, err := server.Status(); err != nil || status.Status != 200 {
			fmt.Printf("Failed with %s\n", err)
		}
	}
}
