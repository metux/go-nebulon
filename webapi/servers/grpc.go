package servers

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/webapi/blockstore/grpc_v1"
)

type GRPCServer struct {
	Conf      base.ServerConfig
	RPCServer *grpc.Server
}

func (server *GRPCServer) Shutdown(force bool) error {
	if force {
		server.RPCServer.Stop()
	} else {
		server.RPCServer.GracefulStop()
	}
	return nil
}

func (server *GRPCServer) Serve() error {
	if server.Conf.BlockStore == nil {
		return fmt.Errorf("no blockstore (%s) for grpc server [%w]", server.Conf.BlockStoreID, base.ErrNoStore)
	}

	server.RPCServer = grpc.NewServer()
	_ = grpc_v1.NewBlockService(server.RPCServer, server.Conf.BlockStore)

	l, err := net.Listen(server.Conf.Proto, ":"+server.Conf.Port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s [%w]", server.Conf.Port, err)
	}

	log.Printf("gRPC server listening at %v", l.Addr())
	if err := server.RPCServer.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}

	return nil
}

func NewGRPCServer(conf base.ServerConfig) *GRPCServer {
	grpc := GRPCServer{
		Conf: conf,
	}
	return &grpc
}
