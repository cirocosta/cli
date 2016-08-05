package createctx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/mitchellh/go-wordwrap"
	"github.com/wedeploy/cli/color"
	"github.com/wedeploy/cli/config"
	"github.com/wedeploy/cli/containers"
	"github.com/wedeploy/cli/projects"
	"github.com/wedeploy/cli/prompt"
)

var (
	// ErrContainerPath indicates an invalid container location
	ErrContainerPath = errors.New("A container immediate parent dir must be the root of a project")

	// ErrProjectPath indicates an invalid project location
	ErrProjectPath = errors.New("A project can not have another project as its parent")

	// ErrInvalidID indicates an invalid resource ID (such as empty string)
	ErrInvalidID = errors.New("Value for resource ID is invalid")

	// ErrResourceExists indicates that two resource can not share the same location
	ErrResourceExists = errors.New("A resource already exists on the root of this location")
)

// New creates a resource
func New() error {
	switch config.Context.Scope {
	case "project":
		return NewContainer()
	case "global":
		return NewProject()
	default:
		return ErrResourceExists
	}
}

// NewContainer creates a container resource
func NewContainer() error {
	var rel string
	var bin []byte

	if config.Context.Scope == "container" {
		return ErrResourceExists
	}

	if config.Context.Scope != "project" {
		return ErrContainerPath
	}

	projectRoot := config.Context.ProjectRoot
	workingDir, err := os.Getwd()

	if err == nil {
		rel, err = filepath.Rel(projectRoot, workingDir)
	}

	if err != nil {
		return err
	}

	// only allow container creation at first subdir level
	if strings.ContainsRune(rel, os.PathSeparator) {
		return ErrContainerPath
	}

	var c = &containers.Container{}

	if rel == "." {
		return ErrResourceExists
	}

	registry, err := containers.GetRegistry()

	if err != nil {
		return errwrap.Wrapf("Can't get the registry: {{err}}", err)
	}

	for pos, r := range registry {
		ne := fmt.Sprintf("%d) %v", pos+1, r.Name)

		p := 80 - len(ne) - len(r.Type) + 1

		if p < 1 {
			p = 1
		}

		fmt.Fprintf(os.Stdout, "%v%v%v\n", ne, pad(p), r.Type)
		fmt.Fprintf(os.Stdout, "%v\n\n", color.Format(color.FgHiBlack, wordwrap.WrapString(r.Description, 80)))
	}

	var option = prompt.Prompt(fmt.Sprintf("\nSelect container type from 1..%d", len(registry)))

	var index int

	index, err = strconv.Atoi(option)

	index--

	if err != nil || index < 0 || index > len(registry) {
		return errors.New("Invalid option.")
	}

	var reg = registry[index]

	c.ID = prompt.Prompt("Id [default: " + reg.ID + "]")

	if c.ID == "" {
		c.ID = reg.ID
	}

	c.Name = prompt.Prompt("Name [default: " + reg.Name + "]")

	if c.Name == "" {
		c.Name = reg.Name
	}

	c.Type = reg.Type

	bin, err = json.MarshalIndent(c, "", "    ")

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(
		filepath.Join(workingDir, "container.json"),
		bin,
		0644)

	return err
}

// NewProject creates a project resource
func NewProject() error {
	var bin []byte

	if config.Context.Scope != "global" {
		return ErrProjectPath
	}

	workingDir, err := os.Getwd()

	if err != nil {
		return err
	}

	var p = &projects.Project{}

	fmt.Println("Creating project:")
	p.ID = prompt.Prompt("ID")

	if p.ID == "" {
		return ErrInvalidID
	}

	p.Name = prompt.Prompt("Name")
	p.CustomDomain = prompt.Prompt("Custom domain")

	bin, err = json.MarshalIndent(p, "", "    ")

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(
		filepath.Join(workingDir, "project.json"),
		bin,
		0644)

	return err
}

func pad(space int) string {
	return strings.Join(make([]string, space), " ")
}
