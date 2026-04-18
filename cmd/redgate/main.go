package main

import (
	"log"

	"github.com/zaki/redgate/internal/proxy"
)

func main() {
	srv, err := proxy.NewServer(":6382", []string{
		"localhost:6379",
		"localhost:6381",
	})
	if err = srv.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
