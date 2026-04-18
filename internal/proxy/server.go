package proxy

import (
	"fmt"
	"hash/fnv"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/zaki/redgate/internal/resp"
)

type Server struct {
	addr      string
	upstreams []string
	conns     map[string]net.Conn
	mu        sync.Mutex
}

func (s *Server) getConn(addr string) (net.Conn, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if conn, ok := s.conns[addr]; ok {
		return conn, nil
	}
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Errorf("check the error %v", err)
		return nil, err
	}
	s.conns[addr] = conn
	return conn, nil
}

func NewServer(addr string, upstreams []string) (*Server, error) {
	s := &Server{addr: addr, upstreams: upstreams, conns: make(map[string]net.Conn)}
	for _, upstream := range upstreams {
		conn, err := net.Dial("tcp", upstream)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to %s %v", upstream, err)
		}
		s.conns[upstream] = conn
		fmt.Println("connected to %v", upstream)
	}
	return s, nil
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

func (s *Server) getUpstream(key string) string {
	h := fnv.New32()
	h.Write([]byte(key))
	idx := h.Sum32() % uint32(len(s.upstreams))
	return s.upstreams[idx]
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
			if len(val.Array) < 2 {
				fmt.Fprint(conn, "-ERR no key provided\r\n")
				continue
			}
			key := string(val.Array[1].Str)

			upstreamAddr := s.getUpstream(key)
			upstream, err := s.getConn(upstreamAddr)
			if err != nil {
				log.Printf("Upstream connection error: %v\n", err)
				return
			}

			_, err = upstream.Write(val.Marshal())
			if err != nil {
				log.Printf("Upstream write error: %v\n", err)
				return
			}

			upstreamReader := resp.NewReader(upstream)
			response, err := upstreamReader.Read()
			if err != nil {
				log.Printf("Upstream read error: %v\n", err)
				return
			}
			conn.Write(response.Marshal())
		}
	}
}
