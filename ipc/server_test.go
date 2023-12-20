package ipc

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	data := "test data"
	respData := "test response data"
	respErr := errors.New("test response error")
	r := newIncomingMessage([]byte(data))

	ch := make(chan bool)
	go func() {
		defer close(ch)

		<-ch
		assert.Equal(t, data, string(r.Data()))
		d, err := r.response()
		assert.ErrorIs(t, respErr, err)
		assert.Equal(t, respData, string(d))
	}()

	r.SetResponse([]byte(respData), respErr)
	ch <- true
	<-ch
}
