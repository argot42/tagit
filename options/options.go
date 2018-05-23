package options

import (
	"fmt"
	"strings"
)

const (
	None = 0
	Yes = iota
	No
)

func YesNo (default_answer uint, format string, a ...interface{}) bool {
	var options string
	yes := []string{ "y", "yes" }

	switch default_answer {
	case None:
		options = " [y/n] "
	case Yes:
		options = " [Y/n] "
		yes = append(yes, "")
	case No:
		options = " [y/N] "
	}

	fmt.Printf(format + options, a...)
	
	var response string
	fmt.Scanln(&response)
	response = strings.ToLower(response)

	return in(response, yes)
}

func in (str string, list []string) bool {
	for _,item := range list {
		if str == item { 
			return true 
		}
	}
	return false
}
