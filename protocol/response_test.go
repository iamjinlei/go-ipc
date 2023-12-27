package protocol_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iamjinlei/go-ipc/protocol"
)

func TestResponse(t *testing.T) {
	data := "test response data"
	respErr := errors.New("test response error")

	// Write-read with no error
	r := protocol.NewResponse([]byte(data), nil)
	buf := bytes.NewBuffer(nil)
	assert.NoError(t, protocol.WriteMsg(buf, r))
	msg, err := protocol.ReadMsg(bytes.NewReader(buf.Bytes()))
	assert.NoError(t, err)
	r, err = msg.Response()
	assert.NoError(t, err)
	assert.Equal(t, data, string(r.Data()))
	assert.NoError(t, r.Error())

	// Write-read with error
	r = protocol.NewResponse([]byte(data), respErr)
	buf = bytes.NewBuffer(nil)
	assert.NoError(t, protocol.WriteMsg(buf, r))
	msg, err = protocol.ReadMsg(bytes.NewReader(buf.Bytes()))
	assert.NoError(t, err)
	r, err = msg.Response()
	assert.NoError(t, err)
	assert.Equal(t, data, string(r.Data()))
	assert.Equal(t, respErr.Error(), r.Error().Error())
}
