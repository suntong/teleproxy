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

	log "github.com/Sirupsen/logrus"
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

	return DoForward(cfg)
}

func DoForward(cfg Config) error {

	log.Printf("%s v %s. Telegram Forwarding Shuttle Bot", progname, version)
	log.Print("Copyright (C) 20118, Tong Sun")
	log.Print("Copyright (C) 20118, 2017-18, Alexey Kovrizhkin <ak@elfire.ru>")

	app := Application{
		Config: &cfg,
		Log:    log.New(),
	}

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
		log.Printf("info: Got signal %v", sig)
		os.Exit(0)
	}()

	app.Run()

	return nil
}
