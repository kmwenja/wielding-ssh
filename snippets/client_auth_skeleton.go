package main

import (
	"net"

	"golang.org/x/crypto/ssh"
)

func main() {
	cfg := &ssh.ClientConfig{
		User: "",
		Auth: []ssh.AuthMethod{
			ssh.Password("secret"),
			ssh.PasswordCallback(func() (string, error) {
				return "secret", nil
			}),
			ssh.PublicKeys(key1, key2), // put as many keys as you like here
			ssh.PublicKeysCallback(func() ([]ssh.Signer, error) {
				return []ssh.Signer{key1, key2}, nil
			}),
			ssh.KeyboardInteractiveChallenge(challengeFunc),
		},
		// HostKeyCallback: ssh.FixedHostKey(key),
		// HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil // accepted
			// return error to refuse
		},
		// there's a bunch of other cool stuff here like cipher suites
	}
}
