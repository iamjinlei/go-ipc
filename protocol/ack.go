package protocol

type Ack struct {
	O bool
	M string
}

func NewAck(ok bool, msg string) *Ack {
	return &Ack{
		O: ok,
		M: msg,
	}
}

func (a *Ack) Ok() bool {
	return a.O
}

func (a *Ack) Msg() string {
	return a.M
}

func (a *Ack) Type() MsgType {
	return TypeAck
}
