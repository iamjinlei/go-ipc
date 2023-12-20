<h1 align="center">go-ipc (Inter-Process Communication)</h>

[![Go Report Card](https://goreportcard.com/badge/github.com/iamjinlei/go-ipc)](https://goreportcard.com/report/github.com/iamjinlei/go-ipc)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Go Reference](https://pkg.go.dev/badge/github.com/iamjinlei/go-ipc.svg)](https://pkg.go.dev/github.com/iamjinlei/go-ipc)

## Intro

[go-ipc](https://github.com/iamjinlei/go-ipc) provides an abstraction of the [Unix Domain Socket](https://en.wikipedia.org/wiki/Unix_domain_socket).
The SDK can be integrated into an application to allow inter-process data exchange on the same operating system.
A n-producer-1-consumer pattern can be achieved using go-ipc.
The n producers use the client to communicate with the master which consume the messages.
If multiple processes are in *master* or *dual* mode,  a leader is elected through local file system's locking mechanism.
When the master goes down, candidates run a new election to select the new master.
The auto reelection ensures the group communication never breaks.
If auto leader election/migration is not desired, simply enable only one master in the group.

There are 3 operating modes:
* Client: a node can send data to another node which is operating in the *master* mode.
* Master: a node can process data received from other clients.
* Dual: a node behaves as client and master at the same time.

## Install
```bash
go get github.com/iamjinlei/go-ipc
```

## Example

The following code snippet provides a incomplete example to demonstrate the idea.

```go
nd, _ := ipc.NewNode(
    ctx,
    ipc.Config{
	    GroupName:     "groupecho",
		EnableDebug:   false,
		IgnoreDialErr: true,
	},
)

go func() {
	for {
		select {
		case msg := <-nd.ReceiveCh():
			fmt.Printf("received message: \"%v\"\n", string(msg.Data()))
			msg.SetResponse([]byte("ok"), nil)
	    case err := <-nd.ErrCh():
		    fmt.Printf("error: %v\n", err)
		}
	}
}()

msg := ipc.NewOutgoingMessage([]byte("hello!"))
nd.SendCh() <- msg
resp, _ := msg.Response()
fmt.Printf("response: \"%v\"\n", string(resp))
}
```

A runnable program can be found under the [examples](https://github.com/iamjinlei/go-ipc/tree/main/examples/groupecho) folder.
Start multiple instances of the groupecho program and randomly kill any to see how the master reelection keeps communication intact.
