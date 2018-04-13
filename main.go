package main

import (
	"log"
	"time"

	"github.com/radisvaliullin/nano-tgbot/tgbot"

	"github.com/radisvaliullin/nano-tgbot/config"
	"github.com/radisvaliullin/nano-tgbot/zlog"
	"go.uber.org/zap"
)

func main() {

	log.Println("Nano-TgBot is start")

	// app config
	conf, err := config.NewAppConfig()
	if err != nil {
		log.Fatal("new app config err - ", err)
	}
	log.Printf("app config: log conf - %+v; bot conf - %+v", conf.Log, conf.Bot)

	// app logger setup
	// don't use global variable and signelton in your app
	// each rule has an exclusion, maybe logger is than exlusion
	// but if you don't want blocking in log operation,
	// you must use self logger for each component of your app
	zlog.SetColorZapGlobalLogger()
	zap.L().Info("zap global logger setup")

	// start bot
	bot := tgbot.NewBot(conf.Bot)
	err = bot.Start()
	if err != nil {
		zap.L().Fatal("tgbot start", zap.Error(err))
	}

	for {
		zap.L().Info("hearbit")
		time.Sleep(time.Second * 10)
	}
}
