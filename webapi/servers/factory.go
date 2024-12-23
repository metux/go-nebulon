package servers

import (
	"fmt"

	"github.com/metux/go-nebulon/core/base"
)

func NewServer(conf base.ServerConfig) (IServer, error) {
	switch conf.Type {
	case "grpc":
		return NewGRPCServer(conf), nil
	case "http":
		return NewHTTPServer(conf), nil
	default:
		return nil, fmt.Errorf("unsupported server type %s [%w]", conf.Type, base.ErrUnsupportedServer)
	}
}
