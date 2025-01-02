package helpers

import (
	"testing"

	"github.com/udhos/equalfile"

	"github.com/metux/go-nebulon/blockstore"
	"github.com/metux/go-nebulon/filestore"
	"github.com/metux/go-nebulon/util"
	"github.com/metux/go-nebulon/wire"
)

const (
	//	filename string = "/home/nekrad/dl/000.capture/elen0_tg/elen_cross-2024-09-04-04-26-03.P00.mkv.00.mux.mp4.tg.mp4"
	filename string = "../go-nebulon"
	tempfile string = "../test1.tmp"
	dirname  string = ".."
)

var (
	fs = filestore.NewFileStore(blockstore.NewSimpleStore("../.storedata"))
)

func Test_PutGet_1(t *testing.T) {
	t.Logf("Storing file: %s\n", filename)
	ref, err := PutFile(fs,
		"test1",
		wire.Header{wire.Header_ContentType: wire.ContentType_MP4},
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

func Test_PutDir_1(t *testing.T) {
	t.Logf("Storing directory: %s\n", dirname)

	//	TracePutDirectory = true
	ref, err := PutDirectory(fs, "", dirname, util.FilterSkipHidden)

	if err != nil {
		t.Fatalf("storing failed: %s", err)
	}
	t.Logf("Stored dir ref %s", ref.Dump())

	if err := CompareTree(fs, dirname, ref, tempfile); err != nil {
		t.Fatalf("CompareTree() failed: %s", err)
	}
	t.Logf("CompareTree() done")
}
