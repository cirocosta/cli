package hooks

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/wedeploy/cli/tdata"
)

type HooksProvider struct {
	Type       string
	Hook       *Hooks
	WantOutput string
	WantErr    string
	WantError  error
}

var HooksCases = []HooksProvider{
	HooksProvider{
		Type: "build",
		Hook: &Hooks{
			BeforeBuild: "echo before build",
			Build:       "echo during build",
			AfterBuild:  "echo after build",
		},
		WantOutput: "before build\nduring build\nafter build\n",
		WantError:  nil,
	},
	HooksProvider{
		Type: Build,
		Hook: &Hooks{
			BeforeBuild: "echo a",
			AfterBuild:  "echo b",
		},
		WantOutput: "a\nb\n",
		WantErr:    "Error: no build hook main action\n",
		WantError:  nil,
	},
	HooksProvider{
		Type:      "not implemented",
		WantError: ErrMissingHook,
	},
}

var (
	bufErrStream bytes.Buffer
	bufOutStream bytes.Buffer
)

func TestMain(m *testing.M) {
	var defaultErrStream = errStream
	var defaultOutStream = outStream
	errStream = &bufErrStream
	outStream = &bufOutStream
	ec := m.Run()
	errStream = defaultErrStream
	outStream = defaultOutStream
	os.Exit(ec)
}

func TestRunHooks(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Not testing hooks.Build() on Windows")
	}

	for _, c := range HooksCases {
		bufErrStream.Reset()
		bufOutStream.Reset()

		if err := c.Hook.Run(c.Type); err != c.WantError {
			t.Errorf("Expected %v, got %v instead", c.WantError, err)
		}

		var gotOutStream = bufOutStream.String()
		var gotErrStream = bufErrStream.String()

		if gotErrStream != c.WantErr {
			t.Errorf("Expected %v, got %v instead", c.WantErr, gotErrStream)
		}

		if gotOutStream != c.WantOutput {
			t.Errorf("Expected %v, got %v instead", c.WantOutput, gotOutStream)
		}
	}
}

func TestRun(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Not testing hooks.Run() on Windows")
	}

	bufErrStream.Reset()
	bufOutStream.Reset()

	if err := Run("openssl md5 hooks.go"); err != nil {
		t.Errorf("Unexpected error %v when running md5 hooks.go", err)
	}

	h := md5.New()

	if _, err := io.WriteString(h, tdata.FromFile("./hooks.go")); err != nil {
		panic(err)
	}

	if !strings.Contains(bufOutStream.String(), fmt.Sprintf("%x", h.Sum(nil))) {
		t.Errorf("Expected Run() test to contain md5 output similar to crypto.md5")
	}

	if bufErrStream.Len() != 0 {
		t.Errorf("Unexpected err output")
	}
}

func TestRunAndExitOnFailureOnSuccess(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Not testing hooks.Run() on Windows")
	}

	bufErrStream.Reset()
	bufOutStream.Reset()

	RunAndExitOnFailure("openssl md5 hooks.go")

	h := md5.New()

	if _, err := io.WriteString(h, tdata.FromFile("./hooks.go")); err != nil {
		panic(err)
	}

	if !strings.Contains(bufOutStream.String(), fmt.Sprintf("%x", h.Sum(nil))) {
		t.Errorf("Expected Run() test to contain md5 output similar to crypto.md5")
	}

	if bufErrStream.Len() != 0 {
		t.Errorf("Unexpected err output")
	}
}

func TestRunAndExitOnFailureFailure(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Not testing hooks.Run() on Windows")
	}

	bufErrStream.Reset()
	bufOutStream.Reset()

	if os.Getenv("PING_CRASHER") == "1" {
		outStream = os.Stdout
		errStream = os.Stderr
		RunAndExitOnFailure("ping")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestRunAndExitOnFailureFailure")
	cmd.Stderr = errStream
	cmd.Stdout = outStream
	cmd.Env = append(os.Environ(), "PING_CRASHER=1")
	err := cmd.Run()

	if err.Error() != "exit status 1" {
		t.Errorf("Expected exit status 1 for ping process, got %v instead", err)
	}

	if bufErrStream.Len() == 0 {
		t.Error("Expected ping output to be piped to stderr")
	}
}
