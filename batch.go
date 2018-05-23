package main

import (
	"fmt"
	"github.com/argot42/tagit/flags"
	"github.com/argot42/tagit/dictionary"
)

func batch (defaultDict dictionary.Dictionary) {
	if flags.Args.Verbose { fmt.Println("Starting batch tagging") }

	for _,item := range flags.Args.Files {

	}
}
