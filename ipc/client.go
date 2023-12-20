package ipc

type OutgoingMessage struct {
	d      []byte
	resp   []byte
	e      error
	doneCh chan bool
}

func NewOutgoingMessage(d []byte) *OutgoingMessage {
	return &OutgoingMessage{
		d:      d,
		doneCh: make(chan bool),
	}
}

func (r *OutgoingMessage) Data() []byte {
	return r.d
}

func (r *OutgoingMessage) setResponse(d []byte, e error) {
	defer close(r.doneCh)
	r.resp = d
	r.e = e
}

func (r *OutgoingMessage) Response() ([]byte, error) {
	<-r.doneCh
	return r.resp, r.e
}
