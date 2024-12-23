package dirtree

import (
	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/filestore"
	"github.com/metux/go-nebulon/helpers"
	"github.com/metux/go-nebulon/tests/common"
	"github.com/metux/go-nebulon/util"
)

func StoreDirTree(fs base.IFileStore) base.BlockRef {
	ref, err := helpers.PutDirectory(fs, "", ".", util.FilterSkipHidden)
	common.PanicX("storeDirTree", err)
	return ref
}

func CompareDirTree(bs base.IBlockStore, ref base.BlockRef, fn string) base.BlockRef {
	fs := filestore.NewFileStore(bs)
	common.PanicX("compareDirTree", helpers.CompareTree(fs, ".", ref, fn))
	return ref
}

func RunDirTree(bs base.IBlockStore, fn string) base.BlockRef {
	fs := filestore.NewFileStore(bs)
	return CompareDirTree(bs, StoreDirTree(fs), fn)
}
