package update

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/equinox-io/equinox"
	"github.com/fatih/color"
	"github.com/launchpad-project/cli/config"
	"github.com/launchpad-project/cli/defaults"
)

const lucFormat = "2006-01-02 15:04:05 -0700 MST"

var cacheNonAvailabilityDays = 4

// AppID is Equinox's app ID for this tool
var AppID = "app_6vxvxHVfPgz"

// ErrMasterVersion triggered when trying to update a development version
var ErrMasterVersion = errors.New(
	"You must update a development version manually with git")

// PublicKey is the public key for the certificate used with Equinox
var PublicKey = []byte(`
-----BEGIN ECDSA PUBLIC KEY-----
MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAE1+wgAlkBJRmtRwmZWfq9fa8dBlJ929hM
BzASHioHo6RP1V+4EKnAxaYXN4eWlgalxQ2BEr8TqYRM+uHPizteVR11wKfsO6S0
GENiOKpfivw5FIiTN14MZeMTagiKJUOq
-----END ECDSA PUBLIC KEY-----
`)

// NotifierCheck enquires equinox if a new version is available
func NotifierCheck() error {
	var csg = config.Stores["global"]

	if !isNotifyOn() {
		return nil
	}

	var nextVersion = csg.Get("cache.next_version")

	if nextVersion != "" && nextVersion != defaults.Version {
		return nil
	}

	// how long since last non availability result?
	var lastUpdate = csg.Get("cache.last_update_check")
	var luc, luce = time.Parse(lucFormat, lastUpdate)

	if luce == nil && time.Since(luc).Hours() < float64(cacheNonAvailabilityDays*24) {
		return nil
	}

	// save, just to be safe (e.g., if the check below breaks)
	csg.SetAndSave("cache.last_update_check", time.Now().String())

	var channel = csg.Get("release_channel")

	if channel == "" {
		channel = "stable"
	}

	var resp, err = check(channel)

	switch err {
	case nil:
		csg.SetAndSave("cache.next_version", resp.ReleaseVersion)
	case equinox.NotAvailableErr:
		csg.SetAndSave("cache.next_version", "")
		err = nil
	}

	return err
}

// Notify is called every time this tool executes to verify if it is outdated
func Notify() {
	var csg = config.Stores["global"]

	if !isNotifyOn() {
		return
	}

	var nextVersion = csg.Get("cache.next_version")

	if nextVersion != "" && nextVersion != defaults.Version {
		var channel = csg.Get("release_channel")
		var cmd = "we update"

		if channel != "" && channel != "stable" {
			cmd += " --channel " + channel
		}

		println(color.RedString(
			`WARNING: WeDeploy CLI tool is outdated. Run "` + cmd + `".`))
	}
}

// Update this tool
func Update(channel string) {
	var csg = config.Stores["global"]

	if defaults.Version == "master" {
		fmt.Fprintln(os.Stderr, ErrMasterVersion)
		os.Exit(1)
	}

	fmt.Println("Trying to update using the", channel, "distribution channel")
	fmt.Println("Current installed version is " + defaults.Version)

	var resp, err = check(channel)
	handleUpdateCheckError(err)

	err = resp.Apply()

	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	csg.Set("release_channel", channel)
	csg.Set("cache.next_version", "")
	csg.Set("cache.last_update_check", time.Now().String())
	csg.Save()

	fmt.Println("Updated to new version:", resp.ReleaseVersion)
}

func check(channel string) (*equinox.Response, error) {
	var opts equinox.Options
	opts.Channel = channel

	if err := opts.SetPublicKeyPEM(PublicKey); err != nil {
		return nil, err
	}

	resp, err := equinox.Check(AppID, opts)

	return &resp, err
}

func handleUpdateCheckError(err error) {
	var csg = config.Stores["global"]

	switch err {
	case nil:
	case equinox.NotAvailableErr:
		csg.Set("cache.next_version", "")
		csg.Set("cache.last_update_check", time.Now().String())
		csg.Save()
		fmt.Println("No updates available.")
		return
	default:
		println("Update failed:", err.Error())
		os.Exit(1)
	}
}

func isNotifyOn() bool {
	var csg = config.Stores["global"]

	if len(os.Args) == 2 && os.Args[1] == "update" {
		return false
	}

	if defaults.Version == "master" || csg.Get("notify_updates") == "false" {
		return false
	}

	return true
}
