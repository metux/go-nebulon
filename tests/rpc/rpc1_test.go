package rpc

import (
	"testing"

	"github.com/metux/go-nebulon/core/wire"
	"github.com/metux/go-nebulon/tests/common"
)

func run_rpc_test(t *testing.T, serverid string, storeid string) {
	store := setupConn(t, serverid, storeid)
	if err := store.Ping(); err != nil {
		t.Fatalf("rpc ping to %s failed: %s\n", storeid, err)
	}
	common.Test_store_load_peek_block(store, wire.RefType_Blob, []byte("hello world foo"))
}

func Test_GRPC_load_peek(t *testing.T) {
	run_rpc_test(t, "grpc-1", "grpc1")
}

func Test_HTTP_load_peek(t *testing.T) {
	run_rpc_test(t, "http-1", "http1")
}
