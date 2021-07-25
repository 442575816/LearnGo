package main

import (
	"fmt"
	"net"
)

type Server struct {
	IP          string
	Port        uint
	TCPListener *net.TCPListener
}

type ServerSocketChannel interface {
	Init(ip string, port uint) bool
	Start() (int, error)
}

func (server *Server) Init(ip string, port uint) bool {
	server.IP = ip
	server.Port = port

	return true
}

func (server *Server) Start() (int, error) {
	if server.IP == "" {
		server.IP = "127.0.0.1"
	}
	var addr = fmt.Sprintf("%s:%d", server.IP, server.Port)
	remoteAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return -1, err
	}

	ln, err := net.ListenTCP("tcp4", remoteAddr)
	if err != nil {
		return -1, err
	}

	server.TCPListener = ln
	fmt.Printf("Listen tcp succ on port %d\n", server.Port)
	return 0, nil
}

// func main() {
// 	var server Server
// 	server.Init("127.0.0.1", 4455)
// 	server.Start()
// }
