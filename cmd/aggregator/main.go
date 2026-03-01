package main

import "umang404sharma/GRL/internal/aggregator"

func main() {
	zone := "zone-a"
	controller := "http://localhost:10000"

	server := aggregator.NewServer(zone, controller)

	server.Start("9000")
}
