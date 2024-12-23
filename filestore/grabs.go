package filestore

import (
	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/core/crypt"
	"github.com/metux/go-nebulon/core/wire"
)

type GrabReader struct {
	ch         base.BlockRefStream
	BlockStore base.IBlockStore
}

func (r GrabReader) traceRefPtr(ref *base.BlockRef) {
	if ref != nil {
		r.traceRef(*ref)
	}
}

func (r GrabReader) traceRef(ref base.BlockRef) {
	//	log.Printf("traceRef: ref=%s\n", ref.Dump())
	switch ref.Type {
	case wire.RefType_Blob:
		// just send the ref, nothing more to do here
		r.ch.SendRef(ref)
	case wire.RefType_RefList:
		// need to parse it and also send it's entries
		bl, err := crypt.BlockListLoadDecrypt(r.BlockStore, ref)
		if err == crypt.ErrNoCryptoKey {
			// perfectly normal situation: encrypted BlockRefList
			// there content is stripped off keys and added as
			// separate entries in the grab list. just treat it
			// like a blob
			r.ch.SendRef(ref)
			return
		}
		r.ch.SendRefErr(ref, err)
		for _, walk := range bl.Refs {
			r.traceRefPtr(walk)
		}
	case wire.RefType_File, wire.RefType_Directory:
		fh, err := crypt.LoadFileHead(r.BlockStore, ref)
		r.ch.SendRefErr(ref, err)
		r.traceRefPtr(fh.Grabs)
	default:
		// how bad: we don't know whether we have to parse it
		// likely not ... but, who knows ?
		r.ch.SendRefErr(ref, ref.UnsupportedTypeError())
	}
}
