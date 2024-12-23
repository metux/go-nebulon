package helpers

import (
	"fmt"

	"github.com/udhos/equalfile"

	"github.com/metux/go-nebulon/core/base"
)

func compareEntry(fs base.IFileStore, path string, entry base.BlockRef, test_temp_file string) error {
	fn := path + "/" + entry.Name

	if entry.IsFile() {
		if _, err := GetFile(fs, test_temp_file, entry); err != nil {
			return err
		}
		cmp := equalfile.New(nil, equalfile.Options{}) // compare using single mode
		if equal, err := cmp.CompareFile(test_temp_file, fn); err != nil || !equal {
			return fmt.Errorf("compare error [%w]", err)
		}
	} else if entry.IsDir() {
		if err := CompareTree(fs, fn, entry, test_temp_file); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unexpected object, neither dir nor file")
	}

	return nil
}

func CompareTree(fs base.IFileStore, path string, ref base.BlockRef, test_temp_file string) error {
	entries, err := fs.ReadDirectory(ref)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if err := compareEntry(fs, path, e, test_temp_file); err != nil {
			return err
		}
	}
	return nil
}
