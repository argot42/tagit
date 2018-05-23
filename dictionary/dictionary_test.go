package dictionary

import (
	"os"
	"testing"
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
		PathTest{ Type: FAIL, Path: "/tmp" },
		PathTest{ Type: SUCCESS, Path: "/tmp/test" },
		PathTest{ Type: SUCCESS, Path: "./test" },
		PathTest{ Type: SUCCESS, Path: "../test" },
	}

	for _,path := range list_of_tests {
		_, err := CreateDictionary(path.Path)	
		
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

		err = os.Remove(path.Path)
		if err != nil {
			t.Errorf("Error while removing %s\n", path.Path)
			continue
		}
		t.Logf("Cleanup, removed %s\n", path.Path)
	}
}
