package blockstore

import (
	"log"

	"github.com/metux/go-nebulon/core/base"
)

type BlockStoreMap map[string](*base.BlockStoreConfig)

func (sc BlockStoreMap) GetStore(id string) base.IBlockStore {
	val, _ := sc[id]
	if val != nil {
		return val.Store
	}
	return nil
}

func (sm BlockStoreMap) initOne(st *base.BlockStoreConfig) bool {
	// walk through the link list and check which ones are already there
	links := map[string]base.IBlockStore{}
	for linkname, linktarget := range st.Links {
		if link_st, ok := sm[linktarget]; ok {
			if link_st.Store == nil {
				return false
			} else {
				links[linkname] = link_st.Store
			}
		} else {
			log.Printf("link target in %s not defined: %s\n", st.Name, linktarget)
			return false
		}
	}

	newst, err := NewStoreByConfig(*st, links)
	if err != nil {
		log.Printf("error creating new store %s: %s\n", st.Name, err)
		return false
	}
	st.Store = newst
	return true
}

func (sm BlockStoreMap) initPass() bool {
	didsome := false
	for _, st := range sm {
		if st.Store != nil {
			continue
		}
		didsome = didsome || sm.initOne(st)
	}
	return didsome
}

func (sm BlockStoreMap) Init() error {
	// some post-load fixups
	for k, v := range sm {
		v.Name = k
	}

	for sm.initPass() {
		// just loop through
	}
	return nil
}
