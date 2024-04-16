package main

import (
	"context"
	"flag"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/rgurov/tcp-pow/pkg/client"
)

func main() {
	host := flag.String("h", "localhost", "host")
	port := flag.String("p", "7771", "port")
	complexity := flag.Int("c", 4, "complexity (count of leading zeros in hash)")
	flag.Parse()

	logger := slog.Default()
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	c := client.NewClient(
		logger,
		*host,
		*port,
		*complexity,
	)

	conn, err := c.Connect(ctx)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	err = c.Solve(conn)
	if err != nil {
		panic(err)
	}

	message, err := c.ReceiveResponse(conn)
	if err != nil {
		panic(err)
	}

	logger.Info("got message from server: " + message)
}
