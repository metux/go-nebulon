package dirtree

import (
	"testing"

	"github.com/metux/go-nebulon/tests/common"
)

func Test_Dirtree_local(t *testing.T) {
	common.SetRoot("../../")
	ref := RunDirTree(common.TestStore("cache1"), ".download.tmp")
	t.Logf("RunDirTree: ref=%+v\n", ref)
}
