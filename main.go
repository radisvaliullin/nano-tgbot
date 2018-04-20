package main

import (
	"log"
	"os"
	"os/signal"
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
	zap.L().Info("tgbot started")
	if err != nil {
		zap.L().Fatal("tgbot start", zap.Error(err))
	}

	// heartbit
	go func() {
		for {
			zap.L().Info("heartbit")
			time.Sleep(time.Second * 10)
		}
	}()

	// handle stop signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for sig := range c {
		zap.L().Info("catch os signal", zap.Any("sig", sig))
		bot.Stop()
		zap.L().Info("wait stop")
		bot.WaitStop()
		return
	}
}
