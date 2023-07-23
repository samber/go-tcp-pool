package main

import (
	"log"
	"time"

	pool "github.com/samber/go-tcp-pool"
)

func main() {
	// ncat -l 9999 -k

	conn, err := pool.Dial("tcp", "localhost:9999")
	if err != nil {
		log.Fatal(err)
	}

	_ = conn.SetPoolSize(10)
	conn.SetMaxRetries(10)
	conn.SetRetryInterval(10 * time.Millisecond)

	_, _ = conn.Write([]byte("Hello, world!\n"))
	_, _ = conn.Write([]byte("Hello, world!\n"))
	_, _ = conn.Write([]byte("Hello, world!\n"))
	_, _ = conn.Write([]byte("Hello, world!\n"))
	_, _ = conn.Write([]byte("Hello, world!\n"))
}
