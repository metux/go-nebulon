package helpers

import (
	"log"

	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/filestore"
)

type SimpleGC struct {
	BlockStore base.IBlockStore
	Grabs      map[string]base.BlockRefStreamEnt
}

func MakeSimpleGC(bs base.IBlockStore) SimpleGC {
	return SimpleGC{
		BlockStore: bs,
		Grabs:      make(map[string]base.BlockRefStreamEnt),
	}
}

func (gc *SimpleGC) Mark(ref base.BlockRef) int {
	fs := filestore.NewFileStore(gc.BlockStore)
	ch := fs.Grabs(ref)
	count := 0
	for walk := range ch {
		oid := walk.Ref.OID()
		count++
		if me, ok := gc.Grabs[oid]; ok {
			log.Printf(" --> duplicate: %s\n", me.Ref.Dump())
		} else {
			if walk.Error != nil {
				log.Printf(" --> received error %s\n", walk.Error)
			} else if !walk.Finished {
				gc.Grabs[oid] = walk
			}
		}
	}
	return count
}

// we should iterate over recently old blocks
func (gc *SimpleGC) ScanKick() ([]base.BlockRef, int) {
	kick := []base.BlockRef{}
	kept := 0
	for ent := range gc.BlockStore.IterateBlocks() {
		oid := ent.Ref.OID()
		id := ent.Ref.Dump()
		if _, ok := gc.Grabs[oid]; ok {
			log.Printf("Block: KEEP %s\n", id)
			kept++
		} else {
			log.Printf("Block: KICK %s\n", id)
			kick = append(kick, ent.Ref)
		}
	}
	return kick, kept
}
