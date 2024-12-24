package wire

import (
	"strings"
	"google.golang.org/protobuf/proto"
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
	return strings.Join(s, " ")
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
