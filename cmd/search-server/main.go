package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	search "github.com/localitas/localitas-app-search"
	"github.com/localitas/localitas-go"
	"github.com/urfave/cli/v3"
)

var (
	version = "dev"
	commit  = "unknown"
)

func envOrFileToken() string {
	if t := os.Getenv("LOCALITAS_API_TOKEN"); t != "" {
		return t
	}
	return client.DefaultToken()
}

func main() {
	app := &cli.Command{
		Name:    "search-server",
		Usage:   "search app server",
		Version: version + " (" + commit + ")",
		Commands: []*cli.Command{
			serveCommand(),
		},
		DefaultCommand: "serve",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return serveAction(ctx, cmd)
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func serveCommand() *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "Start the server",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "listen", Value: "127.0.0.1:0", Usage: "listen address (core reaches the app over loopback; a container passes :8000)"},
			&cli.StringFlag{Name: "core-url", Value: client.DefaultCoreURL(), Usage: "base URL of the Localitas core API"},
			&cli.StringFlag{Name: "base-path", Value: "/", Usage: "URL prefix for <base href>"},
			&cli.StringFlag{Name: "token", Value: envOrFileToken(), Usage: "bearer token for API calls"},
		},
		Action: serveAction,
	}
}

func newClient(cmd *cli.Command) *client.Client {
	c := client.New(cmd.String("core-url"))
	if t := cmd.String("token"); t != "" {
		c = c.WithToken(t)
	}
	return c
}

func serveAction(ctx context.Context, cmd *cli.Command) error {
	basePath := cmd.String("base-path")
	c := newClient(cmd)

	a := search.New(basePath)
	mux := http.NewServeMux()
	a.RegisterRoutes(mux)
	mux.HandleFunc("GET /health.json", search.HandleHealth)

	ln, err := net.Listen("tcp", cmd.String("listen"))
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	addr := ln.Addr().(*net.TCPAddr)
	fmt.Printf("search-server listening on http://localhost:%d\n", addr.Port)

	announcement := c.Announce(ln)

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		if err := announcement.Withdraw(); err != nil {
			slog.Warn("withdraw from core failed", "error", err)
		}
		os.Exit(0)
	}()

	return http.Serve(ln, mux)
}
