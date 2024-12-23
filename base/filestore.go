package base

import (
	"io"
)

type FileStore interface {
    StoreFile (r io.Reader, headers map[string]string) (OID, error)
    ReadFile (oid OID) (io.Reader, map[string]string, error)
}
