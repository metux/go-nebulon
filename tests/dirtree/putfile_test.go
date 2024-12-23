package dirtree

import (
	"testing"

	"github.com/udhos/equalfile"

	"github.com/metux/go-nebulon/core/wire"
	"github.com/metux/go-nebulon/filestore"
	"github.com/metux/go-nebulon/helpers"
	"github.com/metux/go-nebulon/tests/common"
)

const (
	//	putfile_filename string = "../../go.sum"
	putfile_tempfile string = "../test1.tmp"
	//	putfile_filename string = "/home/nekrad/dl/heilung_lifa.mp4"
	putfile_filename string = "/home/nekrad/dl/ed_hd.mp4"
)

func Test_PutGet_1(t *testing.T) {
	common.SetRoot("../../")
	fs := filestore.NewFileStore(common.TestStore("local1"))

	t.Logf("Storing file: %s\n", putfile_filename)
	ref, err := helpers.PutFile(fs,
		"test1",
		wire.Header{wire.Header_ContentType: wire.ContentType_MP4},
		putfile_filename)

	if err != nil {
		t.Fatalf("storing failed: %e", err)
	}
	t.Logf("Stored file ref %s", ref.Dump())

	headers, err := helpers.GetFile(fs, putfile_tempfile, ref)
	if err != nil {
		t.Fatalf("GetFile() failed: %s", err)
	}
	t.Logf("Pulled file: headers=%s", headers)

	cmp := equalfile.New(nil, equalfile.Options{}) // compare using single mode
	equal, err := cmp.CompareFile(putfile_filename, putfile_tempfile)

	if err != nil {
		t.Fatalf("file compare failed: %s", err)
	}

	if !equal {
		t.Fatal("files dont match")
	}

	t.Logf("files matching")
}
