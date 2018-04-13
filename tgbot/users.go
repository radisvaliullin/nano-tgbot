package tgbot

import (
	"gopkg.in/telegram-bot-api.v4"
)

// User -
type User struct {
	updates  chan tgbotapi.Update
	userResp chan<- UserResp
}

// NewUser -
func NewUser(userResp chan<- UserResp) *User {
	u := &User{
		updates:  make(chan tgbotapi.Update, 100),
		userResp: userResp,
	}
	return u
}

// Start -
func (u *User) Start() {
	go u.run()
}

//
func (u *User) run() {
	for upd := range u.updates {
		// send response to user
		ur := UserResp{
			userID: upd.Message.From.ID,
			chatID: upd.Message.Chat.ID,
		}
		u.userResp <- ur
	}
}

// GetUserUpdatesChan -
func (u *User) GetUserUpdatesChan() chan<- tgbotapi.Update {
	return u.updates
}
