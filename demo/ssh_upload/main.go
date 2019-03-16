package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh"
)

var USAGE = `Usage: %s <ssh server address> <ssh user> <ssh keypath> <remote file to download> <local file path of the download>

Example: %s 127.0.0.1:22 myuser ~/.ssh/id_rsa remote_file local_file
`

func main() {
	if len(os.Args) < 6 {
		fmt.Printf(USAGE, os.Args[0], os.Args[0])
		os.Exit(1)
	}
	addr, user, keyPath := os.Args[1], os.Args[2], os.Args[3]
	src, dst := os.Args[4], os.Args[5]

	keyBytes, _ := ioutil.ReadFile(keyPath)
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

	stdin, _ := session.StdinPipe()

	srcFile, _ := os.Open(src)
	defer srcFile.Close()

	session.Start(fmt.Sprintf("cat > %q", dst))
	io.Copy(stdin, srcFile)
	stdin.Close()
	session.Wait()
}
