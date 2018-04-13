package tgbot

import (
	"gopkg.in/telegram-bot-api.v4"
)

// Dispatcher - dispatch messages by users
type Dispatcher struct {
	users map[int]*User

	updates  chan tgbotapi.Update
	userResp chan<- UserResp
}

// NewDispatcher -
func NewDispatcher(userResp chan<- UserResp) *Dispatcher {
	d := &Dispatcher{
		users:    make(map[int]*User),
		updates:  make(chan tgbotapi.Update, 1000),
		userResp: userResp,
	}
	return d
}

// Start -
func (d *Dispatcher) Start() {
	go d.run()
}

// run
func (d *Dispatcher) run() {

	for {
		select {
		case upd := <-d.updates:
			user, ok := d.users[upd.Message.From.ID]
			if !ok {
				user = NewUser(d.userResp)
				d.users[upd.Message.From.ID] = user
				user.Start()
			}
			user.GetUserUpdatesChan() <- upd
		}
	}
}

// GetBotUpdatesChan -
func (d *Dispatcher) GetBotUpdatesChan() chan<- tgbotapi.Update {
	return d.updates
}
