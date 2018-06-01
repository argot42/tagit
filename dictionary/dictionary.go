package dictionary

import (
	"os"
	"io"
	"sort"
	"bufio"
	"strings"
	"path/filepath"
)

type Dictionary struct {
	Path string
	Tags []string
}

func (d *Dictionary) Add (tag string) (err error) {
	// insert into structure
	sortedInsert(tag, &d.Tags)

	// write to file
	f, err := os.Open(d.Path)
	if err != nil {
		return
	}
	defer f.Close()
	err = sortedWrite(tag, f)

	return
}

func sortedInsert (tag string, s *[]string) {
	index := sort.SearchStrings(*s, tag)

	*slc = append(*slc, "")
	copy((*slc)[index+1:], (*slc)[index:])
	(*slc)[index] = tag
}

func sortedWrite (tag string, f *os.File) error {
}

func (d *Dictionary) Empty() bool {
	return len(d.Path) == 0
}

func (d *Dictionary) Name() string {
	return Path
}

func LoadDictionary (path string) (d Dictionary, err error) {
	fileinfo, err := os.Stat(path)
	if err != nil { return }

	if fileinfo.IsDir() {
		path = filepath.Join(path, ".dict.db")
	}

	file, err := os.Open(path)
	if err != nil { return }
	defer file.Close()

	d.Path = path
	reader := bufio.NewReader(file)
	for {
		var line string
		line, err = reader.ReadString('\n')
		if err != nil {
			if err == io.EOF { break }
			return
		}

		d.Tags = append(d.Tags, strings.TrimRight(line, "\r\n"))	
	}

	return d, nil
}

func CreateDictionary (path string) (d Dictionary, err error) {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil { return }
	defer file.Close()

	return
}
