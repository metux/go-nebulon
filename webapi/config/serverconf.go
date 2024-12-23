package config

import (
	"fmt"

	"github.com/metux/go-nebulon/blockstore"
	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/util"
)

type ServerConfigFile struct {
	Servers base.ServerMap           `yaml:"server"`
	Stores  blockstore.BlockStoreMap `yaml:"blockstore"`
}

func (sc *ServerConfigFile) Load(fn string) error {
	if err := util.LoadYaml(fn, sc); err != nil {
		return err
	}
	return nil
}

func (sc *ServerConfigFile) InitStores() error {
	return sc.Stores.Init()
}

func (sc *ServerConfigFile) GetServer(id string) (base.ServerConfig, error) {
	cf, ok := sc.Servers[id]
	if !ok {
		return cf, fmt.Errorf("no server config %s", id)
	}
	cf.BlockStore = sc.Stores.GetStore(cf.BlockStoreID)
	if cf.BlockStore == nil {
		return cf, fmt.Errorf("missing blockstore config %s", cf.BlockStoreID)
	}
	return cf, nil
}

func LoadServerConf(fn string) (ServerConfigFile, error) {
	cf := ServerConfigFile{}
	err := cf.Load(fn)
	return cf, err
}
