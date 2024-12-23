package base

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound          = errors.New("block not found")
	ErrTimeout           = errors.New("i/o timeout")
	ErrOffline           = errors.New("store offline")
	ErrConfig            = errors.New("store configuration error")
	ErrUnsupportedStore  = fmt.Errorf("unsupported store type [%w]", ErrConfig)
	ErrUnsupportedServer = fmt.Errorf("unsupported server type [%w]", ErrConfig)
	ErrNoStore           = fmt.Errorf("no such store [%w]", ErrConfig)
	ErrNotImplemented    = errors.New("not implemented")
	ErrMissingUrl        = fmt.Errorf("missing URL [%w]", ErrConfig)
	ErrInvalidOID        = errors.New("invalid oid")
	ErrInvalidKey        = errors.New("invalid key")
)
