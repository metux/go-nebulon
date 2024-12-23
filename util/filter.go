package util

type FileNameFilter = func(name string, path string) bool

func FilterSkipHidden(name string, path string) bool {
	if len(name) == 0 {
		return false
	}
	return name[0] != '.'
}
