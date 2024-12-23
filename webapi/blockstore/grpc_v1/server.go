package grpc_v1

import (
	"context"

	grpc "google.golang.org/grpc"
	grpc_codes "google.golang.org/grpc/codes"
	grpc_status "google.golang.org/grpc/status"

	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/core/wire"
)

type BlockService struct {
	wire.UnimplementedBlockStoreServer
	BlockStore base.IBlockStore
}

func (s *BlockService) Put(ctx context.Context, req *wire.RPC_PutBlockRequest) (*wire.RPC_PutBlockReply, error) {
	if req == nil {
		return nil, grpc_status.Errorf(grpc_codes.Unimplemented, "req may not be nil")
	}

	ref, err := s.BlockStore.PutBlock(req.Data, req.Reftype)
	if err != nil {
		return nil, err
	}

	return &wire.RPC_PutBlockReply{Ref: &ref}, nil
}

func (s *BlockService) Get(ctx context.Context, req *wire.RPC_GetBlockRequest) (*wire.RPC_GetBlockReply, error) {
	if req == nil {
		return nil, grpc_status.Errorf(grpc_codes.Unimplemented, "req may not be nil")
	}

	data, err := s.BlockStore.GetBlock(*req.Ref)
	if err != nil {
		return nil, err
	}

	return &wire.RPC_GetBlockReply{Data: data}, nil
}

func (s *BlockService) Peek(ctx context.Context, req *wire.RPC_PeekBlockRequest) (*wire.RPC_PeekBlockReply, error) {
	if req == nil {
		return nil, grpc_status.Errorf(grpc_codes.Unimplemented, "req may not be nil")
	}

	return nil, grpc_status.Errorf(grpc_codes.Unimplemented, "PeekBlock not implemented")
}

func (s *BlockService) Keep(ctx context.Context, req *wire.RPC_KeepBlockRequest) (*wire.RPC_KeepBlockReply, error) {
	if req == nil {
		return nil, grpc_status.Errorf(grpc_codes.Unimplemented, "req may not be nil")
	}

	return nil, grpc_status.Errorf(grpc_codes.Unimplemented, "KeepBlock not implemented")
}

func (s *BlockService) Ping(context.Context, *wire.RPC_PingRequest) (*wire.RPC_PingReply, error) {
	return &wire.RPC_PingReply{Msg: "ACK"}, nil
}

func NewBlockService(reg grpc.ServiceRegistrar, blockstore base.IBlockStore) *BlockService {
	if blockstore == nil {
		panic("NewBlockService: blockstore = nil")
	}

	srv := &BlockService{
		BlockStore: blockstore,
	}
	wire.RegisterBlockStoreServer(reg, srv)
	return srv
}
