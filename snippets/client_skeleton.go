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

	// connect to SSH port and perform handshake
	client, _ := ssh.Dial("tcp", "127.0.0.1:22", cfg)
	defer client.Close()

	// connection established!

	// client.HandleChannelOpen() to handle incoming channels from the server
	// client.Listen() for tunnels
	// session, _ := client.NewSession()
	// session can start shells and  run commands
}
