package protocol

import (
	"errors"
)

type Response struct {
	// Use string type instead of error type since the type is lost
	// after gob encoding/decoding
	E string
	D []byte
}

func NewResponse(d []byte, e error) *Response {
	errStr := ""
	if e != nil {
		errStr = e.Error()
	}

	return &Response{
		E: errStr,
		D: d,
	}
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

func (r *Response) Type() MsgType {
	return TypeResponse
}
