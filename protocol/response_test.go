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

	// Encode-decode
	enc, err := protocol.EncodeResponse([]byte(data), respErr)
	assert.NoError(t, err)
	r, err := protocol.DecodeResponse(enc)
	assert.NoError(t, err)
	assert.Equal(t, data, string(r.Data()))
	assert.Equal(t, respErr.Error(), r.Error().Error())

	// Write-read with no error
	buf := bytes.NewBuffer(nil)
	assert.NoError(t, protocol.WriteResponse(buf, []byte(data), nil))
	r, err = protocol.ReadResponse(bytes.NewReader(buf.Bytes()))
	assert.NoError(t, err)
	assert.Equal(t, data, string(r.Data()))
	assert.NoError(t, r.Error())

	// Write-read with error
	buf = bytes.NewBuffer(nil)
	assert.NoError(t, protocol.WriteResponse(buf, []byte(data), respErr))
	r, err = protocol.ReadResponse(bytes.NewReader(buf.Bytes()))
	assert.NoError(t, err)
	assert.Equal(t, data, string(r.Data()))
	assert.Equal(t, respErr.Error(), r.Error().Error())
}
