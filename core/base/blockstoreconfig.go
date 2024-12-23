package base

const (
	BlockStoreType_LocalFS = "local"
	BlockStoreType_HTTP    = "http"
	BlockStoreType_GRPC    = "grpc"
	BlockStoreType_Cache   = "cache"
)

type BlockStoreConfig struct {
	Name   string `yaml:"-"`
	Type   string
	Url    string
	Config map[string]string
	Links  map[string]string
	Store  IBlockStore `yaml:"-"`
}

func (bsc BlockStoreConfig) ID() string {
	return bsc.Type + " " + bsc.Name
}
