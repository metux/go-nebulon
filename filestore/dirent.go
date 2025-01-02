package filestore

import (
	"github.com/metux/go-nebulon/wire"
)

type DirEntry struct {
	Ref wire.BlockRef
}

func (dent DirEntry) IsDir() bool {
	return dent.Ref.Type == wire.RefType_Directory
}

func (dent DirEntry) IsFile() bool {
	return dent.Ref.Type == wire.RefType_File
}

func (dent DirEntry) Name() string {
	return dent.Ref.Name
}

func NewDirEntry(ref wire.BlockRef) DirEntry {
	return DirEntry{Ref: ref}
}
