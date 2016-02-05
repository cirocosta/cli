package update

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	gup "github.com/inconshreveable/go-update"
	"github.com/launchpad-project/cli/config"
	"github.com/launchpad-project/cli/launchpad"
)

var (
	globalConfigStore = config.Stores["global"]
	signature         = []byte(`-----BEGIN/END PUBLIC KEY-----`)
	// ErrMasterVersion for when trying to run launchpad upgrade on a master distribution
	ErrMasterVersion = errors.New("You must upgrade a master version manually with git")
	// ErrPlatformUnsupported for when a release for the given platform was not found
	ErrPlatformUnsupported = errors.New("Build for your platform was not found")
	// ErrPermission for when there is no permission to replace the binary
	ErrPermission = errors.New("Can't replace Launchpad binary")
)

// Release information
type Release struct {
	ID       string `json:"id"`
	Link     string `json:"link"`
	Version  string `json:"version"`
	Platform string `json:"platform"`
	Checksum string `json:"checksum"`
}

// GetLatestReleases lists the latest releases available for the given platform
func GetLatestReleases() []Release {
	var address = globalConfigStore.Get("endpoint") + "/releases/dist/channel"
	var os = runtime.GOOS
	var arch = runtime.GOARCH

	var b = strings.NewReader(fmt.Sprintf(
		`{"filter": [{"platform": {"value": "%s/%s"}}]}`, os, arch,
	))

	var l, err = launchpad.URL(address, b)

	if err != nil {
		panic(err)
	}

	if err := l.Get(); err != nil {
		panic(err)
	}

	releases := *new([]Release)

	if err := l.DecodeJSON(&releases); err != nil {
		panic(err)
	}

	return releases
}

// Update updates the Launchpad CLI to the given release
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

// ToLatest updates the Launchpad CLI to the latest version
func ToLatest() {
	if launchpad.Version == "master" {
		panic(ErrMasterVersion)
	}

	var releases = GetLatestReleases()

	if len(releases) == 0 {
		println("Releases not found.")
	}

	var next = releases[0]

	if next.Version == launchpad.Version {
		fmt.Println("Installed version " + launchpad.Version + " is already the latest version.")
		return
	}

	if err := Update(next); err != nil {
		panic(err)
	}
}
