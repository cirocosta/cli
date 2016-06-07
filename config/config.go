package config

import (
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"

	"github.com/launchpad-project/cli/context"
	"github.com/launchpad-project/cli/defaults"
	"github.com/launchpad-project/cli/user"
	"github.com/launchpad-project/cli/verbose"
)

// Config of the application
type Config struct {
	Username        string    `ini:"username"`
	Password        string    `ini:"password"`
	Token           string    `ini:"token"`
	Local           bool      `ini:"local"`
	NoColor         bool      `ini:"disable_colors"`
	Endpoint        string    `ini:"endpoint"`
	NotifyUpdates   bool      `ini:"notify_updates"`
	ReleaseChannel  string    `ini:"release_channel"`
	LastUpdateCheck string    `ini:"last_update_check"`
	NextVersion     string    `ini:"next_version"`
	Path            string    `ini:"-"`
	file            *ini.File `ini:"-"`
}

var (
	// Global configuration
	Global *Config

	// Context stores the environmental context
	Context *context.Context
)

// Load the configuration
func (c *Config) Load() {
	c.setDefaults()
	c.read()

	if err := c.file.MapTo(c); err != nil {
		panic(err)
	}
}

// Save the configuration
func (c *Config) Save() {
	var err = c.file.ReflectFrom(c)

	if err != nil {
		panic(err)
	}

	if c.NextVersion == "" {
		c.file.Section("").DeleteKey("next_version")
	}

	err = c.file.SaveTo(c.Path)

	if err != nil {
		panic(err)
	}
}

// Setup the environment
func Setup() {
	setupContext()
	setupGlobal()
}

func (c *Config) setDefaults() {
	c.Local = true
	c.Endpoint = defaults.Endpoint
	c.NotifyUpdates = true
	c.ReleaseChannel = "stable"
}

func (c *Config) configExists() bool {
	var _, err = os.Stat(c.Path)

	switch {
	case err == nil:
		return true
	case os.IsNotExist(err):
		return false
	default:
		panic(err)
	}
}

func (c *Config) read() {
	if !c.configExists() {
		verbose.Debug("Config file not found.")
		c.file = ini.Empty()
		c.banner()
		return
	}

	var err error
	c.file, err = ini.Load(c.Path)

	if err != nil {
		println("Error reading configuration file:", err.Error())
		println("Fix " + c.Path + " by hand or erase it.")
		os.Exit(1)
	}

}

func (c *Config) banner() {
	c.file.Section("DEFAULT").Comment = `# Configuration file for WeDeploy CLI
# https://wedeploy.io`
}

func setupContext() {
	var err error
	Context, err = context.Get()

	if err != nil {
		println(err.Error())
		os.Exit(-1)
	}
}

func setupGlobal() {
	Global = &Config{
		Path: filepath.Join(user.GetHomeDir(), ".we"),
	}

	Global.Load()
}

// Teardown resets the configuration environment
func Teardown() {
	teardownContext()
	teardownGlobal()
}

func teardownContext() {
	Context = nil
}

func teardownGlobal() {
	Global = nil
}
