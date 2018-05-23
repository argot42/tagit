package main

import (
	"fmt"
	"os"
	"github.com/argot42/tagit/flags"
	"github.com/argot42/tagit/dictionary"
	"github.com/argot42/tagit/options"
)

func main() {	
	err := flags.Init_flags()
	checkError(err)

	if flags.Args.Verbose { fmt.Println("Loading default dictionary...") }

	defaultDictionary, err := dictionary.LoadDictionary(flags.Args.DB)
	if os.IsNotExist(err) && flags.Args.Interactive || flags.Args.Client {
		if !options.YesNo(options.Yes, "Dictionary not found, create a new one (%s)", flags.Args.DB) {
			break
		}
		defaultDict, err = dictionary.CreateDictionary(flags.Args.DB)
	}
	checkError(err)

	if flags.Args.Client {
		client(defaultDictionary)
	} else {
		batch(defaultDictionary)
	}
}
