package protocol

type Request struct {
	D []byte
}

func NewRequest(d []byte) *Request {
	return &Request{
		D: d,
	}
}

func (r *Request) Data() []byte {
	return r.D
}

func (r *Request) Type() MsgType {
	return TypeRequest
}
