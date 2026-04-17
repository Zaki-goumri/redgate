package main

import (
	"log"

	"github.com/zaki/redgate/internal/proxy"
)

func main() {
	srv := proxy.NewServer(":6382", "localhost:6379")
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
