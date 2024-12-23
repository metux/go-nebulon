package dirtree

import (
	"testing"

	"github.com/metux/go-nebulon/filestore"
	"github.com/metux/go-nebulon/helpers"
	"github.com/metux/go-nebulon/tests/common"
	"github.com/metux/go-nebulon/util"
)

const (
	putdir_tempfile string = "../test-putdir.tmp"
	putdir_dirname  string = "../../"
)

func Test_PutDir_1(t *testing.T) {
	common.SetRoot("../../")

	t.Logf("Storing directory: %s\n", putdir_dirname)

	fs := filestore.NewFileStore(common.TestStore("unittest-helpers-1"))

	//	TracePutDirectory = true
	ref, err := helpers.PutDirectory(fs, "", putdir_dirname, util.FilterSkipHidden)

	if err != nil {
		t.Fatalf("storing failed: %s", err)
	}
	t.Logf("Stored dir ref %s", ref.Dump())

	if err := helpers.CompareTree(fs, putdir_dirname, ref, putdir_tempfile); err != nil {
		t.Fatalf("CompareTree() failed: %s", err)
	}
	t.Logf("CompareTree() done")
}
