package util

import (
	"os"
	"runtime"
)

var Verbose = false

func Debug(message string) {
	if Verbose {
		println(message)
	}
}

/* Function extracted from Viper's util.go */
func GetUserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
