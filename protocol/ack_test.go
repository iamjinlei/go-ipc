package protocol_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/iamjinlei/go-ipc/protocol"
)

func TestAck(t *testing.T) {
	data := "ack message"
	for _, ok := range []bool{true, false} {
		buf := bytes.NewBuffer(nil)
		a := protocol.NewAck(ok, data)
		require.NoError(t, protocol.WriteMsg(buf, a))

		msg, err := protocol.ReadMsg(bytes.NewReader(buf.Bytes()))
		assert.NoError(t, err)
		a, err = msg.Ack()
		assert.NoError(t, err)
		assert.Equal(t, ok, a.Ok())
		assert.Equal(t, data, a.Msg())
	}
}
