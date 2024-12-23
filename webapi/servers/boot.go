package servers

import (
	"github.com/metux/go-nebulon/webapi/config"
)

func BootServer(conffile string, id string) (IServer, error) {
	srvcf, err := config.LoadServerConf(conffile)
	if err != nil {
		return nil, err
	}
	srvcf.InitStores()

	serverconf, err := srvcf.GetServer(id)
	if err != nil {
		return nil, err
	}

	server, err := NewServer(serverconf)
	if err != nil {
		return nil, err
	}

	return server, nil
}
