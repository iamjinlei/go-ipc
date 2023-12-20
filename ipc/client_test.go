package ipc

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	data := "test data"
	respData := "test response data"
	respErr := errors.New("test response error")
	r := NewOutgoingMessage([]byte(data))

	ch := make(chan bool)
	go func() {
		defer close(ch)

		<-ch
		assert.Equal(t, data, string(r.Data()))
		d, err := r.Response()
		assert.ErrorIs(t, respErr, err)
		assert.Equal(t, respData, string(d))
	}()

	r.setResponse([]byte(respData), respErr)
	ch <- true
	<-ch
}
