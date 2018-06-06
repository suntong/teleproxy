//go:generate go-bindata -pkg $GOPACKAGE -prefix ./ -o bindata.go ./messages.tmpl

package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/jessevdk/go-flags"

	"github.com/LeKovr/go-base/log"
)

// -----------------------------------------------------------------------------

// Flags defines local application flags
type Flags struct {
	ChatID   string `long:"group"    description:"Telegram group ID (without -)"`
	Token    string `long:"token"    description:"Bot token"`
	Template string `long:"template" description:"Message template"`
	Command  string `long:"command"  description:"External command file"`
	Version  bool   `long:"version"  description:"Show version and exit"`
}

// Config defines all of application flags
type Config struct {
	Flags
	Log LogConfig `group:"Logging Options"`
}

// -----------------------------------------------------------------------------

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

func old_main() {

	var cfg Config
	log, _ := setUp(&cfg)

	log.Printf("teleproxy v %s. Telegram proxy bot", Version)
	log.Print("Copyright (C) 2017, Alexey Kovrizhkin <ak@elfire.ru>")

	app := Application{
		Config: &cfg,
		Log:    log,
	}

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
		log.Printf("info: Got signal %v", sig)
		os.Exit(0)
	}()

	app.Run()

	os.Exit(0)
}

// -----------------------------------------------------------------------------

func setUp(cfg *Config) (lg log.Logger, err error) {

	_, err = flags.Parse(cfg)
	if err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
			os.Exit(1) // help printed
		} else {
			os.Exit(2) // error message written already
		}
	}
	if cfg.Version {
		// show version & exit
		fmt.Printf("%s\n%s\n%s", Version, Build, Commit)
		os.Exit(0)
	}

	// use all CPU cores for maximum performance
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Create a new instance of the logger
	lg, err = NewLog(cfg.Log)
	exitOnError(nil, err, "Parse loglevel")

	// group id in config > 0 but we need < 0
	cfg.ChatID = "-" + cfg.ChatID

	return
}

// -----------------------------------------------------------------------------

func exitOnError(lg log.Logger, err error, msg string) {
	if err != nil {
		if lg != nil {
			lg.Printf("error: %s error: %s", msg, err.Error())
		} else {
			fmt.Printf("error: %s error: %s", msg, err.Error())
		}
		os.Exit(1)
	}
}
