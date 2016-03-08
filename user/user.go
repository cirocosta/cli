package user

import (
	"os"
	"runtime"
)

// GetHomeDir returns the user's ~ (home)
// Extracted from Viper's util.go GetUserHomeDir method
func GetHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
