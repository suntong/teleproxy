////////////////////////////////////////////////////////////////////////////
// Program: shuttlebot
// Purpose: Telegram Shuttle Bot
// Authors: Tong Sun (c) 2018, All rights reserved
////////////////////////////////////////////////////////////////////////////

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-easygen/cli"
)

////////////////////////////////////////////////////////////////////////////
// forward

func forwardCLI(ctx *cli.Context) error {
	rootArgv = ctx.RootArgv().(*rootT)
	argv := ctx.Argv().(*forwardT)
	fmt.Printf("[forward]:\n  %+v\n  %+v\n  %v\n", rootArgv, argv, ctx.Args())
	Opts.LogLevel, Opts.Version, Opts.Verbose =
		rootArgv.LogLevel, rootArgv.Version, rootArgv.Verbose.Value()
	//return nil

	var cfg Config
	cfg.ChatID = "-" + argv.ChatID[0]
	cfg.Token = argv.Token
	cfg.Template = argv.Template
	cfg.Command = argv.Command

	// Create a new instance of the logger
	lg, err := NewLog(LogConfig{Opts.LogLevel})
	exitOnError(nil, err, "Parse loglevel")

	app := Application{
		Config: &cfg,
		Log:    lg,
	}

	return DoForward(app)
}

func DoForward(app Application) error {

	app.Log.Printf("%s v %s. Telegram Forwarding Shuttle Bot", progname, version)
	app.Log.Print("Copyright (C) 20118, Tong Sun")
	app.Log.Print("Copyright (C) 20118, 2017-18, Alexey Kovrizhkin <ak@elfire.ru>")

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
		app.Log.Printf("info: Got signal %v", sig)
		os.Exit(0)
	}()

	app.Run()

	return nil
}
