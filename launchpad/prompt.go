package launchpad

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/bgentry/speakeasy"
)

func isSecretKey(key string) bool {
	var match, _ = regexp.MatchString("(password|token|secret)", strings.ToLower(key))
	return match
}

func Prompt(param string) string {

	var value string

	if isSecretKey(param) {
		value, err := speakeasy.Ask(param + ": ")

		if err != nil {
			panic(err)
		}

		return value
	}

	fmt.Fprintf(os.Stdout, param+": ")
	fmt.Fscanf(os.Stdin, "%s\n", &value)
	return value
}
