package ipc

type IncomingMessage struct {
	doneCh chan bool
	d      []byte
	e      error
	resp   []byte
}

func newIncomingMessage(d []byte) *IncomingMessage {
	return &IncomingMessage{
		d:      d,
		doneCh: make(chan bool),
	}
}

func (r *IncomingMessage) Data() []byte {
	return r.d
}

func (r *IncomingMessage) SetResponse(d []byte, e error) {
	defer close(r.doneCh)
	r.resp = d
	r.e = e
}

func (r *IncomingMessage) response() ([]byte, error) {
	<-r.doneCh
	return r.resp, r.e
}
