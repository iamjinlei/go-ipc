package transport

import (
	"io"
)

func read(r io.Reader, buf []byte) error {
	nb := 0
	for nb < len(buf) {
		n, err := r.Read(buf[nb:])
		if err != nil {
			return err
		}
		nb += n
	}
	return nil
}

func write(w io.Writer, data []byte) error {
	nb := 0
	for nb < len(data) {
		n, err := w.Write(data[nb:])
		if err != nil {
			return err
		}
		nb += n
	}
	return nil
}
