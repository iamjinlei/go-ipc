package protocol

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io"

	"github.com/iamjinlei/go-ipc/transport"
)

var (
	ErrUnexpectedMsgSize = errors.New("unexpected message size")
	ErrInvalidMsgType    = errors.New("invalid message type")
)

type MsgType byte

const (
	TypeHandshake MsgType = 1
	TypeAck       MsgType = 2
	TypeRequest   MsgType = 3
	TypeResponse  MsgType = 4
)

func (t MsgType) encode() byte {
	return byte(t)
}

type Msg struct {
	t MsgType
	d []byte
}

func (m *Msg) Type() MsgType {
	return m.t
}

func (m *Msg) Handshake() (*Handshake, error) {
	if m.t != TypeHandshake {
		return nil, ErrInvalidMsgType
	}

	var h Handshake
	if err := decode(m.d, &h); err != nil {
		return nil, err
	}

	return &h, nil
}

func (m *Msg) Ack() (*Ack, error) {
	if m.t != TypeAck {
		return nil, ErrInvalidMsgType
	}

	var a Ack
	if err := decode(m.d, &a); err != nil {
		return nil, err
	}

	return &a, nil
}

func (m *Msg) Request() (*Request, error) {
	if m.t != TypeRequest {
		return nil, ErrInvalidMsgType
	}

	var r Request
	if err := decode(m.d, &r); err != nil {
		return nil, err
	}

	return &r, nil
}

func (m *Msg) Response() (*Response, error) {
	if m.t != TypeResponse {
		return nil, ErrInvalidMsgType
	}

	var r Response
	if err := decode(m.d, &r); err != nil {
		return nil, err
	}

	return &r, nil
}

type msg interface {
	Type() MsgType
}

func WriteMsg(w io.Writer, e msg) error {
	pkt, err := encode(e)
	if err != nil {
		return err
	}

	return transport.WritePacket(w, pkt)
}

func ReadMsg(r io.Reader) (*Msg, error) {
	pkt, err := transport.ReadPacket(r)
	if err != nil {
		return nil, err
	}

	if len(pkt) < 1 {
		return nil, ErrUnexpectedMsgSize
	}

	return &Msg{
		t: MsgType(pkt[0]),
		d: pkt[1:],
	}, nil
}

func encode(e msg) ([]byte, error) {
	var buf bytes.Buffer

	if err := buf.WriteByte(e.Type().encode()); err != nil {
		return nil, err
	}

	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(e); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func decode(data []byte, e any) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(e)
}
