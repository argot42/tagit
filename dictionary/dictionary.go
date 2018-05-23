package dictionary

import (
	"os"
	"io"
	"path/filepath"
	"bufio"
	"strings"
)

type Dictionary struct {
	Path string
	Words []string
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

		d.Words = append(d.Words, strings.TrimRight(line, "\r\n"))	
	}

	return d, nil
}

func CreateDictionary (path string) (d Dictionary, err error) {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil { return }
	defer file.Close()

	return
}
