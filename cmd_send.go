////////////////////////////////////////////////////////////////////////////
// Program: shuttlebot
// Purpose: Telegram Shuttle Bot
// Authors: Tong Sun (c) 2018, All rights reserved
////////////////////////////////////////////////////////////////////////////

package main

import (
	"fmt"
	"time"

	"github.com/go-easygen/cli"

	tb "gopkg.in/tucnak/telebot.v2"
)

////////////////////////////////////////////////////////////////////////////
// send

func sendCLI(ctx *cli.Context) error {
	rootArgv = ctx.RootArgv().(*rootT)
	argv := ctx.Argv().(*sendT)
	//fmt.Printf("[send]:\n  %+v\n  %+v\n  %v\n", rootArgv, argv, ctx.Args())
	Opts.LogLevel, Opts.Version, Opts.Verbose =
		rootArgv.LogLevel, rootArgv.Version, rootArgv.Verbose.Value()

	// Create a new instance of the logger
	lg, err := NewLog(LogConfig{Opts.LogLevel})
	exitOnError(nil, err, "Parse loglevel")

	var cfg Config
	cfg.ChatID = argv.ChatID
	cfg.Token = argv.Token

	app := Application{
		Config: &cfg,
		Log:    lg,
	}

	return DoSend(app, argv.File)
}

func DoSend(app Application, theFile string) error {

	fmt.Printf("%s v %s. Telegram File sending Shuttle Bot\n", progname, version)
	fmt.Println("Copyright (C) 2018, Tong Sun")

	app.Log.Printf("info: Connecting to Telegram...")
	bot, err := tb.NewBot(tb.Settings{
		Token:  app.Config.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	exitOnError(app.Log, err, "Bot init")
	app.bot = bot

	app.ChatInit("info: Sending file to:")

	// https://github.com/tucnak/telebot/tree/v2#files
	p := &tb.Photo{File: tb.FromDisk(theFile)}

	for _, chat := range app.Chat {
		app.bot.Send(chat, p)
	}

	return nil
}
