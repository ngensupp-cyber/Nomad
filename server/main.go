package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"nomad-c2/server/db"
	"nomad-c2/server/api"
	"nomad-c2/server/c2"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env if exists
	godotenv.Load()

	webPortEnv := os.Getenv("WEB_PORT")
	if webPortEnv == "" {
		webPortEnv = "8080"
	}
	c2PortEnv := os.Getenv("C2_PORT")
	if c2PortEnv == "" {
		c2PortEnv = "5555"
	}

	port := flag.String("port", webPortEnv, "Web UI port")
	c2port := flag.String("c2port", c2PortEnv, "C2 Listener port")
	flag.Parse()

	fmt.Println(`
    _   __                       __   ______ ___ 
   / | / /____   ____ ___   ____ _/ /  / ____|__ \
  /  |/ // __ \ / __ '__ \ / __ '/ /  / /    __/ /
 / /|  // /_/ // / / / / // /_/ / /  / /___ / __/ 
/_/ |_/ \____//_/ /_/ /_/ \__,_/_/   \____/|____/ 
                                                  
    Command the desert, from anywhere.
    `)

	// Initialize Databases
	db.InitDB()
	db.InitCache()

	// Start C2 Listener
	go c2.StartListener(*c2port)

	// Start Web Server
	router := api.SetupRoutes()
	
	fmt.Printf("[+] Web UI started on http://0.0.0.0:%s\n", *port)
	fmt.Printf("[+] C2 Listener started on port %s\n", *c2port)
	
	log.Fatal(http.ListenAndServe(":"+*port, router))
}
