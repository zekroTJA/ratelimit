package main

import "log"

func main() {
	ws := NewWebServer(":8080")

	log.Println("Starting web server...")
	log.Fatal(ws.Start())
}
