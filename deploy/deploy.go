package deploy

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/launchpad-project/cli/apihelper"
	"github.com/launchpad-project/cli/verbose"
)

func DeployContainer(container string) error {
	tmp, err := ioutil.TempFile(os.TempDir(), "launchpad-cli")

	if err == nil {
		err = ZipContainer(container, tmp.Name())
	}

	if err == nil {
		err = tmp.Close()
	}

	if err != nil {
		return err
	}

	verbose.Debug("Saving container to", tmp.Name())

	err = os.Remove(tmp.Name())

	if err != nil {
		println("Can not remove temporary file " + tmp.Name())
	}

	return deployContainer(container, tmp.Name())
}

func DeployContainers(list []string) {
	for _, container := range list {
		fmt.Println("Deploying", container)
		DeployContainer(container)
	}
}

func deployContainer(container, jar string) error {
	var req = apihelper.URL("/_pods/admin/" + container)

	apihelper.Auth(req)
	// @todo use something else instead of ValidateOrExit
	apihelper.ValidateOrExit(req, req.Put())

	file, err := os.Open(jar)

	if err == nil {
		_, err = io.Copy(file, req.RequestBody)
	}

	if err == nil {
		err = req.Put()
	}

	if err := file.Close(); err != nil {
		println(err.Error())
	}

	return err
}
