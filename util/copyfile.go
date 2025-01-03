package util

import (
	"io"
	"os"
)

const copyChunkSize = 1024 * 1024 // 1M

func CopyStreamToFile(reader io.ReadCloser, fn string) error {
	newf, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer newf.Close()

	buf := make([]byte, 4096)
	for {
		readTotal, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		_, err = newf.Write(buf[:readTotal])
		if err != nil {
			return err
		}
	}
	reader.Close()
	return nil
}
