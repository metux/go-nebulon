package util

import "log"

type FileNameFilter = func(name string, path string) bool

func FilterSkipHidden(name string, path string) bool {
	log.Printf("FilterSkipHidden: name=%s path=%s\n", name, path)
	if len(name) == 0 {
		return false
	}
	return name[0] != '.'
}
