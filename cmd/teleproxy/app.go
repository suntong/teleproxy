package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/tucnak/telebot"

	"github.com/LeKovr/go-base/database"
	"github.com/LeKovr/go-base/logger"
)

// -----------------------------------------------------------------------------

// Application holds app.Say
type Application struct {
	DB       *database.DB
	Log      *logger.Log
	bot      *telebot.Bot
	template *template.Template
	messages chan telebot.Message
}

// Say loads message from template and sais it to chat
func (app Application) Say(code string, chat telebot.Recipient, user Customer, text string) {
	vars := struct {
		Text string
		User Customer
	}{
		text,
		user,
	}
	buf := new(bytes.Buffer)
	//err :=
	app.template.ExecuteTemplate(buf, code, vars)

	app.bot.SendMessage(chat, buf.String(), nil)

}

// -----------------------------------------------------------------------------

// Exec runs external command
func (app Application) Exec(cmd string, chat telebot.Recipient) {

	out, err := exec.Command("./" + cmd + ".sh").Output()
	// Записать в логи результат скрипта
	if err != nil {
		app.Log.Warnf("cmd ERROR: %+v (%s)", err, out)
		app.bot.SendMessage(chat, "*Ошибка:* "+err.Error(), &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
	} else {
		app.Log.Warnf("cmd OUT: %s", out)
	}
	app.bot.SendMessage(chat, string(out), nil)
}

// loadUser sets Customer fields from telebot.User
func (u *Customer) loadUser(c telebot.User) {
	u.FirstName = c.FirstName
	u.LastName = c.LastName
	u.Username = c.Username
}

// SetState sets user Disabled state
func (app Application) SetState(state int, chat telebot.Recipient, user Customer) error {

	sql := "update customer set disabled = ? where id = ? and disabled <> ?"
	res, err := app.DB.Engine.Exec(sql, state, user.ID, state)

	s := fmt.Sprintf("%01d", state)
	if err != nil {
		app.Log.Warnf("SQL1 error: %s", err.Error())
		app.Say("errState"+s, chat, user, err.Error())
		return err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		app.Log.Warnf("SQL2 error: %s", err.Error())
		app.Say("errState"+s, chat, user, err.Error())
		return err
	}

	if aff > 0 {
		app.Say("userState"+s, chat, user, "")
	} else {
		app.Say("userStateKeep", chat, user, "")
	}
	return nil
}

// -----------------------------------------------------------------------------

// Close closes message channel
func (app Application) Close() {
	if app.messages != nil {
		close(app.messages)
	}
}

// -----------------------------------------------------------------------------

// Run does the deal
func (app *Application) Run(cfg Config) {

	bot, err := telebot.NewBot(cfg.Token)
	stopOnError(app.Log, err, "Bot init")
	app.bot = bot

	tmpl, err := template.New("").ParseFiles(cfg.Template)
	stopOnError(app.Log, err, "Template load")
	app.template = tmpl

	app.messages = make(chan telebot.Message)

	bot.Listen(app.messages, 1*time.Second)
	app.Log.Printf("Connected bot %s", bot.Identity.Username)

	group := telebot.Chat{ID: cfg.Group, Type: "group"}
	engine := app.DB.Engine

	for message := range app.messages {
		inChat := message.Chat.ID == cfg.Group
		app.Log.Debugf("Sender: %+v", message.Sender)
		app.Log.Debugf("%s: %s", message.Chat.Title, message.Text)
		sender := Customer{ID: int64(message.Sender.ID)}

		has, _ := engine.Get(&sender) // TODO: err
		if !has {
			// new customer or op
			sender.loadUser(message.Sender)

			if _, err := engine.Insert(&sender); err != nil {
				app.Log.Errorf("User add error: %+v", err)
			}
			//has, err :=
			engine.Get(&sender)
			if !inChat {
				app.Say("info", group, sender, ".new")
			}
		}

		if message.Text == "/hi" {
			// Say Hi to any user
			app.Say("hello", message.Chat, sender, "")

		} else if inChat { // && strings.HasPrefix(message.Text, "/") {
			// group bot commands, always started from /

			if message.Text == "/help" {
				// Operator needs help
				app.Say("helpOp", message.Chat, sender, "")
				continue
			}
			// Customer related commands

			// split customer Code & rest
			reply := strings.SplitN(strings.TrimPrefix(message.Text, "/"), " ", 2)
			if len(reply) == 1 {
				// run internal command
				if stringExists(cfg.Commands, reply[0]) {
					app.Say("cmdRequest", message.Chat, sender, reply[0])
					go app.Exec(reply[0], message.Chat)
				} else {
					app.Say("errNoCmd", message.Chat, sender, reply[0])
				}
				continue
			}

			c, err := strconv.ParseUint(reply[0], 10, 64)
			if err != nil {
				app.Say("errNoDigit", message.Chat, sender, reply[0])
				continue
			}

			var user = Customer{Code: c}
			has, _ := engine.Get(&user) // TODO: err

			if !has {
				// customer not found
				app.Say("errNoUser", message.Chat, sender, reply[0])

			} else if len(reply) == 2 {
				// given customer code & something
				app.Log.Debugf("Customer: %+v", user)
				switch reply[1] {
				case "=":
					// customer info requested
					app.Say("info", message.Chat, user, "")
					continue
				case "-":
					// lock user
					if app.SetState(1, message.Chat, user) != nil {
						continue
					}
				case "+":
					// unlock user
					if app.SetState(0, message.Chat, user) != nil {
						continue
					}
				default:
					// forward reply to customer
					chat := telebot.Chat{ID: user.ID, Type: "private"}
					bot.SendMessage(chat, reply[1], nil)
				}

				// save log
				rec := Record{ID: user.ID, IDFrom: sender.ID, Message: reply[1]}
				if _, err := engine.Insert(&rec); err != nil {
					app.Log.Errorf("Record add error: %+v", err)
				}

			}
		} else if !inChat {
			// Message from customer

			if message.Text == "/start" {
				// bot started
				app.Say("welcome", message.Chat, sender, "")
				continue
			}

			// other message
			if sender.Disabled < 2 {
				rec := Record{ID: sender.ID, Message: message.Text}
				if _, err := engine.Insert(&rec); err != nil {
					app.Log.Errorf("Record add error: %+v", err)
				}

				if sender.Disabled < 1 {
					app.Say("message", group, sender, message.Text)
				} else {
					app.Say("userLocked", message.Chat, sender, "")
				}
			}
		}
		time.Sleep(time.Second) // wait 1 sec always
	}
	app.Log.Info("Exiting")

}

// -----------------------------------------------------------------------------

// Check if str exists in strings slice
func stringExists(strings []string, str string) bool {
	if len(strings) > 0 {
		for _, s := range strings {
			if str == s {
				return true
			}
		}
	}
	return false
}
