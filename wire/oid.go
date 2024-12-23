package wire

import (
	"github.com/metux/go-nebulon/base"
	"google.golang.org/protobuf/proto"
)

func EncapOID(oid base.OID) * OID {
	return &OID{Data: oid.Data}
}

func DecapOID(oid OID) base.OID {
	return base.OID{Data: oid.Data}
}

func EncapOIDRefList(oids []base.OID) OIDRefList {
	l := make([]*OID, len(oids))
	for idx,ent := range oids {
		l[idx] = EncapOID(ent)
	}
	return OIDRefList{Oids: l}
}

func MarshalOIDRefList(oids []base.OID) ([]byte, error) {
	reflist := EncapOIDRefList(oids)
	return proto.Marshal(&reflist)
}
