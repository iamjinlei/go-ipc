package protocol_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/iamjinlei/go-ipc/protocol"
)

func TestRequest(t *testing.T) {
	data := "test request data"

	// Encode-decode
	enc, err := protocol.EncodeRequest([]byte(data))
	assert.NoError(t, err)
	r, err := protocol.DecodeRequest(enc)
	assert.NoError(t, err)
	assert.Equal(t, data, string(r.Data()))

	// Write-read
	buf := bytes.NewBuffer(nil)
	require.NoError(t, protocol.WriteRequest(buf, []byte(data)))

	r, err = protocol.ReadRequest(bytes.NewReader(buf.Bytes()))
	assert.NoError(t, err)
	assert.Equal(t, data, string(r.Data()))
}
