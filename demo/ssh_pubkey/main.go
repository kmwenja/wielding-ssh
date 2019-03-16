package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh"
)

var USAGE = `Usage: %s <ssh server address> <ssh user> <ssh key path>

Example: %s 127.0.0.1:22 myuser ~/.ssh/id_rsa
`

func main() {
	if len(os.Args) < 4 {
		fmt.Printf(USAGE, os.Args[0], os.Args[0])
		os.Exit(1)
	}
	addr, user, keyPath := os.Args[1], os.Args[2], os.Args[3]

	keyBytes, _ := ioutil.ReadFile(keyPath)

	// you can parse passphrase protected keys too
	key, _ := ssh.ParsePrivateKey(keyBytes)

	cfg := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key), // other keys can be added
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, _ := ssh.Dial("tcp", addr, cfg)
	defer client.Close()

	session, _ := client.NewSession()
	defer session.Close()

	startTerminal(session)
}
