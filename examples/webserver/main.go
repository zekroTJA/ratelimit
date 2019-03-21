package main

import "log"

func main() {
	ws := newWebServer(":8080")

	log.Println("Starting web server...")
	log.Fatal(ws.start())
}
