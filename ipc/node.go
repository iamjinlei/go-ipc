package ipc

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/google/uuid"
	cmap "github.com/orcaman/concurrent-map/v2"

	"github.com/iamjinlei/go-ipc/flock"
	"github.com/iamjinlei/go-ipc/protocol"
)

var (
	ErrHandshake = errors.New("handshake error")
)

type Mode string

const (
	UnknownMode Mode = ""
	MasterOnly  Mode = "master"
	ClientOnly  Mode = "client"
	Dual        Mode = "dual"
)

type Node struct {
	ctx     context.Context
	id      string
	cfg     *Config
	lk      *flock.FileLock
	udsPath string
	inCh    chan *IncomingMessage
	outCh   chan *OutgoingMessage
	errCh   chan error
}

func NewNode(
	ctx context.Context,
	c Config,
) (*Node, error) {
	cfg := &c
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	nd := &Node{
		ctx:     ctx,
		id:      fmt.Sprintf("%x", uuid.New()),
		cfg:     cfg,
		lk:      flock.New(fmt.Sprintf("/tmp/%v.lock", cfg.GroupName)),
		udsPath: fmt.Sprintf("/tmp/%v.sock", cfg.GroupName),
		inCh:    make(chan *IncomingMessage, cfg.BufferSize),
		outCh:   make(chan *OutgoingMessage, cfg.BufferSize),
		errCh:   make(chan error, cfg.BufferSize),
	}

	switch cfg.Mode {
	case MasterOnly:
		go nd.serverLoop()
	case ClientOnly:
		go nd.clientLoop()
	case Dual:
		go nd.serverLoop()
		go nd.clientLoop()
	}

	return nd, nil
}

func (n *Node) SendCh() chan *OutgoingMessage {
	return n.outCh
}

func (n *Node) ReceiveCh() chan *IncomingMessage {
	return n.inCh
}

func (n *Node) ErrCh() chan error {
	return n.errCh
}

func (n *Node) clientLoop() {
	n.cfg.log("[Client] processing loop started")

	for ; true; time.Sleep(time.Second) {
		var d net.Dialer
		conn, err := d.DialContext(n.ctx, "unix", n.udsPath)
		if err != nil {
			if err == context.Canceled {
				return
			}

			if n.cfg.IgnoreDialErr {
				if errors.Is(err, syscall.ECONNREFUSED) ||
					errors.Is(err, fs.ErrNotExist) {
					continue
				}
			}

			n.sendErr(err)
			continue
		}

		n.cfg.log("[Client] connection established")

		func(c net.Conn) {
			defer c.Close()

			n.cfg.setWriteDeadline(c)
			if err := protocol.WriteHandshake(c, n.id); err != nil {
				n.sendErr(err)
				return
			}

			n.cfg.setReadDeadline(c)
			if ack, err := protocol.ReadAck(c); err != nil {
				n.sendErr(err)
				return
			} else if !ack.Ok() {
				n.sendErr(ErrHandshake)
				return
			}

			n.cfg.log("[Client] handshake succeeded")

			transceive(n.ctx, c, n.cfg, n.outCh)
		}(conn)
	}
}

func transceive(
	ctx context.Context,
	conn net.Conn,
	cfg *Config,
	outCh chan *OutgoingMessage,
) {
	cfg.log("[Client] transceive loop started")

	for {
		select {
		case <-ctx.Done():
			cfg.log("[Client] context done")
			return

		case req := <-outCh:
			cfg.log("[Client] outgoing message received from queue")

			cfg.setWriteDeadline(conn)
			if err := protocol.WriteRequest(conn, req.Data()); err != nil {
				req.setResponse(nil, err)
				return
			}

			cfg.log("[Client] outgoing message sent")

			cfg.setReadDeadline(conn)
			resp, err := protocol.ReadResponse(conn)
			if err != nil {
				req.setResponse(nil, err)
				return
			}

			cfg.log("[Client] outgoing message response received")

			req.setResponse(resp.Data(), resp.Error())
		}
	}
}

func (n *Node) serverLoop() {
	n.cfg.log("[Server] processing loop started")

	for ; true; time.Sleep(time.Second) {
		locked, err := n.lk.Lock(n.ctx)
		if err != nil {
			if err == context.Canceled {
				return
			}

			n.sendErr(err)
			continue
		}

		if !locked {
			continue
		}

		n.cfg.log("[Server] master lock acquired")

		// Exclusive lock acquired. Start the server loop.
		connMap := cmap.New[net.Conn]()
		func() {
			defer n.lk.Unlock()
			doneCh := make(chan bool)
			defer close(doneCh)

			if err := os.RemoveAll(n.udsPath); err != nil {
				n.sendErr(err)
				return
			}
			lis, err := net.Listen("unix", n.udsPath)
			if err != nil {
				n.sendErr(err)
				return
			}

			n.cfg.log("[Server] listened on %v", n.udsPath)

			// Ensure listener is closed after ctx done or function exits.
			// This should unblock the accept call below.
			go func() {
				select {
				case <-n.ctx.Done():
				case <-doneCh:
				}
				lis.Close()
			}()

			// Accept incoming connections.
			for {
				conn, err := lis.Accept()
				if err != nil {
					n.sendErr(err)
					return
				}

				n.cfg.log("[Server] new connection accepted")

				go func(c net.Conn) {
					defer c.Close()

					n.cfg.setReadDeadline(c)
					h, err := protocol.ReadHandshake(c)
					if err != nil {
						n.sendErr(err)
						return
					}

					n.cfg.setWriteDeadline(c)
					if err := protocol.WriteAck(c, true, "ok"); err != nil {
						n.sendErr(err)
						return
					}

					n.cfg.log("[Server] handshake succeeded")

					peerId := h.ID()
					connMap.Set(peerId, c)
					defer connMap.Remove(peerId)

					if err := serve(
						n.ctx,
						c,
						n.cfg,
						n.inCh,
					); err != nil {
						n.sendErr(err)
					}
				}(conn)
			}
		}()
	}
}

func serve(
	ctx context.Context,
	conn net.Conn,
	cfg *Config,
	inCh chan *IncomingMessage,
) error {
	cfg.log("[Server] serve loop started")

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		cfg.setReadDeadline(conn)
		req, err := protocol.ReadRequest(conn)
		if err != nil {
			return err
		}

		cfg.log("[Server] incoming message received")

		r := newIncomingMessage(req.Data())
		inCh <- r
		cfg.log("[Server] incoming message sent to queue")

		d, err := r.response()
		cfg.setWriteDeadline(conn)
		if err := protocol.WriteResponse(conn, d, err); err != nil {
			return err
		}

		cfg.log("[Server] incoming message response sent")
	}
}

func (n *Node) sendErr(err error) {
	if len(n.errCh) == n.cfg.BufferSize {
		// Pop out the oldest unhandled error in case the buffer is full.
		// This allows application to ignore err handling without blocking
		// the channel.
		select {
		case <-n.errCh:
		default:
		}
	}
	n.errCh <- err
}
