package services

import (
	"encoding/json"
	"fmt"

	"github.com/launchpad-project/cli/config"
	"github.com/launchpad-project/cli/launchpad"
)

// PodConfig structure
type PodConfig struct {
	Visibility     bool              `json:"visibility"`
	AssetsPath     string            `json:"assetsPath"`
	BasePath       string            `json:"basePath"`
	ConfigPath     string            `json:"configPath"`
	JavascriptPath string            `json:"javascriptPath"`
	LibPath        []string          `json:"libPath"`
	Runtimes       []json.RawMessage `json:"runtimes"`
	SocketIo       bool              `json:"socketio"`
	WebPath        string            `json:"webPath"`
}

// Pod structure
type Pod struct {
	PodConfig PodConfig `json:"config"`
	Name      string    `json:"name"`
	Time      int64     `json:"time"`
}

var globalConfig = config.Stores["global"]

// GetPods lists the pods
func GetPods() {
	var address = globalConfig.Get("endpoint") + "/_admin/pods"
	var username = globalConfig.Get("username")
	var password = globalConfig.Get("password")
	var l, err = launchpad.URL(address, nil)

	if err != nil {
		panic(err)
	}

	l.Auth(username, password)
	l.Get()

	pods := *new([]Pod)

	if err := l.DecodeJSON(&pods); err != nil {
		panic(err)
	}

	for _, pod := range pods {
		fmt.Println(pod.Name)
	}
}
