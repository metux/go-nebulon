package rpc

import (
	"testing"

	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/webapi/servers"

	"github.com/metux/go-nebulon/tests/common"
)

func setupConn(t *testing.T, serverid string, storeid string) base.IBlockStore {
	common.SetRoot("../../")
	server, _ := servers.BootServer("../../tests/perseus.conf.yaml", serverid)
	go server.Serve()
	store := common.TestStore(storeid)
	if err := store.Ping(); err != nil {
		t.Fatalf("rpc ping to %s failed: %s\n", storeid, err)
	}
	return store
}
