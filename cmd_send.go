////////////////////////////////////////////////////////////////////////////
// Program: shuttlebot
// Purpose: Telegram Shuttle Bot
// Authors: Tong Sun (c) 2018, All rights reserved
////////////////////////////////////////////////////////////////////////////

package main

import (
	"fmt"

	"github.com/go-easygen/cli"
)

////////////////////////////////////////////////////////////////////////////
// send

func sendCLI(ctx *cli.Context) error {
	rootArgv = ctx.RootArgv().(*rootT)
	argv := ctx.Argv().(*sendT)
	fmt.Printf("[send]:\n  %+v\n  %+v\n  %v\n", rootArgv, argv, ctx.Args())
	Opts.LogLevel, Opts.Version, Opts.Verbose =
		rootArgv.LogLevel, rootArgv.Version, rootArgv.Verbose.Value()
	return nil
}
