package wire

import (
	"google.golang.org/protobuf/proto"
)

func (rl FileHead) Marshal() ([]byte, error) {
	if rl.Magic == "" {
		rl.Magic = "FILE HEAD"
	}
	return proto.Marshal(&rl)
}

func FileHeadUnmarshal(data []byte) (FileHead, error) {
	fh := FileHead{}
	err := proto.Unmarshal(data, &fh)
	return fh, err
}
