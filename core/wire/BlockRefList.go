package wire

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"strings"

	"google.golang.org/protobuf/proto"
)

func (rl BlockRefList) Dump() string {
	s := []string{}
	for _, walk := range rl.Refs {
		if walk == nil {
			s = append(s, "<nil>")
		} else {
			s = append(s, walk.Dump())
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

func (b BlockRefList) Len() int {
	return len(b.Refs)
}

func (b BlockRefList) Less(i, j int) bool {
	// bytes package already implements Comparable for []byte.
	switch bytes.Compare(b.Refs[i].Oid, b.Refs[j].Oid) {
	case -1:
		return true
	case 0, 1:
		return false
	default:
		log.Panic("not fail-able with `bytes.Comparable` bounded [-1, 1].")
		return false
	}
}

func (b BlockRefList) Swap(i, j int) {
	b.Refs[j], b.Refs[i] = b.Refs[i], b.Refs[j]
}

func (b BlockRefList) Sort() {
	sort.Sort(b)
}

// note: ref parameter *must* be value, not pointer, so the function is
// actually creating a copy
func (b *BlockRefList) Add(ref BlockRef) {
	b.Refs = append(b.Refs, &ref)
}
