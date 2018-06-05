package fileproc

import (
	"errors"
	"path/filepath"
)

const (
	ExactFail = iota
)

var (
	ExactMatchErr = errors.New("exact match")
)

type File struct {
	path string
}

func (f *File) Add (tag string) (err error) {
}

func (f *File) Name() string {
	return filepath.Base(f.path)
}

func (f *File) Path() (path string) {
	return f.path
}

func Parse (path string) (file *File, err error) {
}
