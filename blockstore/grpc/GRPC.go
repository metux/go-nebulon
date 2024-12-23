package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/metux/go-nebulon/blockstore/common"
	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/core/wire"
)

type GRPC struct {
	common.StoreBase
	Url    string
	client wire.BlockStoreClient
}

func (rpc *GRPC) Connect() error {
	if rpc.client == nil {
		conn, err := grpc.NewClient(rpc.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			rpc.Logf("failed to connect to gRPC server at localhost:50051: %v", err)
			return err
		}
		rpc.client = wire.NewBlockStoreClient(conn)
		rpc.Logf("Connected GRPC remote: %s\n", rpc.Url)
	}
	return nil
}

func (rpc *GRPC) PutBlock(data []byte, reftype wire.RefType) (base.BlockRef, error) {
	ref := wire.RefForBlock(data, reftype)

	if err := rpc.Connect(); err != nil {
		return ref, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := rpc.client.Put(ctx, &wire.RPC_PutBlockRequest{Data: data, Reftype: reftype})
	if err != nil {
		rpc.Logf("error calling function PutBlock: %v", err)
		return ref, err
	}

	return *r.Ref, nil
}

func (rpc *GRPC) GetBlock(ref base.BlockRef) ([]byte, error) {
	if err := rpc.Connect(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := rpc.client.Get(ctx, &wire.RPC_GetBlockRequest{Ref: &ref})
	if err != nil {
		rpc.Logf("error calling function GetBlock: %v", err)
		return nil, err
	}

	return r.Data, nil
}

func (rpc *GRPC) KeepBlock(ref base.BlockRef) error {
	if err := rpc.Connect(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if _, err := rpc.client.Keep(ctx, &wire.RPC_KeepBlockRequest{Ref: &ref}); err != nil {
		rpc.Logf("error calling function Keep: %v", err)
		return err
	}

	return nil
}

// FIXME: not implemented yet
func (rpc *GRPC) IterateBlocks() base.BlockRefStream {
	ch := make(base.BlockRefStream, IterateChanSize)

	rpc.Logf("IterateBlocks() not implemented yet\n")
	go func() {
		ch.Finish()
	}()

	return ch
}

func (rpc *GRPC) DeleteBlock(ref base.BlockRef) error {
	if err := rpc.Connect(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if _, err := rpc.client.Delete(ctx, &wire.RPC_DeleteBlockRequest{Ref: &ref}); err != nil {
		rpc.Logf("error calling function Delete: %v", err)
		return err
	}

	return nil
}

// FIXME: not implemented yet
func (rpc *GRPC) PeekBlock(ref base.BlockRef, fetch int) (base.BlockInfo, error) {
	return base.BlockInfo{Ref: ref}, nil
}

func (rpc *GRPC) Ping() error {
	if err := rpc.Connect(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := rpc.client.Ping(ctx, &wire.RPC_PingRequest{Msg: "hello there"})
	if err != nil {
		rpc.Logf("error calling function Ping: %v", err)
		return err
	}

	rpc.Logf("ping reply: %s\n", r.Msg)
	return nil
}

func NewByConfig(config base.BlockStoreConfig, links map[string]base.IBlockStore) (*GRPC, error) {
	if config.Url == "" {
		return nil, base.ErrMissingUrl
	}

	return &GRPC{
		StoreBase: common.StoreBase{
			Name: config.ID(),
		},
		Url: config.Url,
	}, nil
}
