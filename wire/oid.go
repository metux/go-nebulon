package wire

import (
	"github.com/metux/go-nebulon/base"
	"google.golang.org/protobuf/proto"
)

func EncapOID(ref base.BlockRef) * BlockRef {
	return &BlockRef{Data: ref.OID.Data, Type: uint32(ref.Type)}
}

func DecapOID(oid BlockRef) base.OID {
	return base.OID{Data: oid.Data}
}

func EncapOIDRefList(oids []base.BlockRef) OIDRefList {
	l := make([]*BlockRef, len(oids))
	for idx,ent := range oids {
		l[idx] = EncapOID(ent)
	}
	return OIDRefList{Magic: "BLOCK REF LIST", Refs: l}
}

func MarshalOIDRefList(oids []base.BlockRef) ([]byte, error) {
	reflist := EncapOIDRefList(oids)
	return proto.Marshal(&reflist)
}
