package e2e_test

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/rgurov/tcp-pow/pkg/client"
	"github.com/rgurov/tcp-pow/pkg/server"
)

const (
	host       = "localhost"
	port       = "8356"
	complexity = 1
)

func TestE2E(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := server.NewServer(slog.Default(), host, port, complexity)
	serverErr := make(chan error)
	go func() {
		err := srv.Start(ctx)
		serverErr <- err
	}()
	clt := client.NewClient(slog.Default(), host, port, complexity)

	conn, err := clt.Connect(ctx)
	if err != nil {
		t.Fatalf("error connecting to server: %s", err.Error())
	}

	err = clt.Solve(conn)
	if err != nil {
		t.Fatalf("error solving puzzle: %s", err.Error())
	}

	msg, err := clt.ReceiveResponse(conn)
	if err != nil {
		t.Fatalf("error receiving response: %s", err.Error())
	}

	cancel()
	err = <-serverErr
	if err != nil {
		t.Fatalf("error from server: %s", err.Error())
	}

	fmt.Println("Done! Message from server: " + msg)
}
