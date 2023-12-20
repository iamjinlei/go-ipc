package transport

import (
	"encoding/binary"
	"io"
)

func WritePacket(w io.Writer, d []byte) error {
	sz := make([]byte, 4)
	binary.LittleEndian.PutUint32(sz, uint32(len(d)))
	if err := write(w, sz); err != nil {
		return err
	}

	return write(w, d)
}

func ReadPacket(r io.Reader) ([]byte, error) {
	sz := make([]byte, 4)
	if err := read(r, sz); err != nil {
		return nil, err
	}
	n := binary.LittleEndian.Uint32(sz)
	data := make([]byte, n)
	if err := read(r, data); err != nil {
		return nil, err
	}

	return data, nil
}
