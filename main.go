package main

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"text/template"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/tucnak/telebot"

	"github.com/LeKovr/go-base/database"
	"github.com/LeKovr/go-base/logger"
)

// -----------------------------------------------------------------------------

// Flags defines local application flags
type Flags struct {
	Group    int64  `long:"group"    description:"Telegram group ID (without -)"`
	Token    string `long:"token"    description:"Bot token"`
	Template string `long:"template" default:"messages.gohtml" description:"Message template"`
	Version  bool   `long:"version"  description:"Show version and exit"`
}

// Config defines all of application flags
type Config struct {
	Flags
	log logger.Flags
	db  database.Flags
}

// Customer - таблица журнала сообщений
type Customer struct {
	Code                          uint64 `xorm:"pk autoincr"`
	ID                            int64  `xorm:"id unique not null"`
	FirstName, LastName, Username string
	Stamp                         time.Time `xorm:"not null created"`
	Disabled                      int       `xorm:"not null"` // 1 - log only, 2 - full
}

// Record - таблица журнала сообщений
type Record struct {
	Stamp   time.Time `xorm:"pk created"`
	ID      int64     `xorm:"id pk not null"`
	IDFrom  int64     `xorm:"id_from"`
	Message string
}

// -----------------------------------------------------------------------------

func main() {

	var cfg Config
	log, db, _ := setUp(&cfg)
	defer log.Close()

	Program := path.Base(os.Args[0])
	log.Infof("%s v %s. Telegram proxy bot", Program, Version)
	log.Println("Copyright (C) 2016, Alexey Kovrizhkin <ak@elfire.ru>")

	run(cfg, log, db)

	os.Exit(0)
}

// -----------------------------------------------------------------------------

func makeConfig(cfg *Config) *flags.Parser {
	p := flags.NewParser(nil, flags.Default)
	_, err := p.AddGroup("Application Options", "", cfg)
	panicIfError(err) // check Flags parse error

	_, err = p.AddGroup("Logging Options", "", &cfg.log)
	panicIfError(err) // check Flags parse error

	_, err = p.AddGroup("Database Options", "", &cfg.db)
	panicIfError(err)

	return p
}

// -----------------------------------------------------------------------------

func setUp(cfg *Config) (log *logger.Log, db *database.DB, err error) {

	p := makeConfig(cfg)

	_, err = p.Parse()
	if err != nil {
		os.Exit(1) // error message written already
	}
	if cfg.Version {
		// show version & exit
		fmt.Printf("%s\n%s\n%s", Version, Build, Commit)
		os.Exit(0)
	}

	// use all CPU cores for maximum performance
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Create a new instance of the logger
	log, err = logger.New(logger.Dest(cfg.log.Dest), logger.Level(cfg.log.Level))
	panicIfError(err) // check Flags parse error

	// Setup database
	db, err = database.New(cfg.db.Driver, cfg.db.Connect, database.Debug(cfg.db.Debug))
	panicIfError(err) // check Flags parse error

	// Sync database
	err = db.Engine.Sync(new(Customer))
	if err != nil {
		log.Fatalf("DB sync error: %v", err)
	}
	err = db.Engine.Sync(new(Record))
	if err != nil {
		log.Fatalf("DB sync error: %v", err)
	}

	// option without "-"
	cfg.Group = cfg.Group * -1

	return
}

// -----------------------------------------------------------------------------

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

// Application holds app.Say
type Application struct {
	Bot      *telebot.Bot
	DB       *database.DB
	Template *template.Template
	Log      *logger.Log
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
	app.Template.ExecuteTemplate(buf, code, vars)

	app.Bot.SendMessage(chat, buf.String(), nil)

}

func (u *Customer) loadUser(c telebot.User) {
	u.FirstName = c.FirstName
	u.LastName = c.LastName
	u.Username = c.Username
}

// SetState sets user Disabled state
func (app Application) SetState(state int, chat telebot.Recipient, user Customer) error {
	//	affected, err := app.DB.Engine.Id(user.Code).Cols("disabled").Update(Customer{Disabled: state})

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

func run(cfg Config, log *logger.Log, db *database.DB) {

	bot, err := telebot.NewBot(cfg.Token)
	panicIfError(err)

	tmpl, err := template.New("").ParseFiles(cfg.Template)
	panicIfError(err)

	app := Application{
		Log:      log,
		Bot:      bot,
		Template: tmpl,
		DB:       db,
	}

	engine := db.Engine

	messages := make(chan telebot.Message)

	bot.Listen(messages, 1*time.Second)
	log.Printf("Connected bot %s", bot.Identity.Username)

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
		log.Infof("Got signal %v", sig)
		close(messages)
	}()

	group := telebot.Chat{ID: cfg.Group, Type: "group"}

	for message := range messages {
		inChat := message.Chat.ID == cfg.Group
		log.Debugf("Sender: %+v", message.Sender)
		log.Debugf("%s: %s", message.Chat.Title, message.Text)
		sender := Customer{ID: int64(message.Sender.ID)}

		has, _ := engine.Get(&sender) // TODO: err
		if !has {
			// new customer or op
			sender.loadUser(message.Sender)

			if _, err := db.Engine.Insert(&sender); err != nil {
				log.Errorf("User add error: %+v", err)
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
				log.Debugf("Customer: %+v", user)
				if reply[1] == "=" {
					// customer info requested
					app.Say("info", message.Chat, user, "")
					continue
				} else if reply[1] == "" {
					continue
				}
				if reply[1] == "-" {
					// lock user
					if app.SetState(1, message.Chat, user) != nil {
						continue
					}

				} else if reply[1] == "+" {
					// unlock user
					if app.SetState(0, message.Chat, user) != nil {
						continue
					}

				} else {
					// forward reply to customer
					chat := telebot.Chat{ID: user.ID, Type: "private"}
					bot.SendMessage(chat, reply[1], nil)
				}

				// save log
				rec := Record{ID: user.ID, IDFrom: sender.ID, Message: reply[1]}
				if _, err := db.Engine.Insert(&rec); err != nil {
					log.Errorf("Record add error: %+v", err)
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
				if _, err := db.Engine.Insert(&rec); err != nil {
					log.Errorf("Record add error: %+v", err)
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
	log.Info("Exiting")

}
