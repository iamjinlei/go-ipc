package protocol_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/iamjinlei/go-ipc/protocol"
)

func TestHandshake(t *testing.T) {
	id := "test_id"

	// Encode-decode
	enc, err := protocol.EncodeHandshake(id)
	assert.NoError(t, err)
	h, err := protocol.DecodeHandshake(enc)
	assert.NoError(t, err)
	assert.Equal(t, id, h.ID())

	// Write-read
	buf := bytes.NewBuffer(nil)
	require.NoError(t, protocol.WriteHandshake(buf, id))

	h, err = protocol.ReadHandshake(bytes.NewReader(buf.Bytes()))
	assert.NoError(t, err)
	assert.Equal(t, id, h.ID())
}
