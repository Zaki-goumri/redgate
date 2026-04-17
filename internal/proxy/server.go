package proxy

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/zaki/redgate/internal/resp"
)

type Server struct {
	addr     string
	upstream string
}

func NewServer(addr string, upstream string) *Server {
	return &Server{addr, upstream}
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
	upstream, err := net.Dial("tcp", s.upstream)
	if err != nil {
		log.Printf("Failed to connect to upstream Redis: %v\n", err)
		return
	}
	defer upstream.Close()

	reader := resp.NewReader(conn)
	upstreamReader := resp.NewReader(upstream)

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

		default:
			_, err := upstream.Write(val.Marshal())
			if err != nil {
				log.Printf("Upstream write error: %v\n", err)
				return
			}

			response, err := upstreamReader.Read()
			if err != nil {
				log.Printf("Upstream read error: %v\n", err)
				return
			}

			conn.Write(response.Marshal())
		}
	}
}
