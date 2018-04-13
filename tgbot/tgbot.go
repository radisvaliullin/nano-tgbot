package tgbot

import (
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
}

// NewBot - return new bot object
// if need we also can return error value
func NewBot(conf *BotConf) *Bot {
	if conf == nil {
		conf = &BotConf{}
	}
	return &Bot{conf: conf}
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

	// run bot updates
	go b.updates()

	// run, handle updates
	go b.run()

	return nil
}

// Stop -
func (b *Bot) Stop() {
	b.updStopSig <- struct{}{}
}

// WaitStop -
func (b *Bot) WaitStop() {
	<-b.updHdlDone
}

// run - not public method for run bot in goroutine
func (b *Bot) run() {
	// defer done signal
	defer func() { b.updHdlDone <- struct{}{} }()

	// do something
	for upd := range b.updChan {

		if upd.Message == nil {
			continue
		}

		// resp message
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, b.conf.DefaultMessage)
		msg.ReplyToMessageID = upd.Message.MessageID

		_, err := b.botAPI.Send(msg)
		if err != nil {
			zap.L().Error("tgbot: send resp message", zap.Error(err))
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
