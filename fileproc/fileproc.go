package fileproc

import (
	"github.com/argot42/tagit/utils"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var r = regexp.MustCompile(`^([\p{L}\s\w]*)\[?([\p{L}\s\w]*)\]?$`)

type File struct {
	path     string
	name     string
	ext      string
	brackets bool
	tags     []string
}

func (f *File) Add(tag string) {
	utils.StringSortedInsert(tag, &f.tags)
}

func (f *File) Write() error {
	return os.Rename(f.path, filepath.Join(
		filepath.Dir(f.path),
		f.name+"["+strings.Join(f.tags, " ")+"]"+f.ext))
}

func (f *File) Name() string {
	return f.name + f.ext
}

func (f *File) Path() (path string) {
	return f.path
}

func Parse(path string) (file *File, err error) {
	fi, err := os.Stat(path)
	if err != nil {
		return
	}

	name, extension := utils.SplitNameExt(fi.Name())
	file = &File{
		path: path,
		ext:  extension,
	}

	match := r.FindStringSubmatch(name)

	file.name = match[1] // parsed filename
	if len(match[2]) > 0 {
		file.brackets = true
		file.tags = strings.Split(match[2], " ")
	}

	return
}
