package local

import (
	"encoding/hex"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/metux/go-nebulon/blockstore/common"
	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/core/wire"
)

type LocalFS struct {
	common.StoreBase
	Path  string
	Error error
}

func (s LocalFS) FilenameForRef(ref base.BlockRef) string {
	return s.Path + "/objects/" + ref.Type.String() + "/" + ref.OID()
}

func (s LocalFS) PutBlock(data []byte, reftype wire.RefType) (base.BlockRef, error) {
	ref := wire.RefForBlock(data, reftype)
	fn := s.FilenameForRef(ref)
	os.MkdirAll(filepath.Dir(fn), os.ModePerm)

	// dont write already existing objects
	// 2do: we could check equality here
	// 2do: should touch() it
	if _, err := os.Stat(fn); err == nil {
		if TraceWrite {
			log.Printf("object already exists %s\n", fn)
		}
		return ref, nil
	}
	// fixme: should write to temp and later rename
	err := os.WriteFile(fn, data, 0644)
	return ref, err
}

func (s LocalFS) xlateErr(err error) error {
	if errors.Is(err, os.ErrNotExist) {
		return base.ErrNotFound
	}
	if errors.Is(err, os.ErrInvalid) {
		return base.ErrNotFound
	}
	if errors.Is(err, os.ErrClosed) {
		return base.ErrNotFound
	}
	if errors.Is(err, os.ErrDeadlineExceeded) {
		return base.ErrTimeout
	}
	return err
}

func (s LocalFS) GetBlock(ref base.BlockRef) ([]byte, error) {
	data, err := os.ReadFile(s.FilenameForRef(ref))
	return data, s.xlateErr(err)
}

func (s LocalFS) KeepBlock(ref base.BlockRef) error {
	fileName := s.FilenameForRef(ref)
	_, err := os.Stat(fileName)
	if err != nil {
		return s.xlateErr(err)
	}

	currentTime := time.Now().Local()
	return s.xlateErr(os.Chtimes(fileName, currentTime, currentTime))
}

func (s LocalFS) IterateBlocks() base.BlockRefStream {
	ch := make(base.BlockRefStream, IterateChanSize)

	go func() {
		for reftype, subdir := range wire.RefType_name {
			dirp, _ := os.Open(s.Path + "/objects/" + subdir)
			if dirp == nil {
				continue
			}
			defer dirp.Close()
			entries, err := dirp.Readdir(IterateDirSize)
			for err == nil {
				for _, ent := range entries {
					oid, err := hex.DecodeString(ent.Name())
					if err == nil {
						ref := base.BlockRef{
							Oid:  oid,
							Type: wire.RefType(reftype),
						}
						ch.SendRef(ref)
					} else {
						ch.SendError(err)
					}
				}
				entries, err = dirp.Readdir(IterateDirSize)
			}
			ch.SendError(err)
		}
		ch.Finish()
	}()

	return ch
}

func (s LocalFS) DeleteBlock(ref base.BlockRef) error {
	return s.xlateErr(os.Remove(s.FilenameForRef(ref)))
}

func (s LocalFS) PeekBlock(ref base.BlockRef, fetch int) (base.BlockInfo, error) {
	fn := s.FilenameForRef(ref)
	inf := base.BlockInfo{Ref: ref}

	st, err := os.Stat(fn)
	if err != nil {
		return inf, s.xlateErr(err)
	}

	inf.Size = st.Size()
	inf.Present = st.Size()
	inf.ModTime = st.ModTime()

	return inf, nil
}

func (s LocalFS) Ping() error {
	return nil
}

func NewByConfig(config base.BlockStoreConfig, links map[string]base.IBlockStore) (*LocalFS, error) {
	if config.Url == "" {
		return nil, base.ErrMissingUrl
	}
	return &LocalFS{
		StoreBase: common.StoreBase{
			Name: config.ID(),
		},
		Path: strings.Replace(config.Url, "${ROOT}", VirtualRoot, -1),
	}, nil
}
