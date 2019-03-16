package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

var USAGE = `Usage: %s <ssh server address> <ssh user> <path to ssh private key> <command to run>.

Example: %s 127.0.0.1:22 myuser ~/.ssh/id_rsa date
`

func main() {
	if len(os.Args) < 4 {
		fmt.Printf(USAGE, os.Args[0], os.Args[0])
		os.Exit(1)
	}
	addr := os.Args[1]
	user := os.Args[2]
	privKeyPath := os.Args[3]
	cmd := strings.Join(os.Args[4:], " ")

	keyBytes, _ := ioutil.ReadFile(privKeyPath)
	key, _ := ssh.ParsePrivateKey(keyBytes)

	cfg := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, _ := ssh.Dial("tcp", addr, cfg)
	defer client.Close()

	session, _ := client.NewSession()
	defer session.Close()

	stdout, _ := session.StdoutPipe()
	stderr, _ := session.StderrPipe()
	stdin, _ := session.StdinPipe()

	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)
	go io.Copy(stdin, os.Stdin)

	session.Start(cmd)
	session.Wait()
}
