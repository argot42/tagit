package options

import (
	"fmt"
	"strings"
	"strconv"
	"errors"
)

// default answer constants
const None = -1
const ( // YesNo
	Yes = iota
	No
)

func YesNo (defaultAnswer int, format string, a ...interface{}) bool {
	var options string
	yes := []string{ "y", "yes" }

	switch defaultAnswer {
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

func ChooseNumeric (defaultAnswer int, options interface{}, format string, a ...interface{}) int {
	switch o := options.(type) {
	case []string:
		for _, v := range o { fmt.Println(v) }
	case []int:
		for _, v := range o { fmt.Println(v) }
	case []float32:
		for _, v := range o { fmt.Println(v) }
	case []float64:
		for _, v := range o { fmt.Println(v) }
	default:
		panic(errors.New("Type " + fmt.Sprint(o) + " not supported"))
	}

	var optionBox string

	switch defaultAnswer {
	case None:
		optionBox = ": "
	default:
		optionBox = "[" + strconv.Itoa(defaultAnswer) + "] "
	}

	input := defaultAnswer
	fmt.Printf(format + optionBox, a)
	fmt.Scanln(&input)

	return input
}
