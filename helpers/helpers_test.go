package helpers

import (
	"testing"

	"github.com/udhos/equalfile"

	"github.com/metux/go-nebulon/blockstore"
	"github.com/metux/go-nebulon/filestore"
)

const (
//	filename string = "/home/nekrad/dl/000.capture/elen0_tg/elen_cross-2024-09-04-04-26-03.P00.mkv.00.mux.mp4.tg.mp4"
	filename string = "../go-nebulon"
	tempfile string = "../test1.tmp"
)

var (
	fs = filestore.NewFileStore(blockstore.NewSimpleStore("../.storedata"))
)

func Test_PutGet_1(t *testing.T) {
	t.Logf("Storing file: %s\n", filename)
	ref, err := StoreFile(fs,
		"test1",
		map[string]string{"Content-Type": "video/mp4"},
		filename)

	if err != nil {
		t.Fatalf("storing failed: %e", err)
	}
	t.Logf("Stored file ref %s", ref.Dump())

	headers, err := GetFile(fs, tempfile, ref)
	if err != nil {
		t.Fatalf("GetFile() failed: %s", err)
	}
	t.Logf("Pulled file: headers=%s", headers)

	cmp := equalfile.New(nil, equalfile.Options{}) // compare using single mode
	equal, err := cmp.CompareFile(filename, tempfile)

	if err != nil {
		t.Fatalf("file compare failed: %s", err)
	}

	if !equal {
		t.Fatal("files dont match")
	}

	t.Logf("files matching")
}
