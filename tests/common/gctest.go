package common

import (
	"log"

	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/helpers"
)

func testGC(bs base.IBlockStore, r base.BlockRef) {
	gc := helpers.MakeSimpleGC(bs)
	count := gc.Mark(r)
	log.Printf("grab count: %d\n", count)

	kick, kept := gc.ScanKick()

	log.Printf("Kept: %d\n", kept)
	log.Printf("Kick: %d\n", len(kick))

	for _, k := range kick {
		err := bs.DeleteBlock(k)
		if err == nil {
			log.Printf("DELETED %s\n", k.Dump())
		} else {
			log.Printf("ERROR   %s -- %s\n", k.Dump(), err)
		}
	}
}
