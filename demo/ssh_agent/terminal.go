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

func startTerminal(s *ssh.Session) error {
	// get pipes connected to stdin, stdout and stderr on this session
	stdout, err := s.StdoutPipe()
	if err != nil {
		return fmt.Errorf("could not get stdout: %s", err)
	}

	stderr, err := s.StderrPipe()
	if err != nil {
		return fmt.Errorf("could not get stderr: %s", err)
	}

	stdin, err := s.StdinPipe()
	if err != nil {
		return fmt.Errorf("could not get stdin: %s", err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1, // make sure this is 1, otherwise you won't see what you're typing
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	// get the current terminal size
	h, w, err := getWinSize()
	if err != nil {
		return err
	}

	// ask for a remote pty running the same terminal emualor as the user's
	// with the same height and width as the user's current terminal size
	if err := s.RequestPty(os.Getenv("TERM"), h, w, modes); err != nil {
		return fmt.Errorf("could not get tty: %s", err)
	}

	// start a shell on that pty
	if err := s.Shell(); err != nil {
		return fmt.Errorf("could not start shell: %s", err)
	}

	// reset the local terminal to "raw" mode so that bytes sent to
	// stdin are sent as-is to the remote pty
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal("could not configure stdin: ", err)
	}
	// after we exit, restore the previous mode of the terminal
	defer terminal.Restore(int(os.Stdin.Fd()), oldState)

	// listen out for signals, the SIGWINCH is used to sync window sizes with the remote pty
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGWINCH)
	defer signal.Stop(signals)
	go func() {
		for sig := range signals {
			switch sig {
			case syscall.SIGWINCH:
				h, w, err := getWinSize()
				if err != nil {
					log.Fatal(err)
				}
				s.WindowChange(h, w)
			case os.Interrupt:
				stdin.Write([]byte("\x03"))
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
	return nil
}

// getWinSize retrieves the current terminal's dimensions
func getWinSize() (h, w int, err error) {
	type wsz struct {
		Rows uint16
		Cols uint16
		X    uint16
		Y    uint16
	}
	var ws wsz
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
