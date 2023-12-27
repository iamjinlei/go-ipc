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

	h := protocol.NewHandshake(id)
	buf := bytes.NewBuffer(nil)
	require.NoError(t, protocol.WriteMsg(buf, h))

	msg, err := protocol.ReadMsg(bytes.NewReader(buf.Bytes()))
	assert.NoError(t, err)
	h, err = msg.Handshake()
	assert.NoError(t, err)
	assert.Equal(t, id, h.ID())
}
