package wire

import (
	"testing"

	"google.golang.org/protobuf/proto"
)

func Test_BlockRef_sizes(t *testing.T) {
	ref := BlockRef{}
	bin, _ := proto.Marshal(&ref)
	t.Logf("Null ref: encoded size: %d\n", len(bin))

	ref = RefForBlock([]byte{}, RefType_Blob)
	bin, _ = proto.Marshal(&ref)
	t.Logf("Block ref: oid size: %d\n", len(ref.Oid))
	t.Logf("Block ref: encoded size: %d\n", len(bin))

	ref = RefForBlock([]byte{}, RefType_Blob)
	ref.Key = ref.Oid
	bin, _ = proto.Marshal(&ref)
	t.Logf("Block ref w/ key: encoded size: %d\n", len(bin))

	ref = RefForBlock([]byte{}, RefType_Blob)
	ref.Offset = 17
	ref.Limit = 229
	bin, _ = proto.Marshal(&ref)
	t.Logf("Block ref w/ range ref: encoded size: %d\n", len(bin))
}
