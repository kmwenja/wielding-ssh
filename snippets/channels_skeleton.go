package main

import "golang.org/x/crypto/ssh"

func main() {
	// conn, channels, globalReqs

	for newChannel := range channels {
		// 'session' channels power shells, command execs, file copying
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.Prohibited, "unsupported channel type")
		}
		channel, channelReqs, _ := newChannel.Accept()

		// channelReqs have to be processed or the channel will hang
		// channel.Read to read from the channel (io.Reader)
		// channel.Write to write to the channel (io.Writer)
		// channel.SendRequest() to send channel requests to the other side
	}
}
