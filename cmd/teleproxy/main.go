package main

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"
	"time"

	"github.com/jessevdk/go-flags"

	"github.com/LeKovr/go-base/database"
	"github.com/LeKovr/go-base/logger"
)

// -----------------------------------------------------------------------------

// Flags defines local application flags
type Flags struct {
	Group    int64    `long:"group"    description:"Telegram group ID (without -)"`
	Token    string   `long:"token"    description:"Bot token"`
	Template string   `long:"template" default:"messages.gohtml" description:"Message template"`
	Commands []string `long:"command"  description:"Allowed command(s)"`
	Version  bool     `long:"version"  description:"Show version and exit"`
}

// Config defines all of application flags
type Config struct {
	Flags
	Logger logger.Flags   `group:"Logging Options"`
	DB     database.Flags `group:"Database Options"`
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

func main() {

	var cfg Config
	log, db, _ := setUp(&cfg)
	defer log.Close()

	Program := path.Base(os.Args[0])
	log.Infof("%s v %s. Telegram proxy bot", Program, Version)
	log.Println("Copyright (C) 2017, Alexey Kovrizhkin <ak@elfire.ru>")

	app := Application{
		Log: log,
		DB:  db,
	}

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
		log.Infof("Got signal %v", sig)
		app.Close()
	}()

	app.Run(cfg)

	os.Exit(0)
}

// -----------------------------------------------------------------------------

func setUp(cfg *Config) (log *logger.Log, db *database.DB, err error) {

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
	log, err = logger.New(logger.Dest(cfg.Logger.Dest), logger.Level(cfg.Logger.Level))
	if err != nil {
		panic("Logger init error: " + err.Error())
	}

	// Setup database
	db, err = database.New(cfg.DB.Driver, cfg.DB.Connect, database.Debug(cfg.DB.Debug))
	stopOnError(log, err, "DB init")

	// Sync database
	err = db.Engine.Sync(new(Customer))
	stopOnError(log, err, "DB customer sync")
	err = db.Engine.Sync(new(Record))
	stopOnError(log, err, "DB record sync")

	// group id in config > 0 but we need < 0
	cfg.Group = cfg.Group * -1

	return
}

// -----------------------------------------------------------------------------

// stopOnError used internally for fatal errors checking
func stopOnError(log *logger.Log, err error, info string) {
	if err != nil {
		log.Fatalf("Error with %s: %v", info, err)
	}

}
