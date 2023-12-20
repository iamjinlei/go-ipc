package protocol

import (
	"errors"
	"io"

	"github.com/iamjinlei/go-ipc/transport"
)

type Response struct {
	// Use string type instead of error type since the type is lost
	// after gob encoding/decoding
	E string
	D []byte
}

func (r *Response) Data() []byte {
	return r.D
}

func (r *Response) Error() error {
	if r.E == "" {
		return nil
	}

	return errors.New(r.E)
}

func (r *Response) Encode() ([]byte, error) {
	return encode(r)
}

func EncodeResponse(d []byte, e error) ([]byte, error) {
	errStr := ""
	if e != nil {
		errStr = e.Error()
	}

	r := &Response{
		E: errStr,
		D: d,
	}
	return r.Encode()
}

func DecodeResponse(d []byte) (*Response, error) {
	var resp Response
	if err := decode(d, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func WriteResponse(w io.Writer, d []byte, e error) error {
	pkt, err := EncodeResponse(d, e)
	if err != nil {
		return err
	}

	return transport.WritePacket(w, pkt)
}

func ReadResponse(r io.Reader) (*Response, error) {
	pkt, err := transport.ReadPacket(r)
	if err != nil {
		return nil, err
	}

	return DecodeResponse(pkt)
}
