package fileproc

import (
	"os"
	"regexp"
	"errors"
	"path/filepath"
)

var r = regexp.MustCompile(`^[\w\s]+\[([\w\s]*)\][\.\w]*`)

type File struct {
	path string
	name string
	brackets bool
	tags []string
}

func (f *File) Add (tag string) (err error) {
}

func (f *File) Name() string {
	return f.name
}

func (f *File) Path() (path string) {
	return f.path
}

func Parse (path string) (file *File, err error) {
	fi, err := os.Stat(path)
	if err != nil { return }

	match := r.FindAllStringSubmatch(fi.Name(), -1)

	switch len(match) {
	case 0:
		file = &File{
			path: path,
			name: fi.Name(),
		}
	default:
		file = &File{
			path: path,
			name: fi.Name(),
			brackets: true,
		}
		if len(match[0]) > 1 { file.tags = match[0][1:] }
	}

	return
}
