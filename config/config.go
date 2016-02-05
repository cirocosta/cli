package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/launchpad-project/cli/util"
)

// Store is the structure for the config object
type Store struct {
	Name             string
	Path             string
	ConfigurableKeys map[string]bool
	Data             map[string]interface{}
}

// ErrConfigKeyNotFound is used when a key is not found on the config
var ErrConfigKeyNotFound = errors.New("key not found")

// ErrConfigKeyNotConfigurable is used when a key is not configurable
var ErrConfigKeyNotConfigurable = errors.New("key not configurable")

var global = &Store{
	Name: "global",
	Path: util.GetUserHomeDir() + "/.launchpad.json",
	// only string values should be configurable
	ConfigurableKeys: map[string]bool{
		"username": true,
		"password": true,
		"endpoint": false,
	},
}

var app = &Store{
	Name: "app",
	Path: "./launchpad.json",
	ConfigurableKeys: map[string]bool{
		"name":        true,
		"description": true,
		"domain":      true,
	},
}

// Stores sets the map of available config stores
var Stores = map[string]*Store{
	"global": global,
	"app":    app,
}

// Save saves the config file
func (s *Store) Save() {
	bin, err := json.MarshalIndent(s.Data, "", "    ")

	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(s.Path, bin, 0644); err != nil {
		panic(err)
	}
}

// Init reads the config file
func (s *Store) Init() {
	content, err := ioutil.ReadFile(s.Path)

	if err != nil && !os.IsNotExist(err) {
		print("Fatal error reading global configuration file " + s.Path)
		panic(err)
	}

	json.Unmarshal(content, &s.Data)
}

// GetString gets a string from the config object
func (s *Store) GetString(key string) (string, error) {
	var keyPath = strings.Split(key, ".")
	var parent = s.Data

	for pos, subPath := range keyPath {
		if _, exists := parent[subPath]; !exists {
			return "", ErrConfigKeyNotFound
		}

		if pos != len(keyPath)-1 {
			parent = parent[subPath].(map[string]interface{})
			continue
		}

		switch parent[subPath].(type) {
		case nil:
			return "null", nil
		case string, int, int64, float64:
			return fmt.Sprintf("%v", parent[subPath]), nil
		default:
			return "", ErrConfigKeyNotFound
		}
	}

	return "", ErrConfigKeyNotFound
}

// Get gets a string from the config object and panics when not found
func (s *Store) Get(key string) string {
	value, err := s.GetString(key)

	if err != nil {
		panic(err)
	}

	return value
}

// GetInterface gets an interface from the config object
func (s *Store) GetInterface(key string) (interface{}, error) {
	var keyPath = strings.Split(key, ".")
	var parent = s.Data

	for pos, subPath := range keyPath {
		_, exists := parent[subPath]

		if !exists {
			return "", ErrConfigKeyNotFound
		}

		if pos != len(keyPath)-1 {
			parent = parent[subPath].(map[string]interface{})
			continue
		}

		return parent[subPath], nil
	}

	return "", ErrConfigKeyNotFound
}

// Set sets the value for a given key
func (s *Store) Set(key, value string) error {
	if s.Data == nil {
		s.Data = make(map[string]interface{})
	}

	var keyPath = strings.Split(key, ".")
	var parent = s.Data

	for pos, subPath := range keyPath {
		if pos == len(keyPath)-1 {
			parent[subPath] = value
			continue
		}

		switch parent[subPath].(type) {
		case map[string]interface{}:
		default:
			parent[subPath] = make(map[string]interface{})
		}

		parent = parent[subPath].(map[string]interface{})
	}

	return nil
}

// SetPublicKey sets the value for a given key
// (or fail if it the key is not public)
func (s *Store) SetPublicKey(key, value string) error {
	if !s.ConfigurableKeys[key] {
		return ErrConfigKeyNotConfigurable
	}

	return s.Set(key, value)
}

// SetAndSave sets the value for a given key and save config
func (s *Store) SetAndSave(key, value string) {
	if err := s.Set(key, value); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	s.Save()
}

// SetAndSavePublicKey sets the value for a given key and save config
// (or fail if it the key is not public)
func (s *Store) SetAndSavePublicKey(key, value string) {
	if err := s.SetPublicKey(key, value); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	s.Save()
}

// Setup the global and app configs.
func Setup() {
	global.Init()
	app.Init()
}
