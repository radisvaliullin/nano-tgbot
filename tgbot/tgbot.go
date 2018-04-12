package tgbot

import (
	"time"
)

// BotConf - out nano bot configs
type BotConf struct {
	Token string
}

// Bot - our nano telegram bot (useless)
type Bot struct {
	conf *BotConf
}

// NewBot - return new bot object
// if need we also can return error value
func NewBot(conf *BotConf) *Bot {
	if conf == nil {
		conf = &BotConf{}
	}
	return &Bot{conf}
}

// Start - starts out bot, and if need init addition components
// prefer to reciver use one character name
func (b *Bot) Start() error {

	// run
	go b.run()

	return nil
}

// run - not public method for run bot in goroutine
func (b *Bot) run() {

	// do something
	for {
		time.Sleep(time.Second)
	}
}
