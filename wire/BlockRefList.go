package wire

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"strings"
)

func (rl BlockRefList) Dump() string {
	s := []string{}
	for _, walk := range rl.Refs {
		if walk == nil {
			s = append(s, "<nil>")
		} else {
			s = append(s, walk.HexKey())
		}
	}
	return fmt.Sprintf("(%d) %s", len(rl.Refs), strings.Join(s, " "))
}

func (rl BlockRefList) Count() int {
	return len(rl.Refs)
}

func (rl BlockRefList) Marshal() ([]byte, error) {
	// FIXME: perhaps we should sort fist
	if rl.Magic == "" {
		rl.Magic = "BLOCK REF LIST"
	}
	return proto.Marshal(&rl)
}

func (rl *BlockRefList) Unmarshal(data []byte) error {
	err := proto.Unmarshal(data, rl)
	return err
}

func (rl *BlockRefList) AddRef(ref BlockRef) {
	rl.Refs = append(rl.Refs, &ref)
}
