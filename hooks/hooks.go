package hooks

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/errwrap"
	"github.com/wedeploy/cli/verbose"
)

// Hooks (after / deploy / main action)
type Hooks struct {
	BeforeBuild  string `json:"before_build"`
	Build        string `json:"build"`
	AfterBuild   string `json:"after_build"`
	BeforeDeploy string `json:"before_deploy"`
	Deploy       string `json:"deploy"`
	AfterDeploy  string `json:"after_deploy"`
}

// HookError struct
type HookError struct {
	Command string
	Err     error
}

func (he HookError) Error() string {
	return fmt.Sprintf("Command %v failure: %v", he.Command, he.Err.Error())
}

// Build is 'build' hook
const Build = "build"

var (
	// ErrMissingHook is used when the hook is missing
	ErrMissingHook = errors.New("Missing hook.")

	outStream io.Writer = os.Stdout
	errStream io.Writer = os.Stderr
)

// Run invokes the hooks for the given hook type on working directory
func (h *Hooks) Run(hookType string, wdir string) error {
	var owd, err = os.Getwd()

	if err != nil {
		return errwrap.Wrapf("Can't get current working dir on hooks run: {{err}}", err)
	}

	if wdir != "" {
		if err = os.Chdir(wdir); err != nil {
			return err
		}
	}

	switch hookType {
	case "build":
		err = h.runBuild()
	default:
		err = ErrMissingHook
	}

	if wdir != "" {
		if ech := os.Chdir(owd); ech != nil {
			fmt.Fprintf(os.Stderr, "Multiple errors: %v\n", err)
			panic(ech)
		}
	}

	return err
}

func (h *Hooks) runBuild() error {
	if h.Build == "" && (h.BeforeBuild != "" || h.AfterBuild != "") {
		fmt.Fprintf(errStream, "Error: no build hook main action\n")
	}

	var steps = []string{
		h.BeforeBuild,
		h.Build,
		h.AfterBuild,
	}

	for _, eachStep := range steps {
		if eachStep == "" {
			continue
		}

		var err = Run(eachStep)

		if err != nil {
			return HookError{
				Command: eachStep,
				Err:     err,
			}
		}
	}

	return nil
}

// Run a process synchronously inheriting stderr and stdout
func Run(command string) error {
	verbose.Debug("> " + command)
	return run(command)
}
