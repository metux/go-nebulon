package wireutil

import (
	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/wire"
	"google.golang.org/protobuf/proto"
)

func EncapOID(ref base.BlockRef) * wire.BlockRef {
	return &wire.BlockRef{Data: ref.OID.Data, Type: uint32(ref.Type)}
}

func DecapOID(oid wire.BlockRef) base.OID {
	return base.OID{Data: oid.Data}
}

func EncapOIDRefList(oids []base.BlockRef) wire.OIDRefList {
	l := make([]*wire.BlockRef, len(oids))
	for idx,ent := range oids {
		l[idx] = &wire.BlockRef{Data: ent.OID.Data, Type: uint32(ent.Type)}
	}
	return wire.OIDRefList{Magic: "BLOCK REF LIST", Refs: l}
}

func MarshalOIDRefList(oids []base.BlockRef) ([]byte, error) {
	reflist := EncapOIDRefList(oids)
	return proto.Marshal(&reflist)
}
