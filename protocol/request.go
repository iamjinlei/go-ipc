package protocol

import (
	"io"

	"github.com/iamjinlei/go-ipc/transport"
)

type Request struct {
	D []byte
}

func (r *Request) Data() []byte {
	return r.D
}

func (r *Request) Encode() ([]byte, error) {
	return encode(r)
}

func EncodeRequest(d []byte) ([]byte, error) {
	r := &Request{
		D: d,
	}
	return r.Encode()
}

func DecodeRequest(d []byte) (*Request, error) {
	var req Request
	if err := decode(d, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

func WriteRequest(w io.Writer, d []byte) error {
	pkt, err := EncodeRequest(d)
	if err != nil {
		return err
	}

	return transport.WritePacket(w, pkt)
}

func ReadRequest(r io.Reader) (*Request, error) {
	pkt, err := transport.ReadPacket(r)
	if err != nil {
		return nil, err
	}

	return DecodeRequest(pkt)
}
