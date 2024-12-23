package blockstore

import (
	"log"
	"os"

	"github.com/metux/go-nebulon/core/base"
	"gopkg.in/yaml.v2"
)

type StoreConf struct {
	Stores BlockStoreMap `yaml:"blockstore"`
}

func (sc *StoreConf) Load(fn string) error {
	yamlFile, err := os.ReadFile(fn)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, sc)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return err
	}
	return sc.Stores.Init()
}

func (sc StoreConf) GetStore(id string) base.IBlockStore {
	return sc.Stores.GetStore(id)
}
