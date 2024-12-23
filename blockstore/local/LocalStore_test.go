package local

import (
	"testing"

	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/core/wire"
)

const (
	TestStorePath = "../.tmp/store.unit-test-1"
)

func Test_Put_Peek_1(t *testing.T) {
	bs, _ := NewByConfig(base.BlockStoreConfig{
		Type:   base.BlockStoreType_LocalFS,
		Url:    TestStorePath,
		Config: nil,
	}, nil)
	ref, err := bs.PutBlock([]byte(TestStorePath), wire.RefType_Blob)

	if err != nil {
		t.Fatalf("storing failed: %s\n", err)
	}

	t.Logf("Stored ref %s\n", ref.Dump())

	bi, err := bs.PeekBlock(ref, 1)
	if err != nil {
		t.Fatalf("Peek failed for %s: %s\n", ref.Dump(), err)
	}

	t.Logf("Peek result for %s: %+v\n", ref.Dump(), bi)
}
