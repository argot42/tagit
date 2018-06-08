package fileproc

import (
	"testing"
)

type FileTest struct {
	Type int
	Path string
}

const (
	SUCCESS = iota
	FAIL
)

func TestParse (t *testing.T) {
	list_of_tests := []FileTest{
		FileTest{ Type: SUCCESS, Path: "/tmp/test" },
		FileTest{ Type: FAIL, Path: "/tmp/test/foo" },
		FileTest{ Type: FAIL, Path: "/tmp/test_not_existent" },
		FileTest{ Type: SUCCESS, Path: "/tmp/test/foo0" },
		FileTest{ Type: SUCCESS, Path: "/tmp/test/foo1 [a b c]" },
		FileTest{ Type: SUCCESS, Path: "/tmp/test/foo2 [y x z].txt" },
		FileTest{ Type: SUCCESS, Path: "/tmp/test/foo3.txt" },
	}

	for _, f := range list_of_tests {
		file, err := Parse(f.Path)
		if err != nil {
			if f.Type == SUCCESS { t.Errorf("%s - %s\n", f.Path, err.Error()) }
			continue
		}

		t.Logf("%+v", file)

		if f.Type == FAIL { t.Errorf("%s - Should have failed\n", f.Path) }
	}
}

func TestWrite (t *testing.T) {
	list_of_tests := []FileTest{	
		FileTest{ Type: FAIL, Path: "/tmp/test/foo" },
		FileTest{ Type: FAIL, Path: "/tmp/test_not_existent" },
		FileTest{ Type: SUCCESS, Path: "/tmp/test/foo0" },
		FileTest{ Type: SUCCESS, Path: "/tmp/test/foo1 [a b c]" },
		FileTest{ Type: SUCCESS, Path: "/tmp/test/foo2 [y x z].txt" },
		FileTest{ Type: SUCCESS, Path: "/tmp/test/foo3.txt" },
	}

	for _, f := range list_of_tests {
		file, err := Parse(f.Path)
		if err != nil {
			if f.Type == SUCCESS { t.Errorf("%s - %s\n", f.Path, err.Error()) }
			continue
		}

		file.Add("hola")
		file.Add("chau")
		file.Add("nos vemos")

		t.Logf("%+v", file)

		if err := file.Write(); err != nil {
			if f.Type == SUCCESS { t.Errorf("%s - %s\n", f.Path, err.Error()) }
		}

		if f.Type == FAIL { t.Errorf("%s - Should have failed\n", f.Path) }
	}
}
