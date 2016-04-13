package deploy

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/dustin/go-humanize"
	"github.com/launchpad-project/api.go"
	"github.com/launchpad-project/cli/apihelper"
	"github.com/launchpad-project/cli/config"
	"github.com/launchpad-project/cli/containers"
	"github.com/launchpad-project/cli/hooks"
	"github.com/launchpad-project/cli/pod"
	"github.com/launchpad-project/cli/progress"
	"github.com/launchpad-project/cli/projects"
	"github.com/launchpad-project/cli/verbose"
)

// Deploy holds the information of a POD to be zipped or deployed
type Deploy struct {
	Container     containers.Container
	ContainerPath string
	PackageSize   int64
	progress      *progress.Bar
}

// DeployFlags modifiers
type DeployFlags struct {
	Hooks bool
}

// ErrDeploy is a generic error triggered when any deploy error happens
var ErrDeploy = errors.New("Error during deploy")

// All deploys a list of containers on the given context
func All(list []string, df *DeployFlags) (err error) {
	var wg sync.WaitGroup
	var el []error

	wg.Add(len(list))

	for _, container := range list {
		go func(c string) {
			el = append(el, Only(c, df))
			wg.Done()
		}(container)
	}

	wg.Wait()
	progress.Stop()

	for _, e := range el {
		if e == nil {
			continue
		}

		println(e.Error())
		err = ErrDeploy
	}

	return err
}

// Only PODify a container and deploys it to Launchpad
func Only(container string, df *DeployFlags) error {
	var deploy, err = New(container)

	if err != nil {
		return err
	}

	var projectID = config.Stores["project"].Get("id")

	err = validateOrCreateProject(projectID, config.Stores["project"].Get("name"))

	if err != nil {
		return err
	}

	err = validateOrCreateContainer(projectID, deploy.Container)

	if err != nil {
		return err
	}

	return runDeploy(deploy, df)
}

// New Deploy instance
func New(container string) (*Deploy, error) {
	var deploy = &Deploy{
		ContainerPath: path.Join(config.Context.ProjectRoot, container),
		progress:      progress.New(container),
	}

	var err = containers.GetConfig(deploy.ContainerPath, &deploy.Container)

	if err != nil {
		return nil, err
	}

	return deploy, err
}

// Zip packages a POD to a .pod package
func Zip(dest, container string) error {
	var deploy, err = New(container)

	if err == nil {
		err = deploy.Zip(dest)
	}

	return err
}

// Deploy POD to Launchpad
func (d *Deploy) Deploy(pod string) (err error) {
	var projectID = config.Stores["project"].Get("id")
	var u = path.Join("api/push", projectID, d.Container.ID)
	var req = apihelper.URL(u)
	var file io.Reader

	apihelper.Auth(req)

	w := &writeCounter{
		progress: d.progress,
		Size:     uint64(d.PackageSize),
	}

	d.progress.Reset("Uploading", "")
	file, err = os.Open(pod)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, errCreate := writer.CreateFormFile("pod", "container.pod")

	if errCreate != nil {
		return errCreate
	}

	_, errCopy := io.Copy(part, file)

	if errCopy != nil {
		return errCopy
	}

	writer.Close()

	if err == nil {
		req.Body(io.TeeReader(body, w))
	}

	req.Headers.Set("Content-Type", writer.FormDataContentType())

	if err == nil {
		err = apihelper.Validate(req, req.Post())
	}

	if err == nil || err == launchpad.ErrUnexpectedResponse {
		d.progress.Append = fmt.Sprintf(
			"%s (Complete)",
			humanize.Bytes(uint64(d.PackageSize)))
		d.progress.Set(progress.Total)
	}

	if err == nil {
		fmt.Printf(fmt.Sprintf("Ready! %v.%v.liferay.io\n", d.Container.ID, projectID))
	}

	return err
}

// Only PODify a container and deploys it to Launchpad
func (d *Deploy) Only() error {
	tmp, err := ioutil.TempFile(os.TempDir(), "launchpad-cli")

	err = d.Zip(tmp.Name())

	if err == nil {
		err = d.Deploy(tmp.Name())
	}

	if tmp != nil {
		os.Remove(tmp.Name())
	}

	return err
}

// Zip packages a POD to a .pod package
func (d *Deploy) Zip(dest string) (err error) {
	d.progress.Reset("Zipping", "")
	dest, _ = filepath.Abs(dest)

	var ignorePatterns = append(d.Container.DeployIgnore, pod.CommonIgnorePatterns...)

	d.PackageSize, err = pod.Compress(
		dest,
		d.ContainerPath,
		ignorePatterns,
		d.progress)

	if err == nil {
		d.progress.Set(progress.Total)
	}

	verbose.Debug("Saving container to", dest)

	return err
}

func runDeploy(deploy *Deploy, df *DeployFlags) (err error) {
	var ch = deploy.Container.Hooks

	if df.Hooks && ch != nil && ch.BeforeDeploy != "" {
		err = hooks.Run(ch.BeforeDeploy)
	}

	if err == nil {
		err = deploy.Only()
	}

	if err == nil && df.Hooks && ch != nil && ch.AfterDeploy != "" {
		err = hooks.Run(ch.AfterDeploy)
	}

	return err
}

func validateOrCreateContainer(projectID string, c containers.Container) (err error) {
	err = containers.Validate(projectID, c.ID)

	if err == containers.ErrContainerAlreadyExists {
		return nil
	}

	if err == nil {
		err = containers.Install(projectID, c)

		if err == nil {
			fmt.Println("New container installed")
		}
	}

	return err
}

func validateOrCreateProject(projectID, projectName string) (err error) {
	err = projects.Validate(projectID)

	if err == projects.ErrProjectAlreadyExists {
		return nil
	}

	if err == nil {
		err = projects.Create(projectID, config.Stores["project"].Get("name"))

		if err == nil {
			fmt.Println("New project created")
		}
	}

	return err
}
