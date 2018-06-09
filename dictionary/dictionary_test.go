package dictionary

import (
	"os"
	"sort"
	"testing"
	"path/filepath"
)

type PathTest struct {
	Type int
	Path string
}

const (
	SUCCESS = 0
	FAIL = iota
)

func TestNewDict (t *testing.T) {
	list_of_tests := []PathTest{ 
		PathTest{ Type: FAIL, Path: "/etc" },
		PathTest{ Type: FAIL, Path: "~/programming" },
		PathTest{ Type: FAIL, Path: "~" },
		PathTest{ Type: FAIL, Path: ".." },
		PathTest{ Type: FAIL, Path: "." },
		PathTest{ Type: SUCCESS, Path: "/tmp" },
		PathTest{ Type: SUCCESS, Path: "/tmp/" },
		PathTest{ Type: SUCCESS, Path: "/tmp/.dict.db" },
	}

	for _,path := range list_of_tests {
		d, err := LoadDictionary(path.Path)

		if err != nil {
			if path.Type == SUCCESS {
				t.Errorf("%s - %s\n", path.Path, err.Error())
			}

			continue
		}

		if path.Type == FAIL {
			t.Errorf("%s - Should have failed\n", path.Path)
			continue
		}

		t.Log(d)
	}
}

func TestCreateDict (t *testing.T) {
	list_of_tests := []PathTest{ 
		PathTest{ Type: SUCCESS, Path: "/tmp" },
		PathTest{ Type: SUCCESS, Path: "/tmp/test" },
		//PathTest{ Type: SUCCESS, Path: "./test" },
		//PathTest{ Type: SUCCESS, Path: "../test" },
	}

	for _,path := range list_of_tests {
		dictPath := filepath.Join(path.Path, ".dict.db")

		_, err := CreateDictionary(dictPath)
		
		if err != nil {
			if path.Type == SUCCESS {
				t.Errorf("%s - %s\n", path.Path, err.Error())
			}

			continue
		}

		if path.Type == FAIL {
			t.Errorf("%s - Should have failed\n", path.Path)
			continue
		}

		err = os.Remove(dictPath)
		if err != nil {
			t.Errorf("Error while removing %s\n", path.Path)
			continue
		}
		t.Logf("Cleanup, removed %s\n", dictPath)
	}
}

func TestAddTag (t *testing.T) {
	dict, err := LoadDictionary("/tmp/test/test001")
	if err != nil { t.Fatal(err) }

	tags := []string{ "foo0", "a", "c", "b", "x", "i", "f", "d", "漢" }
	//tags := []string{ "e", "b", "a" }
	for _, s := range tags {
		err = dict.Add(s)

		if err != nil {
			t.Error(err)
			continue
		}
	}
}

func TestLoadDict (t *testing.T) {
	dict, err := LoadDictionary("/tmp/test001")
	if err != nil { t.Fatal(err) }

	t.Log(dict.Tags)

	if !sort.StringsAreSorted(dict.Tags) {
		t.Fatal("not sorted")
	}
}
