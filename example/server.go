package main

import "github.com/paopaol/gocache"

func main() {
	// server := gocache.NewRestServer()
	server := gocache.NewTcpBnfServer()
	server.Run(":12346")
}
