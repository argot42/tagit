package dictionary

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/argot42/tagit/utils"
	"github.com/renstrom/fuzzysearch/fuzzy"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const BUFFER = 2048

const (
	ExactFail = iota
)

var (
	ExactMatchErr = errors.New("exact match")
)

type Dictionary struct {
	Path string
	Tags []string
}

func (d *Dictionary) Add(tag string) (err error) {
	// if no initialized return error
	if !d.Initialized() {
		return errors.New("The dictionary is not initialized")
	}

	// insert into structure
	utils.StringSortedInsert(tag, &d.Tags)

	// write to file
	f, err := os.OpenFile(d.Path, os.O_RDWR, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	err = sortedWrite(tag, f)

	return
}

func sortedWrite(tag string, f *os.File) (err error) {
	position, err := findPos([]byte(tag), f)
	if err != nil {
		return
	}

	err = writeInto([]byte(tag+"\n"), position, f)

	return
}

func findPos(tag []byte, file *os.File) (position int64, err error) {
	reader := bufio.NewReader(file)

	var read int
	for {
		var b []byte
		position += int64(read)

		b, err = reader.ReadBytes('\n')
		if err != nil {
			break
		}

		if bytes.Compare(tag, b) <= 0 {
			break
		}

		read = len(b)
	}

	if err != nil && err != io.EOF {
		return
	}

	return position, nil
}

func writeInto(tail []byte, position int64, file *os.File) (err error) {
	total := len(tail)
	b := make([]byte, BUFFER)

	_, err = file.Seek(position, 0)
	if err != nil {
		return
	}

	for {
		n, err := file.Read(b)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		total += n
		tail = append(tail, b[:n]...)
	}

	_, err = file.WriteAt(tail, position)

	return
}

func (d *Dictionary) Initialized() bool {
	return len(d.Path) != 0
}

func (d *Dictionary) Name() string {
	return d.Path
}

func LoadDictionary(path string) (d Dictionary, err error) {
	fileinfo, err := os.Stat(path)
	if err != nil {
		return
	}

	if fileinfo.IsDir() {
		path = filepath.Join(path, ".dict.db")
	}

	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	d.Path = path
	reader := bufio.NewReader(file)
	for {
		var line string
		line, err = reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}

		d.Tags = append(d.Tags, strings.TrimRight(line, "\r\n"))
	}

	return d, nil
}

func (d *Dictionary) FFind(tag string, flag int) (matches []string, err error) {
	switch flag {
	case ExactFail:
		rank := fuzzy.RankFindFold(tag, d.Tags)
		for _, v := range rank {
			if v.Distance == 0 {
				return []string{}, ExactMatchErr
			}
			matches = append(matches, v.Target)
		}

	default:
		matches = fuzzy.FindFold(tag, d.Tags)
	}

	return
}

func CreateDictionary(path string) (d Dictionary, err error) {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0664)
	if err != nil {
		return
	}
	defer file.Close()

	d.Path = path
	return
}
