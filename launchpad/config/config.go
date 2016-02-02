package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/launchpad-project/cli/launchpad/util"
	"github.com/spf13/viper"
)

type Store struct {
	Name             string
	Path             string
	ConfigurableKeys map[string]bool
	Data             *viper.Viper
}

var ErrConfigKeyNotFound = errors.New("key not found")
var ErrConfigKeyNotConfigurable = errors.New("key not configurable")

var global = &Store{
	Name: "global",
	Path: util.GetUserHomeDir() + "/.launchpad.json",
	ConfigurableKeys: map[string]bool{
		"username": true,
		"password": true,
		"endpoint": true,
	},
}

var app = &Store{
	Name: "app",
	Path: "./launchpad.json",
	ConfigurableKeys: map[string]bool{
		"name": true,
	},
}

var Stores = map[string]*Store{
	"global": global,
	"app":    app,
}

func (s *Store) Save() {
	bin, err := json.MarshalIndent(s.toMap(), "", "    ")

	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(s.Path, bin, 0644); err != nil {
		panic(err)
	}
}

func (s *Store) toMap() interface{} {
	var settings map[string]interface{}

	switch s.Name {
	case "global":
		settings = map[string]interface{}{
			"username": s.Data.GetString("username"),
			"password": s.Data.GetString("password"),
			"endpoint": s.Data.GetString("endpoint"),
		}
	case "app":
		settings = map[string]interface{}{
			"name": s.Data.GetString("name"),
		}
	}

	return settings
}

func (s *Store) Get(key string) (string, error) {
	if !s.Data.IsSet(key) || len(s.Data.GetStringMapString(key)) != 0 {
		return "", ErrConfigKeyNotFound
	}

	var value = s.Data.Get(key)

	return fmt.Sprint(value), nil
}

func (s *Store) Set(key, value string) error {
	if !s.ConfigurableKeys[key] {
		return ErrConfigKeyNotConfigurable
	}

	s.Data.Set(key, value)
	return nil
}

func (s *Store) SetAndSave(key, value string) {
	if err := s.Set(key, value); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	s.Save()
}

func setupGlobalConfig() {
	var Data = viper.New()
	Data.SetConfigName(".launchpad")
	Data.SetConfigType("json")
	Data.AddConfigPath("$HOME")

	if err := Data.ReadInConfig(); err != nil {
		print("Fatal error reading global configuration file (.launchpad.json).")
		panic(err)
	}

	global.Data = Data
}

func setupAppConfig() {
	var Data = viper.New()
	Data.SetConfigName("launchpad")
	Data.SetConfigType("json")
	Data.AddConfigPath(".")

	if err := Data.ReadInConfig(); err != nil && !os.IsNotExist(err) {
		print("Fatal error reading project configuration file (launchpad.json).")
		panic(err)
	}

	app.Data = Data
}

func Setup() {
	setupGlobalConfig()
	setupAppConfig()
}
