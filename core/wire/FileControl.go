package wire

import (
	"google.golang.org/protobuf/proto"
)

func (fc FileControl) Marshal() ([]byte, error) {
	return proto.Marshal(&fc)
}

func FileControlUnmarshal(data []byte) (FileControl, error) {
	// needs to be in this order
	fc := FileControl{}
	err := proto.Unmarshal(data, &fc)
	return fc, err
}
