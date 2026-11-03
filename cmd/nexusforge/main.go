package main

import (
 "log"
 "os"
 "github.com/niknight1403/Nexusforge/internal/server"
)

func main() {
 addr := os.Getenv("NEXUSFORGE_ADDR")
 if addr == "" { addr = "127.0.0.1:8080" }
 log.Printf("NexusForge listening at %s", addr)
 if err := server.Listen(addr); err != nil { log.Fatal(err) }
}
