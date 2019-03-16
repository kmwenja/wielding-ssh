Wielding SSH Talk
=================

Material for the Wielding SSH Talk held at the [Nairobi LUG](https://groups.google.com/forum/#!forum/nairobi-gnu) meetup on 2nd March 2019.

The talk was recorded and published to [Youtube](https://youtu.be/k6gAf6WU5fo).

You can also view the talk online by going to https://talks.godoc.org/github.com/kmwenja/wielding-ssh/talk.slide.

Run the presentation locally:
-----------------------------

1. Setup a Go environment: https://golang.org/doc/install
2. `go get -u golang.org/x/tools/cmd/present`
3. In this directory, run `present -orighost localhost -notes`
4. Visit http://localhost:3999 in your browser.

[Present Docs](https://godoc.org/golang.org/x/tools/cmd/present)

Run the demo programs:
----------------------

*The demo examples have no error checking and are not prod ready!*

1. Cd into a demo example e.g. `cd demo/ssh_pubkey/`
2. Build the source code: `go build`. This will build a binary named after the directory that you're in.
3. Run the built binary e.g. `./ssh_pubkey`. Most of the examples take in arguments in order to run ok. Run the binary without any arguments to see what arguments to provide (usage).

References:
-----------

Docs - https://godoc.org/golang.org/x/crypto/ssh
Code - https://github.com/golang/crypto/tree/master/ssh
High Level APIS - https://github.com/gliderlabs/ssh
RFC 4251 (start here) - https://tools.ietf.org/html/rfc4251
SSH Chat - https://github.com/shazow/ssh-chat
SSH Tron - https://github.com/zachlatta/sshtron
Gitea - https://github.com/go-gitea/gitea
Gogs - https://github.com/gogs/gogs
Hashicorp's Vault - https://github.com/hashicorp/vault
Gravitational's Teleport - https://github.com/gravitational/teleport
