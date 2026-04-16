package proxy

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/zaki/redgate/internal/resp"
)

type Server struct {
	addr string
}

func NewServer(addr string) *Server {
	return &Server{addr}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	log.Printf("Redgate frontend listening on %s\n", s.addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v\n", err)
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("Client connected: %s\n", conn.RemoteAddr())
	reader := resp.NewReader(conn)
	for {
		val, err := reader.Read()
		if err != nil {
			log.Printf("Client disconnected: %v\n", err)
			return
		}

		if val.Typ != resp.Array || len(val.Array) == 0 {
			continue
		}

		var args []string
		for _, v := range val.Array {
			args = append(args, string(v.Str))
		}
		fmt.Printf("Received: %s\n", strings.Join(args, " "))

		cmd := strings.ToUpper(string(val.Array[0].Str))

		switch cmd {
		case "PING":
			fmt.Fprint(conn, "+PONG\r\n")
		case "ECHO":
			if len(val.Array) < 2 {
				fmt.Fprint(conn, "-ERR wrong number of arguments for 'echo' command\r\n")
				continue
			}
			arg := string(val.Array[1].Str)
			fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(arg), arg)
		case "COMMAND":
			if len(val.Array) >= 2 {
				sub := strings.ToUpper(string(val.Array[1].Str))
				if sub == "COUNT" {
					fmt.Fprint(conn, ":0\r\n")
				} else {
					fmt.Fprint(conn, "*0\r\n")
				}
			} else {
				fmt.Fprint(conn, "*0\r\n")
			}
		case "SET":
			fmt.Fprint(conn, "+OK\r\n")
		case "GET":
			fmt.Fprint(conn, "$-1\r\n")
		default:
			fmt.Fprintf(conn, "-ERR unknown command '%s'\r\n", cmd)
		}
	}
}
