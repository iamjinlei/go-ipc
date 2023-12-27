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

	r := protocol.NewRequest([]byte(data))
	buf := bytes.NewBuffer(nil)
	require.NoError(t, protocol.WriteMsg(buf, r))

	msg, err := protocol.ReadMsg(bytes.NewReader(buf.Bytes()))
	assert.NoError(t, err)
	r, err = msg.Request()
	assert.NoError(t, err)
	assert.Equal(t, data, string(r.Data()))
}
