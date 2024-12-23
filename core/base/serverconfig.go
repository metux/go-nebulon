package base

type ServerConfig struct {
	Name         string `yaml:"name"`
	Type         string `yaml:"type"`
	Proto        string `yaml:"proto"`
	Port         string `yaml:"port"`
	BlockStoreID string `yaml:"blockstore"`

	BlockStore IBlockStore `yaml:"-"`
}

type ServerMap map[string]ServerConfig
