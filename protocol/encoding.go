package protocol

import (
	"bytes"
	"encoding/gob"
)

func encode(e any) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(e); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func decode(data []byte, e any) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(e)
}
