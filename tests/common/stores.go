package common

import (
	"fmt"
	"log"

	"github.com/metux/go-nebulon/blockstore"
	"github.com/metux/go-nebulon/blockstore/local"
	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/webapi/config"
)

var (
	rootDir string = "."

	loaded  bool                    = false
	conf    blockstore.StoreConf    = blockstore.StoreConf{}
	srvconf config.ServerConfigFile = config.ServerConfigFile{}
)

func SetRoot(s string) {
	rootDir = s
	local.VirtualRoot = s
}

func load() {
	PanicX("config load error", LoadConfig(rootDir+"/tests/perseus.conf.yaml"))
}

func LoadConfig(fn string) error {
	log.Printf("loading config: %s\n", fn)
	if !loaded {
		if err := conf.Load(fn); err != nil {
			log.Printf("Loader error: %s\n", err)
			return err
		}
	}
	return nil
}

func TestStore(name string) base.IBlockStore {
	load()
	bs := conf.GetStore(name)
	if bs == nil {
		panic(fmt.Errorf("TestStore: cant get store %s\n", name))
	}
	return bs
}

func PanicX(prefix string, err error) {
	if err != nil {
		panic(fmt.Errorf(prefix+" [%w]", err))
	}
}
