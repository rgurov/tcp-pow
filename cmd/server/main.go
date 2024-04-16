package main

import (
	"context"
	"flag"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/rgurov/tcp-pow/pkg/server"
)

func main() {
	host := flag.String("h", "localhost", "host")
	port := flag.String("p", "7771", "port")
	complexity := flag.Int("c", 4, "complexity (count of leading zeros in hash)")
	flag.Parse()

	logger := slog.Default()
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	srv := server.NewServer(
		logger,
		*host,
		*port,
		*complexity,
	)

	err := srv.Start(ctx)
	if err != nil {
		panic(err)
	}
}
