package wire

import (
	"google.golang.org/protobuf/proto"
)

func (frame AnnounceFrame) Marshal() ([]byte, error) {
	return proto.Marshal(&frame)
}

func AnnounceFrameUnmarshal(data []byte) (AnnounceFrame, error) {
	// needs to be in this order
	frame := AnnounceFrame{}
	err := proto.Unmarshal(data, &frame)
	return frame, err
}

func (payload AnnouncePayload) Marshal() ([]byte, error) {
	return proto.Marshal(&payload)
}

func AnnouncePayloadUnmarshal(data []byte) (AnnouncePayload, error) {
	// needs to be in this order
	payload := AnnouncePayload{}
	err := proto.Unmarshal(data, &payload)
	return payload, err
}
