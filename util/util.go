package util

import (
	"encoding/json"
	"os"
	"reflect"
	"runtime"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

// Verbose flag
var Verbose = false

// Debug message only available on verbose mode
func Debug(message string) {
	if Verbose {
		println(message)
	}
}

// GetUserHomeDir returns the user's ~ (home)
// Extracted from Viper's util.go
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

func AssertJSON(t *testing.T, want string, got interface{}) {
	var wantJSON interface{}
	var gotMap interface{}

	bin, err := json.Marshal(got)

	if err != nil {
		t.Error(err)
	}

	if err = json.Unmarshal([]byte(want), &wantJSON); err != nil {
		t.Errorf("Wanted value %s isn't JSON.", want)
	}

	if err = json.Unmarshal(bin, &gotMap); err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(wantJSON, gotMap) {
		t.Errorf("Expected %s, got %s instead\n\n%s",
			want, string(bin), pretty.Compare(wantJSON, gotMap))
	}
}
