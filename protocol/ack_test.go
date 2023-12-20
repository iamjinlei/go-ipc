package protocol_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/iamjinlei/go-ipc/protocol"
)

func TestAck(t *testing.T) {
	msg := "ack message"
	for _, ok := range []bool{true, false} {
		// Encode-decode
		enc, err := protocol.EncodeAck(ok, msg)
		assert.NoError(t, err)
		a, err := protocol.DecodeAck(enc)
		assert.NoError(t, err)
		assert.Equal(t, ok, a.Ok())
		assert.Equal(t, msg, a.Msg())

		// Write-read
		buf := bytes.NewBuffer(nil)
		require.NoError(t, protocol.WriteAck(buf, ok, msg))

		a, err = protocol.ReadAck(bytes.NewReader(buf.Bytes()))
		assert.NoError(t, err)
		assert.Equal(t, ok, a.Ok())
		assert.Equal(t, msg, a.Msg())
	}
}
