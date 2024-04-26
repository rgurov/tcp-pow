package client

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"net"
	"time"

	"github.com/rgurov/tcp-pow/pkg/puzzle"
	"github.com/rgurov/tcp-pow/pkg/server"
)

const (
	timeout = time.Minute
)

type Client struct {
	logger     *slog.Logger
	address    string
	complexity int
}

func NewClient(
	logger *slog.Logger,
	host, port string,
	complexity int,
) *Client {
	return &Client{
		logger:     logger,
		address:    host + ":" + port,
		complexity: complexity,
	}
}

func (c *Client) Connect(ctx context.Context) (net.Conn, error) {
	dialer := net.Dialer{
		Timeout: timeout,
	}

	conn, err := dialer.DialContext(ctx, "tcp", c.address)
	if err != nil {
		return nil, fmt.Errorf("error connecting to %s: %w", c.address, err)
	}
	c.logger.Info("connected")

	return conn, nil
}

func (c *Client) Solve(conn net.Conn) error {
	hash := make([]byte, puzzle.PuzzleSize)
	_, err := conn.Read(hash)
	if err != nil && err != io.EOF {
		return fmt.Errorf("error receiving puzzle hash: %w", err)
	}

	c.logger.Info("puzzle hash received")

	clientPuzzle := puzzle.NewPuzzle([puzzle.PuzzleSize]byte(hash), c.complexity)
	solver := puzzle.NewPuzzleSolver(clientPuzzle)

	c.logger.Info("solving the puzzle...")
	for !solver.Solve() {
	}

	solution, err := solver.GetSolution()
	if err != nil {
		return fmt.Errorf("error getting solution: %w", err)
	}

	solvedHash := sha256.New()
	solvedHash.Write(hash)
	solvedHash.Write(solution[:])
	solvedSum := solvedHash.Sum(nil)

	c.logger.Info("puzzle solved, hash=" + hex.EncodeToString(solvedSum))

	_, err = conn.Write(solution[:])
	if err != nil {
		return fmt.Errorf("error sending solution: %w", err)
	}

	return nil
}

func (c *Client) ReceiveResponse(conn net.Conn) (string, error) {
	messageLength := make([]byte, server.MessageLengthSize)
	_, err := conn.Read(messageLength)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("error receiving message length: %w", err)
	}

	length := binary.LittleEndian.Uint32(messageLength)
	message := make([]byte, length)
	fmt.Println(length)

	_, err = conn.Read(message)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("error receiving message: %w", err)
	}

	return string(message), nil
}
