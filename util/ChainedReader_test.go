package util

import (
	"io"
	"strings"
	"testing"
)

func Test_ChainedReader_1(t *testing.T) {
	want := "hello world"

	chained := NewChainedReader(
		strings.NewReader("hello"),
		strings.NewReader(" "),
		strings.NewReader("world"),
	)

	ret, err := io.ReadAll(chained)
	if err != nil {
		t.Errorf("io.ReadAll() error %s\n", err)
	}

	got := string(ret)

	if got != want {
		t.Errorf("expected \"%s\" got=\"%s\"\n", want, got)
	}
}
