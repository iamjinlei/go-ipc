package protocol

type Handshake struct {
	Id string
}

func NewHandshake(id string) *Handshake {
	return &Handshake{
		Id: id,
	}
}

func (h *Handshake) ID() string {
	return h.Id
}

func (h *Handshake) Type() MsgType {
	return TypeHandshake
}
