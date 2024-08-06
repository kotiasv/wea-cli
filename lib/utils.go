package lib

import (
	"fmt"
	"os"
	"strings"
)

func HandleErr(err error, condition bool) {
	if condition {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func Join(s []string, sep string) string {
	return strings.Join(s, sep)
}
