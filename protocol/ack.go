package protocol

import (
	"io"

	"github.com/iamjinlei/go-ipc/transport"
)

type Ack struct {
	O bool
	M string
}

func (a *Ack) Ok() bool {
	return a.O
}

func (a *Ack) Msg() string {
	return a.M
}

func (a *Ack) Encode() ([]byte, error) {
	return encode(a)
}

func EncodeAck(ok bool, msg string) ([]byte, error) {
	a := &Ack{
		O: ok,
		M: msg,
	}
	return a.Encode()
}

func DecodeAck(d []byte) (*Ack, error) {
	var a Ack
	if err := decode(d, &a); err != nil {
		return nil, err
	}

	return &a, nil
}

func WriteAck(w io.Writer, ok bool, msg string) error {
	pkt, err := EncodeAck(ok, msg)
	if err != nil {
		return err
	}

	return transport.WritePacket(w, pkt)
}

func ReadAck(r io.Reader) (*Ack, error) {
	pkt, err := transport.ReadPacket(r)
	if err != nil {
		return nil, err
	}

	return DecodeAck(pkt)
}
