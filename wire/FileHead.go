package wire

import (
	"google.golang.org/protobuf/proto"
)

/*
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
*/

func (rl FileHead) Marshal() ([]byte, error) {
	// FIXME: perhaps we should sort fist
	if rl.Magic == "" {
		rl.Magic = "FILE HEAD"
	}
	return proto.Marshal(&rl)
}

func (rl *FileHead) Unmarshal(data []byte) error {
	err := proto.Unmarshal(data, rl)
	return err
}
