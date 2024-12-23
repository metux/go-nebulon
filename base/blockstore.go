package base

type BlockStore interface {
	StoreBlock([]byte) (OID, error)
	LoadBlock(OID) ([]byte, error)
}
