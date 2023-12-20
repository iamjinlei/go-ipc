package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/iamjinlei/go-ipc/ipc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	nd, err := ipc.NewNode(
		ctx,
		ipc.Config{
			GroupName:     "groupecho",
			EnableDebug:   false,
			IgnoreDialErr: true,
		},
	)
	if err != nil {
		fmt.Printf("Error creating node: %v\n", err)
		return
	}

	rand.Seed(time.Now().UnixNano())
	whoami := rand.Intn(16)

	fmt.Printf("Node#%v started\n", whoami)
	ticker := time.NewTicker(time.Second)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case msg := <-nd.ReceiveCh():
				fmt.Printf("[#%v] received message: \"%v\"\n",
					whoami,
					string(msg.Data()),
				)
				msg.SetResponse(
					[]byte(fmt.Sprintf("Node#%v: ok", whoami)),
					nil,
				)
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return

		case err := <-nd.ErrCh():
			fmt.Printf("[#%v] error: %v\n", whoami, err)

		case <-ticker.C:
			msg := ipc.NewOutgoingMessage(newMsg(whoami))
			fmt.Printf("[#%v] new message: %v\n",
				whoami,
				string(msg.Data()),
			)

			nd.SendCh() <- msg
			resp, err := msg.Response()

			if err != nil {
				fmt.Printf("[#%v] received error: %v\n",
					whoami,
					err,
				)
			} else {
				fmt.Printf("[#%v] received response: \"%v\"\n",
					whoami,
					string(resp),
				)
			}
		}
	}
}

func newMsg(id int) []byte {
	return []byte(fmt.Sprintf(
		"Node#%v: my clock is %v",
		id,
		time.Now().UnixMilli(),
	))
}
