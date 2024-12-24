package base

// these must be equal to the ones in wire proto

type RefType uint32

const (
	RefType_Blob    = RefType(0)
	RefType_RefList = RefType(1)
)

type BlockRef struct {
	OID  OID
	Type RefType
	Key  []byte
}
