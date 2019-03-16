package main

func main() {
	// conn, channels, globalReqs

	type Request struct {
		Type      string
		WantReply bool
		Payload   []byte
	}

	for req := range globalReqs {
		// to accept, set accept to true, otherwise false
		req.Reply(accept, payload)
	}
}
