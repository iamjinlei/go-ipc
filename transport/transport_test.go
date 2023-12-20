package transport_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iamjinlei/go-ipc/transport"
)

func TestTransport(t *testing.T) {
	data := "test request data"
	buf := bytes.NewBuffer(nil)
	assert.NoError(t, transport.WritePacket(buf, []byte(data)))
	d, err := transport.ReadPacket(bytes.NewReader(buf.Bytes()))
	assert.NoError(t, err)
	assert.Equal(t, data, string(d))
}
