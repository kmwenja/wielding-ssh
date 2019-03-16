package main

import (
	"net"

	"golang.org/x/crypto/ssh"
)

func main() {
	cfg := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			return nil, nil // accepts login
		},
		PublicKeyCallback: func(c ssh.ConnMetadata, pubkey ssh.PublicKey) (*ssh.Permissions, error) {
			return nil, nil // accepts login
		},
	}
	cfg.AddHostKey(key)

	// listen on SSH port
	listener, _ := net.Listen("tcp", "0.0.0.0:22")
	defer listener.Close()

	// receive TCP connection
	c, _ := listener.Accept()

	// perform handshake
	conn, channels, globalreqs, _ := ssh.NewServerConn(c, cfg)

	// connection established!
}
