package config

import (
	"go/build"
	"log"
	"os"
	"path"

	"gopkg.in/yaml.v2"

	"github.com/radisvaliullin/nano-tgbot/tgbot"
)

const (
	pathProjRoot = "github.com/radisvaliullin/nano-tgbot"
	pathExamConf = "config/example.config.yml"

	// ENV - default enviroment variables
	envAppBotToken = "APP_BOT_TOKEN"
	envAppConfPath = "APP_CONFPATH"
)

// LogConf -
type LogConf struct {
	Level string
}

// AppConfig - all app configs
type AppConfig struct {
	Log *LogConf
	Bot *tgbot.BotConf
}

// NewAppConfig - new app conf
func NewAppConfig() (*AppConfig, error) {

	c := &AppConfig{}

	// config file path
	// try get conf file path from env, else use default example conf
	confPath := os.Getenv(envAppConfPath)
	if confPath == "" {
		confPath = pathExamConf
	}
	fullConfPath := path.Join(getGOPATH(), "src", pathProjRoot, confPath)

	// open config file
	f, err := os.Open(fullConfPath)
	if err != nil {
		log.Println("config: conf file open err - ", err)
		return nil, err
	}

	// parse config
	// we use yaml lib
	// (exist many lib for work with config, example - https://github.com/spf13/viper)
	err = yaml.NewDecoder(f).Decode(c)
	if err != nil {
		log.Print("config: yaml conf file decode err - ", err)
		return nil, err
	}

	// Get bot token from env
	c.Bot.Token = os.Getenv(envAppBotToken)

	return c, nil
}

func getGOPATH() string {
	p := os.Getenv("GOPATH")
	if p == "" {
		p = build.Default.GOPATH
	}
	return p
}
