package wire

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/metux/go-nebulon/util"
)

func ParseBlockRef(str string) BlockRef {
	elems := strings.Split(str, ":")
	ref := BlockRef{}

	if len(elems) > 0 {
		v, _ := RefType_value[elems[0]]
		ref.Type = RefType(v)
		if len(elems) > 1 {
			ref.Oid, _ = hex.DecodeString(elems[1])
			if len(elems) > 2 {
				ref.Key, _ = hex.DecodeString(elems[2])
			}
		}
	}
	return ref
}

func RefForBlock(data []byte, t RefType) BlockRef {
	return BlockRef{
		Type: t,
		Oid:  util.ContentKey(data),
	}
}

func (ref BlockRef) Dump() string {
	s := fmt.Sprintf("%s:%X:%X", ref.Type, ref.Oid, ref.Key)
	if ref.Name != "" {
		s = s + " <" + ref.Name + ">"
	}
	return s
}

func (ref BlockRef) OID() string {
	return fmt.Sprintf("%X", ref.Oid)
}

func (ref BlockRef) IsDir() bool {
	return ref.Type == RefType_Directory
}

func (ref BlockRef) IsFile() bool {
	return ref.Type == RefType_File
}

// convert a normal block ref to a grab ref:
// kick out the key and other unncessary data
// fix: we might need to keep the encryption type, in case of unencrypted
// blocks are compressed
func (ref BlockRef) ToGrab() BlockRef {
	switch ref.Type {
	case RefType_RefList:
		// needs to be rewritten to Blob, since encrypted
		//		ref.Type = RefType_Blob
	case RefType_Blob, RefType_File, RefType_Directory:
		// leave them as they are
	default:
		// unhandled
	}

	to := BlockRef{
		Oid:    ref.Oid,
		Type:   ref.Type,
		Cipher: ref.Cipher,
	}

	return to
}

func (ref BlockRef) UnsupportedTypeError() error {
	return fmt.Errorf("unsupported ref type %s", ref.Type)
}
