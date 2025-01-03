package util

import (
	"bytes"
	"io"
	"log"
)

type ChainedReader struct {
	Readers []io.Reader
}

func (s *ChainedReader) Read(p []byte) (n int, err error) {
	if len(s.Readers) == 0 {
		return 0, io.EOF
	}

	n, e := s.Readers[0].Read(p)
	if n != 0 && e == nil {
		return n, e
	}

	s.Readers = s.Readers[1:]
	return s.Read(p)
}

func (s *ChainedReader) Close() error {
	log.Println("ChainedReader::Close()\n")
	return nil
}

func (s *ChainedReader) AddBytes(p []byte) {
	s.AddReader(bytes.NewReader(p))
}

func (s *ChainedReader) AddReader(r io.Reader) {
	s.Readers = append(s.Readers, r)
}

func NewChainedReader(arg ...io.Reader) io.Reader {
	r := ChainedReader{
		Readers: arg,
	}
	return &r
}
