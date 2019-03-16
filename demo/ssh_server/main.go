package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"

	ptyUtil "github.com/kr/pty"
	"golang.org/x/crypto/ssh"
)

var USAGE = `Usage: %s <listen addr>

This server will, when the client logs in, start a shell for the user that run the
server in the same directory the server is run.

To ssh to this server, either use the username:password combination of "testuser:hello"
or use the private keyfile "test" that is in the source code directory for this server.

For example (assuming the server is running above:

Running the server: %s 127.0.0.1:2022

With password: ssh -p 2022 testuser@127.0.0.1
With public key: ssh -p 2022 -i test 127.0.0.1
`

func main() {
	if len(os.Args) < 2 {
		fmt.Printf(USAGE, os.Args[0], os.Args[0])
		os.Exit(1)
	}
	addr := os.Args[1]

	authorizedPubKeyBytes, _ := ioutil.ReadFile("test.pub")
	authorizedPubKey, _, _, _, _ := ssh.ParseAuthorizedKey(authorizedPubKeyBytes)

	cfg := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			if c.User() == "testuser" && string(pass) == "hello" {
				return nil, nil // allow login
			}
			return nil, fmt.Errorf("password rejected for %q", c.User())
		},
		PublicKeyCallback: func(c ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error) {
			// check if the fingerprints match
			if string(pubKey.Marshal()) == string(authorizedPubKey.Marshal()) {
				return nil, nil // allow login
			}
			return nil, fmt.Errorf("public key rejected for %q", c.User())
		},
	}

	hostPrivKeyBytes, _ := ioutil.ReadFile("host_key")
	hostPrivKey, _ := ssh.ParsePrivateKey(hostPrivKeyBytes)
	cfg.AddHostKey(hostPrivKey)

	listener, _ := net.Listen("tcp", addr)
	defer listener.Close()
	for {
		c, _ := listener.Accept()
		go handleNewConnection(c, cfg)
	}
}

func handleNewConnection(c net.Conn, cfg *ssh.ServerConfig) {
	defer c.Close()

	// handshake performed here
	conn, chans, reqs, _ := ssh.NewServerConn(c, cfg)
	defer conn.Close()

	// connection established!

	// don't handle global requests for now
	go ssh.DiscardRequests(reqs)

	// handle incoming channels
	for newChannel := range chans {
		log.Printf("Channel Type: %s", newChannel.ChannelType())

		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.Prohibited, "unsupported channel type")
			continue
		}

		channel, requests, _ := newChannel.Accept()
		go handleSessionChannel(channel, requests)
	}
}

type ptyReqMsg struct {
	Term     string
	Columns  uint32
	Rows     uint32
	Width    uint32
	Height   uint32
	Modelist string
}

type winChgMsg struct {
	Columns uint32
	Rows    uint32
	Width   uint32
	Height  uint32
}

type envMsg struct {
	Name  string
	Value string
}

type execMsg struct {
	Command string
}

func handleSessionChannel(channel ssh.Channel, reqs <-chan *ssh.Request) {
	defer channel.Close()
	// TODO send exit-status on return

	// pty-req, env, shell, window-change
	// pty-req, env, exec, window-change
	// exec

	c := exec.Cmd{}
	var m ptyReqMsg
	var pty *os.File

	syncWinSize := func() {
		for req := range reqs {
			req.Reply(req.Type == "window-change", nil)
			var wm winChgMsg
			ssh.Unmarshal(req.Payload, &wm)
			log.Printf("Win Change Req: %s %v", req.Type, req.Payload)
			ptyUtil.Setsize(pty, &ptyUtil.Winsize{
				Rows: uint16(wm.Rows),
				Cols: uint16(wm.Columns),
				X:    uint16(wm.Width),
				Y:    uint16(wm.Height),
			})
		}
	}

	startPty := func() {
		// TODO handle termios
		pty, _ = ptyUtil.StartWithSize(&c, &ptyUtil.Winsize{
			Rows: uint16(m.Rows),
			Cols: uint16(m.Columns),
			X:    uint16(m.Width),
			Y:    uint16(m.Height),
		})
		go syncWinSize()
		go io.Copy(pty, channel)
		io.Copy(channel, pty)
	}

	for req := range reqs {
		log.Printf("Req: %s", req.Type)

		switch req.Type {
		case "pty-req":
			req.Reply(true, nil)
			ssh.Unmarshal(req.Payload, &m)
			c.Env = append(c.Env, fmt.Sprintf("TERM=%s", m.Term))
		case "env":
			req.Reply(true, nil)
			var em envMsg
			ssh.Unmarshal(req.Payload, &em)
			c.Env = append(c.Env, fmt.Sprintf("%s=%s", em.Name, em.Value))
		case "shell":
			req.Reply(true, nil)
			c.Path = os.Getenv("SHELL")
			startPty()
			return
		case "exec":
			req.Reply(true, nil)
			var em execMsg
			ssh.Unmarshal(req.Payload, &em)
			path, _ := exec.LookPath(em.Command)
			c.Path = path
			log.Printf("Current pty-req: %v", m)
			if m.Term != "" {
				startPty()
				return
			} else {
				stdin, _ := c.StdinPipe()
				stdout, _ := c.StdoutPipe()
				stderr, _ := c.StderrPipe()
				go io.Copy(channel, stdout)
				go io.Copy(channel, stderr)
				go io.Copy(stdin, channel)
				c.Run()
				return
			}
		}
	}
}
