package config

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"gopkg.in/ini.v1"

	"github.com/hashicorp/errwrap"
	"github.com/wedeploy/cli/defaults"
	"github.com/wedeploy/cli/user"
	"github.com/wedeploy/cli/usercontext"
	"github.com/wedeploy/cli/verbose"
)

// RemoteConfig for a remote
type RemoteConfig struct {
	URL        string
	URLComment string
	Comment    string
}

// Remotes (list of alternative endpoints)
type Remotes struct {
	list remotesList
}

type remotesList map[string]RemoteConfig

// List remotes
func (r *Remotes) List() []string {
	var keys = make([]string, 0, len(r.list))

	for k := range r.list {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}

// Get a given remote by name
func (r *Remotes) Get(name string) (RemoteConfig, bool) {
	remote, ok := r.list[name]
	return remote, ok
}

// Set a remote
func (r *Remotes) Set(name string, url string, comment ...string) {
	// make sure to use # by default, instead of ;
	if len(comment) != 0 {
		comment = append([]string{"#"}, comment...)
	}

	r.list[name] = RemoteConfig{
		URL:     url,
		Comment: strings.Join(comment, " "),
	}
}

// Del deletes a remote by name
func (r *Remotes) Del(name string) {
	delete(r.list, name)
}

// Config of the application
type Config struct {
	Username        string    `ini:"username"`
	Password        string    `ini:"password"`
	Token           string    `ini:"token"`
	Local           bool      `ini:"local"`
	LocalPort       int       `ini:"local_port"`
	NoColor         bool      `ini:"disable_colors"`
	Endpoint        string    `ini:"endpoint"`
	NotifyUpdates   bool      `ini:"notify_updates"`
	ReleaseChannel  string    `ini:"release_channel"`
	LastUpdateCheck string    `ini:"last_update_check"`
	NextVersion     string    `ini:"next_version"`
	Path            string    `ini:"-"`
	Remotes         Remotes   `ini:"-"`
	file            *ini.File `ini:"-"`
}

var (
	// Global configuration
	Global *Config

	// Context stores the environmental context
	Context *usercontext.Context
)

// Load the configuration
func (c *Config) Load() error {
	switch c.configExists() {
	case true:
		if err := c.read(); err != nil {
			return err
		}
	default:
		verbose.Debug("Config file not found.")
		c.file = ini.Empty()
		c.banner()
	}

	c.load()
	return nil
}

// Save the configuration
func (c *Config) Save() error {
	var cfg = c.file
	var err = cfg.ReflectFrom(c)

	if err != nil {
		return errwrap.Wrapf("Can't load configuration: {{err}}", err)
	}

	c.updateRemotes()
	c.simplify()

	err = cfg.SaveToIndent(c.Path, "    ")

	if err != nil {
		return errwrap.Wrapf("Can't save configuration: {{err}}", err)
	}

	return nil
}

// Setup the environment
func Setup() error {
	if err := setupGlobal(); err != nil {
		return errwrap.Wrapf("Error setting up global config: {{err}}", err)
	}

	if err := setupContext(); err != nil {
		return errwrap.Wrapf("Error setting up context: {{err}}", err)
	}

	return nil
}

func (c *Config) setDefaults() {
	c.Local = true
	c.LocalPort = 8080
	c.Endpoint = defaults.Endpoint
	c.NotifyUpdates = true
	c.ReleaseChannel = "stable"

	// By design Windows users should see no color unless they enable it
	// Issue #51.
	// https://github.com/wedeploy/cli/issues/51
	if runtime.GOOS == "windows" {
		c.NoColor = true
	}
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

func (c *Config) load() {
	c.setDefaults()

	if err := c.file.MapTo(c); err != nil {
		panic(err)
	}
}

func (c *Config) read() error {
	var err error
	c.file, err = ini.Load(c.Path)

	if err != nil {
		return errwrap.Wrapf("Error reading configuration file: {{err}}\n"+
			"Fix "+c.Path+" by hand or erase it.", err)
	}

	c.readRemotes()
	return nil
}

func parseRemoteSectionName(parsable string) (parsed string, is bool) {
	var r = regexp.MustCompile(`remote \"(.*)\"`)
	var matches = r.FindStringSubmatch(parsable)

	if len(matches) == 2 {
		parsed = matches[1]
		is = true
	}

	return parsed, is
}

func (c *Config) listRemotes() []string {
	var list = []string{}

	for _, k := range c.file.SectionStrings() {
		var key, is = parseRemoteSectionName(k)

		if !is {
			continue
		}

		list = append(list, key)
	}

	return list
}

func (c *Config) getRemote(key string) *ini.Section {
	return c.file.Section(`remote "` + key + `"`)
}

func (c *Config) deleteRemote(key string) {
	c.file.DeleteSection(`remote "` + key + `"`)
}

func (c *Config) readRemotes() {
	c.Remotes = Remotes{
		list: remotesList{},
	}

	for _, k := range c.listRemotes() {
		s := c.getRemote(k)
		u := s.Key("url")
		URLComment := strings.TrimPrefix(u.Comment, "#")
		comment := strings.TrimPrefix(s.Comment, "#")

		c.Remotes.list[k] = RemoteConfig{
			URL:        u.String(),
			URLComment: strings.TrimSpace(URLComment),
			Comment:    strings.TrimSpace(comment),
		}
	}
}

func (c *Config) simplify() {
	var mainSection = c.file.Section("")
	var omitempty = []string{"next_version", "last_update_check"}

	for _, k := range omitempty {
		var key = mainSection.Key(k)
		if key.Value() == "" && key.Comment == "" {
			mainSection.DeleteKey(k)
		}
	}
}

func (c *Config) simplifyRemotes() {
	for _, k := range c.listRemotes() {
		s := c.getRemote(k)
		u := s.Key("url")

		if u.Value() == "" && u.Comment == "" {
			s.DeleteKey("url")
		}

		if len(s.Keys()) == 0 && s.Comment == "" {
			c.deleteRemote(k)
		}
	}
}

func (c *Config) updateRemotes() {
	for _, k := range c.listRemotes() {
		if _, ok := c.Remotes.list[k]; !ok {
			c.deleteRemote(k)
		}
	}

	for k, v := range c.Remotes.list {
		s := c.getRemote(k)
		key := s.Key("url")
		key.SetValue(v.URL)
		key.Comment = v.URLComment
		s.Comment = v.Comment
	}

	c.simplifyRemotes()
}

func (c *Config) banner() {
	c.file.Section("DEFAULT").Comment = `# Configuration file for WeDeploy CLI
# https://wedeploy.io`
}

func setupContext() (err error) {
	Context, err = usercontext.Get()

	if err != nil {
		err = errwrap.Wrapf("Fatal context setup failure: {{err}}", err)
	}

	return err
}

func setupGlobal() error {
	Global = &Config{
		Path: filepath.Join(user.GetHomeDir(), ".we"),
	}

	return Global.Load()
}

// Teardown resets the configuration environment
func Teardown() {
	teardownGlobal()
	teardownContext()
}

func teardownContext() {
	Context = nil
}

func teardownGlobal() {
	Global = nil
}
