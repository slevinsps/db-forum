package utils

import "fmt"

func PrintDebug(str string, args ...interface{}) {
	return
	for i := range args {
		str += fmt.Sprintf("%v", args[i])
	}
	fmt.Println(str)
}
