package base

import (
	"io"
)

type BlockRefStreamEnt struct {
	Ref      BlockRef
	Error    error
	Finished bool
}

type BlockRefStream chan BlockRefStreamEnt

func (ch BlockRefStream) SendRef(ref BlockRef) {
	ch <- BlockRefStreamEnt{Ref: ref}
}

func (ch BlockRefStream) SendError(err error) {
	if err != nil && err != io.EOF {
		ch <- BlockRefStreamEnt{Error: err}
	}
}

func (ch BlockRefStream) SendRefErr(ref BlockRef, err error) {
	ch.SendRef(ref)
	ch.SendError(err)
}

func (ch BlockRefStream) Finish() {
	ch <- BlockRefStreamEnt{Finished: true}
	close(ch)
}
