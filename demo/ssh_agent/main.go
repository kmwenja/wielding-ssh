package main

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

var USAGE = `Usage: %s <ssh server address> <ssh user>.

Ensure that your SSH agent is running before running this.

Example: %s 127.0.0.1:22 myuser
`

func main() {
	if len(os.Args) < 3 {
		fmt.Printf(USAGE, os.Args[0], os.Args[0])
		os.Exit(1)
	}
	addr := os.Args[1]
	user := os.Args[2]

	// START OMIT
	c, _ := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	defer c.Close()

	a := agent.NewClient(c)
	// get all the pubkeys from the agent
	signers, _ := a.Signers()

	// use the keys
	authMethod := ssh.PublicKeys(signers...)
	// END OMIT

	cfg := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, _ := ssh.Dial("tcp", addr, cfg)
	defer client.Close()

	session, _ := client.NewSession()
	defer session.Close()

	startTerminal(session)
}
