package cache

import (
	"github.com/metux/go-nebulon/blockstore/common"
	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/core/wire"
)

const (
	StoreType = base.BlockStoreType_Cache
)

type Cache struct {
	common.StoreBase
	remote base.IBlockStore
	cache  base.IBlockStore
}

func (c Cache) PutBlock(data []byte, reftype wire.RefType) (base.BlockRef, error) {
	ref, err := c.cache.PutBlock(data, reftype)

	// push to remote in background
	go func() {
		remote_ref, remote_err := c.remote.PutBlock(data, reftype)
		if remote_err != nil {
			c.Logf("remote push error: %s\n", remote_err)
		} else {
			c.Logf("remote pushed: %s\n", remote_ref.Dump())
		}
	}()

	return ref, err
}

func (c Cache) GetBlock(ref base.BlockRef) ([]byte, error) {
	data, err := c.cache.GetBlock(ref)
	if err == nil {
		c.Logf("serving locally stored %s\n", ref.Dump())
		return data, err
	}

	return c.remote.GetBlock(ref)
}

func (c Cache) KeepBlock(ref base.BlockRef) error {
	err := c.cache.KeepBlock(ref)

	go func() {
		remote_err := c.remote.KeepBlock(ref)
		if remote_err != nil {
			c.Logf("remote keep err %s\n", remote_err)
		}
	}()

	return err
}

func (c Cache) IterateBlocks() base.BlockRefStream {
	// FIXME: should also query the remote ?
	return c.cache.IterateBlocks()
}

func (c Cache) DeleteBlock(ref base.BlockRef) error {
	return c.cache.DeleteBlock(ref)
}

func (c Cache) PeekBlock(ref base.BlockRef, fetch int) (base.BlockInfo, error) {
	inf, err := c.cache.PeekBlock(ref, fetch)
	if err == nil {
		return inf, err
	}
	return c.remote.PeekBlock(ref, fetch)
}

func (c Cache) Ping() error {
	return nil
}

func NewByConfig(config base.BlockStoreConfig, links map[string]base.IBlockStore) (*Cache, error) {
	val_remote, _ := links["remote"]
	val_cache, _ := links["cache"]

	store := Cache{
		StoreBase: common.StoreBase{
			Name: config.ID(),
		},
		remote: val_remote,
		cache:  val_cache,
	}

	if val_remote == nil {
		store.Logf("missing 'remote' link")
		return nil, base.ErrConfig
	}
	if val_cache == nil {
		store.Logf("missing 'cache' link")
		return nil, base.ErrConfig
	}

	return &store, nil
}
