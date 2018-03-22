package main

import (
	"github.com/davidbelliott/nyx/config"
	"github.com/davidbelliott/nyx/http"
	"log"
	"os"
	"time"
	"flag"
)

func main() {
	flag.Parse()
	c, err := config.Load()
	if err != nil {
		log.Printf("Could not read configuration: %s\n", err)
		return
	}

	log.Println("Starting Server")
	if err := http.Start(c); err != nil {
		log.Printf("Could not start server or server crashed: %s\n", err)
		log.Printf("Waiting 10 seconds before dying...")
		time.Sleep(10 * time.Second)
		log.Printf("Exiting")
		os.Exit(1)
		return
	}
	os.Exit(0)
}
