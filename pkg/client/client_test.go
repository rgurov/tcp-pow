package client_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log/slog"
	"net"
	"testing"

	"github.com/rgurov/tcp-pow/pkg/client"
	"github.com/rgurov/tcp-pow/pkg/puzzle"
	"github.com/rgurov/tcp-pow/pkg/server"
)

const (
	host       = "localhost"
	port       = "8356"
	complexity = 1
)

func TestClientConnect(t *testing.T) {
	listener, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		t.Fatalf("cant listen on %s:%s: %s", host, port, err.Error())
	}
	defer listener.Close()

	client := client.NewClient(slog.Default(), host, port, complexity)

	serverErr := make(chan error)
	go func() {
		_, err := listener.Accept()
		serverErr <- err
	}()

	ctx := context.Background()
	clientConn, err := client.Connect(ctx)
	if err != nil {
		t.Fatalf("can't connect to server: %s", err.Error())
	}
	defer clientConn.Close()

	err = <-serverErr
	if err != nil {
		t.Fatalf("can't accept client connection: %s", err.Error())
	}
}

func TestClientSolve(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	client := client.NewClient(slog.Default(), host, port, complexity)

	serverPuzzle := puzzle.NewRandomPuzzle(complexity)
	serverPuzzleHash := serverPuzzle.GetInitialHash()
	solution := make([]byte, puzzle.SolutionSize)

	serverErr := make(chan error, 2)
	go func() {
		// sending puzzle from server
		_, err := serverConn.Write(serverPuzzleHash[:])
		serverErr <- err

		// reading solution
		_, err = serverConn.Read(solution)
		serverErr <- err
	}()

	err := client.Solve(clientConn)
	if err != nil {
		t.Fatalf("error while solving puzzle: %s", err.Error())
	}

	err = <-serverErr
	if err != nil {
		t.Fatalf("error while writing server conn: %s", err.Error())
	}

	err = <-serverErr
	if err != nil && err != io.EOF {
		t.Fatalf("error reading solution: %s", err.Error())
	}

	if !serverPuzzle.IsValidSolution([puzzle.SolutionSize]byte(solution)) {
		hash := sha256.New()
		hash.Write(solution)
		hash.Write(serverPuzzleHash[:])
		t.Fatalf("client sent invalid solution: %s", hex.EncodeToString(hash.Sum(nil)))
	}
}

func TestClientReceiveResponse(t *testing.T) {
	const message = "Hello, World!"

	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	client := client.NewClient(slog.Default(), host, port, complexity)

	serverErr := make(chan error, 2)
	go func() {
		messageSize := make([]byte, server.MessageLengthSize)
		messageSize[0] = byte(len(message))

		// sending message size to client
		_, err := serverConn.Write(messageSize)
		serverErr <- err

		// sending message to client
		_, err = serverConn.Write([]byte(message))
		serverErr <- err
	}()

	receivedMessage, err := client.ReceiveResponse(clientConn)
	if err != nil {
		t.Fatalf("error receiving reseponse: %s", err.Error())
	}

	err = <-serverErr
	if err != nil {
		t.Fatalf("error writing message size: %s", err.Error())
	}

	err = <-serverErr
	if err != nil {
		t.Fatalf("error writing message: %s", err.Error())
	}

	if receivedMessage != message {
		t.Fatalf("received message is not equal to original message, got: %s want: %s", receivedMessage, message)
	}
}
