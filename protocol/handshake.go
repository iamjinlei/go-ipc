package protocol

import (
	"io"

	"github.com/iamjinlei/go-ipc/transport"
)

type Handshake struct {
	Id string
}

func (h *Handshake) ID() string {
	return h.Id
}

func (h *Handshake) Encode() ([]byte, error) {
	return encode(h)
}

func EncodeHandshake(id string) ([]byte, error) {
	h := &Handshake{
		Id: id,
	}
	return h.Encode()
}

func DecodeHandshake(d []byte) (*Handshake, error) {
	var h Handshake
	if err := decode(d, &h); err != nil {
		return nil, err
	}

	return &h, nil
}

func WriteHandshake(w io.Writer, id string) error {
	pkt, err := EncodeHandshake(id)
	if err != nil {
		return err
	}

	return transport.WritePacket(w, pkt)
}

func ReadHandshake(r io.Reader) (*Handshake, error) {
	pkt, err := transport.ReadPacket(r)
	if err != nil {
		return nil, err
	}

	return DecodeHandshake(pkt)
}
