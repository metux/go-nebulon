package dirtree

import (
	"errors"
	"testing"

	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/filestore"
	"github.com/metux/go-nebulon/tests/common"
)

const (
	keep_test_temp_file = ".download.tmp"
)

// FIXME: move this to test cases
func Test_with_keep(t *testing.T) {

	common.SetRoot("../../")

	fs := filestore.NewFileStore(common.TestStore("cache1"))
	ref := StoreDirTree(fs)

	t.Logf("LOCAL HTTP TEST -- %s\n", ref.Dump())
	CompareDirTree(common.TestStore("cache1"), ref, keep_test_temp_file)
	t.Logf("COMPARE TREE OKAY")
	common.PanicX("KeepBlock", common.TestStore("cache1").KeepBlock(ref))

	err := common.TestStore("cache1").KeepBlock(base.BlockRef{})
	if errors.Is(err, base.ErrNotFound) {
		t.Logf("test keepblock: BLOCK NOT FOUND - %s\n", err)
	} else {
		t.Logf("keepblock err: %s\n", err)
	}

	newref := StoreDirTree(filestore.NewFileStore(common.TestStore("cache1")))
	t.Logf("stored to http: newref=%s\n", newref.Dump())

	inf, err := common.TestStore("cache1").PeekBlock(newref, 12)
	t.Logf("HTTP PEEK inf=%+v err=%s\n", inf, err)

	t.Logf("iterating ...\n")
	for ent := range common.TestStore("cache1").IterateBlocks() {
		if !ent.Finished {
			t.Logf("==> HTTP ITERATE %s %s\n", ent.Ref.Dump(), ent.Error)
		}
	}
	t.Logf("done interating\n")
}
