package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"unsafe"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

var USAGE = `Usage: %s <ssh server address> <ssh user> <ssh password>

Example: %s 127.0.0.1:22 myuser mypassword
`

func main() {
	if len(os.Args) < 4 {
		fmt.Printf(USAGE, os.Args[0], os.Args[0])
		os.Exit(1)
	}
	addr, user, password := os.Args[1], os.Args[2], os.Args[3]

	// configure the connection
	cfg := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
			// ssh.PasswordCallback(prompt func() (secret string, err error))
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // disable strict host checking
		// HostKeyCallback: ssh.FixedHostKey(key),
	}

	// connect to SSH port
	client, _ := ssh.Dial("tcp", addr, cfg)
	defer client.Close()

	// start a new session
	session, _ := client.NewSession()
	defer session.Close()

	startTerminal(session)
}

func startTerminal(s *ssh.Session) error {
	// get pipes connected to stdin, stdout and stderr on this session
	stdout, _ := s.StdoutPipe()
	stderr, _ := s.StderrPipe()
	stdin, _ := s.StdinPipe()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1, // make sure this is 1, otherwise you won't see what you're typing
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	// get the current terminal size
	h, w, _ := getWinSize()

	// ask for a remote pty running the same terminal emualor as the user's
	// with the same height and width as the user's current terminal size
	s.RequestPty(os.Getenv("TERM"), h, w, modes)

	// start a shell on that pty
	s.Shell()

	// reset the local terminal to "raw" mode so that bytes sent to
	// stdin are sent as-is to the remote pty
	oldState, _ := terminal.MakeRaw(int(os.Stdin.Fd()))
	// after we exit, restore the previous mode of the terminal
	defer terminal.Restore(int(os.Stdin.Fd()), oldState)

	// listen out for signals, the SIGWINCH is used to sync window sizes with the remote pty
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGWINCH)
	defer signal.Stop(signals)
	go func() {
		for sig := range signals {
			switch sig {
			case syscall.SIGWINCH:
				h, w, _ := getWinSize()
				s.WindowChange(h, w)
			default:
				log.Printf("unknown signal: %v", sig)
			}
		}
	}()

	// redirect stdout and stderr to local stdout and stderr resp.
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)

	// redirect the local stdin to the remote pty's stdin
	io.Copy(stdin, os.Stdin)
	stdin.Close()
	return nil
}

// getWinSize retrieves the current terminal's dimensions.
func getWinSize() (h, w int, err error) {
	type wsz struct {
		Rows uint16
		Cols uint16
		X    uint16
		Y    uint16
	}
	var ws wsz
	// to get current terminal window size one has to make a system call
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		os.Stdin.Fd(),
		syscall.TIOCGWINSZ,
		uintptr(unsafe.Pointer(&ws)))
	if errno != 0 {
		return -1, -1, fmt.Errorf("could not get window size: %s", syscall.Errno(errno))
	}

	return int(ws.Rows), int(ws.Cols), nil
}
