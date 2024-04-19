package main

import (
	"context"
	"flag"
	"log/slog"
	"net"
	"os/signal"
	"syscall"
	"time"

	"github.com/rgurov/tcp-pow/pkg/client"
)

func main() {
	host := flag.String("h", "localhost", "host")
	port := flag.String("p", "7771", "port")
	complexity := flag.Int("c", 6, "complexity (count of leading zeros in hash)")
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

	var conn net.Conn
	err := withRetries(
		func() error {
			var err error
			conn, err = c.Connect(ctx)
			return err
		},
		3,
		time.Second*3,
	)

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
	logger.Info("done! waiting for signal to close")
	<-ctx.Done()
}

func withRetries(f func() error, retries int, delay time.Duration) error {
	var err error
	for i := 0; i < retries; i++ {
		err = f()
		if err == nil {
			return nil
		}
	}
	return err
}
