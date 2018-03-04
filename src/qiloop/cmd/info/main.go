package main

import (
	"encoding/json"
	"fmt"
	"log"
	service0 "qiloop/meta/stage2"
	directory "qiloop/meta/stage3"
	"qiloop/net"
	"qiloop/value"
)

func main() {
	endpoint := ":9559"
	conn, err := net.NewClient(endpoint)
	if err != nil {
		log.Fatalf("failed to connect %s: %s", endpoint, err)
	}
	server := service0.Server{net.NewProxy(conn, 0, 0)}
	permissions := map[string]value.Value{
		"ClientServerSocket":    value.Bool(true),
		"MessageFlags":          value.Bool(true),
		"MetaObjectCache":       value.Bool(true),
		"RemoteCancelableCalls": value.Bool(true),
	}
	permissions, err = server.Authenticate(permissions)
	if err != nil {
		log.Fatalf("authentication failed: %s", err)
	}

	directory := directory.Directory{net.NewProxy(conn, 1, 1)}
	services, err := directory.Services()
	if err != nil {
		log.Fatalf("failed to list services: %s", err)
	}

	json, err := json.MarshalIndent(services, "", "    ")
	if err != nil {
		log.Fatalf("json encoding failed: %s", err)
	}
	fmt.Println(string(json))
}
