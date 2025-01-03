package filestore

import (
	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/blockcrypt"
	"github.com/metux/go-nebulon/wire"
)

type readerBase struct {
	// the underlying BlockStore to write into
	BlockStore base.BlockStore
}

func (r readerBase) loadBlockList(ref wire.BlockRef) (wire.BlockRefList, error) {
	reflist := wire.BlockRefList{}

	data, err := r.loadBlock(ref)

	// note do it in separate steps, since reflist is changed here
	err = reflist.Unmarshal(data)
	return reflist, err
}

func (r readerBase) loadBlock(ref wire.BlockRef) ([]byte, error) {
	return blockcrypt.BlockLoadDecrypt(r.BlockStore, ref)
}

func (r *readerBase) loadFileControl(ref wire.BlockRef) (wire.FileControl, error) {
	return blockcrypt.LoadFileControl(r.BlockStore, ref)
}
