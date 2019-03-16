package main

import "fmt"

func main() {
	cfg := &ssh.ClientConfig{
		User: "",
		Auth: []ssh.AuthMethod{
			// add authentication methods here
		},
		HostKeyCallback: // authenticate server,
	}

	// connect to SSH port
	c, _ := net.Dial("tcp", "127.0.0.1:22")
	defer c.Close()

	// perform handshake
	conn, channels, globalReqs, _ := ssh.NewClient(c, "127.0.0.1:22", cfg)
	defer conn.Close()

	// connection established!
}
