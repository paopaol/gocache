package main

import "github.com/gocache"


func main() {
	server := gocache.NewRestServer()
	server.Run(":9997")
}
