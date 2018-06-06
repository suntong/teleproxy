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
	return DoForward()
}

func DoForward() error {
	var cfg Config
	log, _ := setUp(&cfg)

	log.Printf("%s v %s. Telegram Forwarding Shuttle Bot", progname, version)
	log.Print("Copyright (C) 20118, Tong Sun\n\t2017-18, Alexey Kovrizhkin <ak@elfire.ru>")

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

	return nil
}
