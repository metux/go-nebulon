package wire

import (
	"google.golang.org/protobuf/proto"
)

func (fc FileControl) Marshal() ([]byte, error) {
	return proto.Marshal(&fc)
}

func (fc *FileControl) Unmarshal(data []byte) error {
	return proto.Unmarshal(data, fc)
}
