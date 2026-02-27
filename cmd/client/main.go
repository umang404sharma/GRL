package main

import "umang404sharma/GRL/internal/client"

func main() {
	hostID := "client-1"
	aggregator := "http://localhost:9000"

	server := client.NewServer(hostID, aggregator)
	server.Start("8001")
}
