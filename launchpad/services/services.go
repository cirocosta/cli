package services

import (
	"encoding/json"
	"fmt"

	"github.com/launchpad-project/cli/launchpad/client"
	"github.com/launchpad-project/cli/launchpad/config"
)

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

type Pod struct {
	PodConfig PodConfig `json:"config"`
	Name      string    `json:"name"`
	Time      int64     `json:"time"`
}

var globalConfig = config.Stores["global"]

func GetPods() {
	var address = globalConfig.Data.GetString("endpoint") + "/_admin/pods"
	var username = globalConfig.Data.GetString("username")
	var password = globalConfig.Data.GetString("password")
	var l = client.Url(address, nil)

	l.Auth(username, password)
	l.Get()

	pods := *new([]Pod)

	if err := l.ResponseJson(&pods); err != nil {
		panic(err)
	}

	for _, pod := range pods {
		fmt.Println(pod.Name)
	}
}
