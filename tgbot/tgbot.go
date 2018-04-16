package tgbot

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gopkg.in/telegram-bot-api.v4"
)

// BotConf - out nano bot configs
type BotConf struct {
	// keep token only env
	Token string

	//
	DefaultMessage string
}

// Bot - our nano telegram bot (useless)
type Bot struct {
	conf *BotConf

	botAPI     *tgbotapi.BotAPI
	updChan    chan tgbotapi.Update
	updStopSig chan struct{}
	updHdlDone chan struct{}
	runStopSig chan struct{}
	runDone	   chan struct{}

	userResp chan UserResp

	disp *Dispatcher
}

// NewBot - return new bot object
// if need we also can return error value
func NewBot(conf *BotConf) *Bot {
	if conf == nil {
		conf = &BotConf{}
	}
	b := &Bot{
		conf:     conf,
		userResp: make(chan UserResp, 1000),
	}
	return b
}

// Start - starts out bot, and if need init addition components
// prefer to reciver use one character name
func (b *Bot) Start() error {

	// init bot api
	bot, err := tgbotapi.NewBotAPI(b.conf.Token)
	if err != nil {
		zap.L().Error("tgbot: new bot api", zap.Error(err))
		return err
	}
	b.botAPI = bot

	// updates (income messages from users)
	b.updChan = make(chan tgbotapi.Update, b.botAPI.Buffer)
	b.updStopSig = make(chan struct{})
	b.updHdlDone = make(chan struct{})
	b.runStopSig = make(chan struct{})
	b.runDone	 = make(chan struct{})

	// start dispatcher
	dsp := NewDispatcher(b.userResp)
	dsp.Start()
	b.disp = dsp

	// run bot updates
	go b.updates()

	// run, handle updates
	go b.run()

	return nil
}

// Stop -
func (b *Bot) Stop() {
	b.runStopSig <- struct{}{}
}

// WaitStop -
func (b *Bot) WaitStop() {
	<-b.updHdlDone
	<-b.runDone
}

// run - not public method for run bot in goroutine
func (b *Bot) run() {
	// defer done signal
	defer func() { b.runDone <- struct{}{} }()

	// do something
	for {
		select {
		case upd := <-b.updChan:
			if upd.Message == nil {
				continue
			}
			b.disp.GetBotUpdatesChan() <- upd

		case ur := <-b.userResp:
			// resp message
			msg := tgbotapi.NewMessage(ur.chatID, fmt.Sprintf("%v - %s", ur.userID, b.conf.DefaultMessage))

			_, err := b.botAPI.Send(msg)
			if err != nil {
				zap.L().Error("tgbot: send resp message", zap.Error(err))
			}
		case <-b.runStopSig:
			b.updStopSig <- struct{}{}
			b.runDone <- struct{}{}
			close(b.userResp)
			return
		}
	}
}

//
func (b *Bot) updates() {

	// updates config
	uc := tgbotapi.NewUpdate(0)
	uc.Timeout = 10

	for {
		select {
		case <-b.updStopSig:
			close(b.updChan)
			b.updHdlDone <- struct{}{}
			return

		default:
			updates, err := b.botAPI.GetUpdates(uc)
			if err != nil {
				zap.L().Error("tgbot: get updates", zap.Error(err))
				zap.L().Info("Failed to get updates, retrying in 3 seconds...")
				time.Sleep(time.Second * 3)
				continue
			}

			for _, update := range updates {
				if update.UpdateID >= uc.Offset {
					uc.Offset = update.UpdateID + 1
					b.updChan <- update
				}
			}
		}
	}
}
