package flags

import (
	"os/user"
	"flag"
	"path/filepath"
)

type StrFlagSlice []string

func (sfs *StrFlagSlice) String() (s string) {
	for _,str := range *sfs { s += str + " " }
	return
}

func (sfs *StrFlagSlice) Set (value string) error {
	*sfs = append(*sfs, value)
	return nil
}

type Flags struct {
	DB string
	Verbose, Interactive, Client bool
	Tags StrFlagSlice
	Files []string
}

var Args Flags

func Init_flags() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}
	flag.StringVar(&Args.DB, "d", filepath.Join(usr.HomeDir, ".local/share/tagit/dict.db"), "default dictionary path for fuzzy matching")

	flag.BoolVar(&Args.Verbose, "v", false, "verbose output")
	flag.BoolVar(&Args.Interactive, "i", false, "ask for confirmation after performing actions")
	flag.BoolVar(&Args.Client, "ii", false, "fire up the CLI client")

	flag.Var(&Args.Tags, "t", "list of tags")

	flag.Parse()

	Args.Files = flag.Args()
	return nil
}