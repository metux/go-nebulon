package common

import (
	"bytes"
	"fmt"
	"log"

	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/core/wire"
)

func Test_store_load_peek_block(bs base.IBlockStore, reftype wire.RefType, data []byte) {
	PanicX("ping error", bs.Ping())
	ref, err := bs.PutBlock(data, reftype)
	PanicX("Peek test: put failed", err)
	log.Printf("Peek test: stored %s\n", ref.Dump())
	data2, err := bs.GetBlock(ref)
	PanicX("Peek test: get failed", err)
	if bytes.Compare(data, data2) != 0 {
		PanicX("Peek test", fmt.Errorf("block bytes dont match"))
	}
	log.Printf("Peek test: compare ok")
}
