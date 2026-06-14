package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/localitas/localitas-app-search"
)

var (
	version = "dev"
	commit  = "unknown"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "version") {
		fmt.Printf("search-server %s (commit: %s)\n", version, commit)
		os.Exit(0)
	}

	var (
		listen   = flag.String("listen", ":0", "listen address")
		basePath = flag.String("base-path", "/", "URL prefix for <base href>")
	)
	flag.Parse()

	app := search.New(*basePath)

	mux := http.NewServeMux()
	app.RegisterRoutes(mux)
	mux.HandleFunc("GET /health.json", search.HandleHealth)

	ln, err := net.Listen("tcp", *listen)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	addr := ln.Addr().(*net.TCPAddr)
	fmt.Printf("search-server listening on http://localhost:%d\n", addr.Port)

	shutdown, err := search.BroadcastMDNS(addr.Port, search.DefaultHealth.Name)
	if err != nil {
		log.Printf("⚠️  mDNS broadcast failed: %v", err)
	}

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Println("shutting down...")
		if shutdown != nil {
			shutdown()
		}
		os.Exit(0)
	}()

	if err := http.Serve(ln, mux); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
