package util

import (
	"io"
)

func CloseReader(r io.Reader) error {
	if closer, ok := r.(io.ReadCloser); ok {
		return closer.Close()
	}
	return nil
}
