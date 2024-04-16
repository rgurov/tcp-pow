package server

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net"
	"strings"

	"github.com/rgurov/tcp-pow/pkg/puzzle"
	"github.com/rgurov/tcp-pow/static"
)

const (
	MessageLengthSize = 4
)

type ConnectionState int

const (
	ConnectionStateInit ConnectionState = iota
	ConnectionStateWaitingSolution
)

type Server struct {
	logger       *slog.Logger
	address      string
	complexity   int
	responseList []string
}

func NewServer(
	logger *slog.Logger,
	host, port string,
	complexity int,
) *Server {
	address := host + ":" + port
	responseList := strings.Split(static.WordsOfWisdom, "\n")
	return &Server{
		logger:       logger,
		address:      address,
		complexity:   complexity,
		responseList: responseList,
	}
}

func (s *Server) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("can't listen on %s: %w", s.address, err)
	}

	s.logger.Info("server started on " + s.address)

	closeErr := make(chan error)
	go func() {
		<-ctx.Done()
		closeErr <- listener.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			return <-closeErr
		default:
			conn, err := listener.Accept()
			if err != nil {
				continue
			}

			go s.handleConnection(ctx, conn)
		}
	}
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	logger := s.logger.With(slog.String("ip", conn.RemoteAddr().String()))
	logger.Info("new connection")

	state := ConnectionStateInit
	clientPuzzle := puzzle.NewRandomPuzzle(s.complexity)

	defer func() {
		conn.Close()
		logger.Info("connection closed")
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			switch state {
			case ConnectionStateInit:
				initialHash := clientPuzzle.GetInitialHash()
				_, err := conn.Write(initialHash[:])

				if err != nil {
					logger.Error("error sending initial hash, closing connection: " + err.Error())
					return
				}
				logger.Info("puzzle sent, hash=" + hex.EncodeToString(initialHash[:]))
				state = ConnectionStateWaitingSolution
			case ConnectionStateWaitingSolution:
				solution := make([]byte, puzzle.SolutionSize)
				_, err := conn.Read(solution)
				if err != nil && err != io.EOF {
					logger.Error("error reading solution, closing connection: " + err.Error())
					return
				}

				if clientPuzzle.IsValidSolution([8]byte(solution)) {
					logger.Info("valid solution received")

					err := connWriteString(conn, s.chooseRandomResponse())
					if err != nil {
						logger.Error("error sending word of wisdom, closing connection: " + err.Error())
						return
					}
					logger.Info("word of wisdom sent")
				} else {
					logger.Info("invalid solution received")
				}

				return
			}
		}
	}
}

func (s *Server) chooseRandomResponse() string {
	return s.responseList[rand.Intn(len(s.responseList))]
}

func connWriteString(conn net.Conn, message string) error {
	bs := make([]byte, MessageLengthSize)
	binary.LittleEndian.PutUint32(bs, uint32(len(message)))

	_, err := conn.Write(bs)
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte(message))
	return err
}
