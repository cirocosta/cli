package update

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	gup "github.com/inconshreveable/go-update"
	"github.com/launchpad-project/cli/launchpad"
	"github.com/launchpad-project/cli/launchpad/client"
	"github.com/launchpad-project/cli/launchpad/config"
)

var (
	globalConfigStore      = config.Stores["global"]
	signature              = []byte(`-----BEGIN/END PUBLIC KEY-----`)
	ErrNewestVersion       = errors.New("Newest version already instaled")
	ErrMasterVersion       = errors.New("You can't upgrade a master version")
	ErrPlatformUnsupported = errors.New("Build for your platform was not found")
	ErrPermission          = errors.New("Can't replace Launchpad binary")
)

type Release struct {
	Id       string `json:"id"`
	Link     string `json:"link"`
	Version  string `json:"version"`
	Platform string `json:"platform"`
	Checksum string `json:"checksum"`
}

func PostUpdate() {
	if launchpad.Version == "master" {
		return
	}
}

func GetReleases() []Release {
	var address = globalConfigStore.Data.GetString("endpoint") + "/releases/dist/channel"
	var os = runtime.GOOS
	var arch = runtime.GOARCH

	var b = strings.NewReader(fmt.Sprintf(
		`{"filter": [{"platform": {"value": "%s_%s"}}]}`, os, arch,
	))

	var l = client.Url(address, b)
	l.Get()

	releases := *new([]Release)

	if err := l.ResponseJson(&releases); err != nil {
		panic(err)
	}

	return releases
}

func Update(release Release) error {
	resp, err := http.Get(release.Link)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	checksum, err := hex.DecodeString(release.Checksum)

	if err != nil {
		return err
	}

	err = gup.Apply(resp.Body, gup.Options{
		Checksum:  checksum,
		Signature: signature,
	})

	return err
}

func UpdateToLatest() {
	if launchpad.Version == "master" {
		panic(ErrMasterVersion)
	}

	var releases = GetReleases()

	if len(releases) == 0 {
		println("Releases not found.")
	}

	var next = releases[0]

	if next.Version == launchpad.Version {
		fmt.Println("Installed version " + launchpad.Version + " is already the last version.")
		return
	}

	if err := Update(next); err != nil {
		panic(err)
	}
}
