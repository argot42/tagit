package utils

import (
	"sort"
	"path/filepath"
	"strings"
)

func Prepend (s string, slice []string) []string {
	new_slice := make([]string, len(slice)+1)
	new_slice[0] = s
	copy(new_slice[1:], slice)

	return new_slice
}

func StringSortedInsert(s string, slc *[]string) {
	index := sort.SearchStrings(*slc, s)

	*slc = append(*slc, "")
	copy((*slc)[index+1:], (*slc)[index:])
	(*slc)[index] = s
}

func SplitNameExt (path string) (name, ext string) {
	name = filepath.Base(path)

	if len(name) < 3 { return }

	index := -1
	for i:=1; i<len(name); i++ {
		if name[i] == '.' {
			index = i
			break
		}
	}
	if index < 0 { return }

	ext = name[index:]

	return strings.TrimSuffix(name, ext), ext
}
