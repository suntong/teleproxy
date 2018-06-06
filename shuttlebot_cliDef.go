////////////////////////////////////////////////////////////////////////////
// Program: shuttlebot
// Purpose: Telegram Shuttle Bot
// Authors: Tong Sun (c) 2018, All rights reserved
////////////////////////////////////////////////////////////////////////////

package main

import (
	"github.com/go-easygen/cli"
)

////////////////////////////////////////////////////////////////////////////
// Constant and data type/structure definitions

//==========================================================================
// shuttlebot

type rootT struct {
	cli.Helper
	LogLevel string      `cli:"log_level" usage:"logging level"`
	Verbose  cli.Counter `cli:"v,verbose" usage:"Verbose mode (Multiple -v options increase the verbosity.)"`
	Version  bool        `cli:"V,version" usage:"Show version and exit"`
}

var root = &cli.Command{
	Name:   "shuttlebot",
	Desc:   "Telegram Shuttle Bot\nVersion " + version + " built on " + date,
	Text:   "Toolkit to transfer things for Telegram",
	Global: true,
	Argv:   func() interface{} { return new(rootT) },
	Fn:     shuttlebot,

	NumArg: cli.AtLeast(1),
}

// Template for main starts here
////////////////////////////////////////////////////////////////////////////
// Constant and data type/structure definitions

// The OptsT type defines all the configurable options from cli.
//  type OptsT struct {
//  	LogLevel	string
//  	Verbose	cli.Counter
//  	Version	bool
//  	Verbose int
//  }

////////////////////////////////////////////////////////////////////////////
// Global variables definitions

//  var (
//          progname  = "shuttlebot"
//          version   = "0.1.0"
//          date = "2018-06-06"

//  	rootArgv *rootT
//  	// Opts store all the configurable options
//  	Opts OptsT
//  )

////////////////////////////////////////////////////////////////////////////
// Function definitions

// Function main
//  func main() {
//  	cli.SetUsageStyle(cli.DenseNormalStyle) // left-right, for up-down, use ManualStyle
//  	//NOTE: You can set any writer implements io.Writer
//  	// default writer is os.Stdout
//  	if err := cli.Root(root,
//  		cli.Tree(forwardDef),
//  		cli.Tree(sendDef)).Run(os.Args[1:]); err != nil {
//  		fmt.Fprintln(os.Stderr, err)
//  	}
//  	fmt.Println("")
//  }

// Template for main dispatcher starts here
//==========================================================================
// Main dispatcher

//  func shuttlebot(ctx *cli.Context) error {
//  	ctx.JSON(ctx.RootArgv())
//  	ctx.JSON(ctx.Argv())
//  	fmt.Println()

//  	return nil
//  }

// Template for CLI handling starts here

////////////////////////////////////////////////////////////////////////////
// forward

//  func forwardCLI(ctx *cli.Context) error {
//  	rootArgv = ctx.RootArgv().(*rootT)
//  	argv := ctx.Argv().(*forwardT)
//  	fmt.Printf("[forward]:\n  %+v\n  %+v\n  %v\n", rootArgv, argv, ctx.Args())
//  	Opts.LogLevel, Opts.Verbose, Opts.Version, Opts.Verbose =
//  		rootArgv.LogLevel, rootArgv.Verbose, rootArgv.Version, rootArgv.Verbose.Value()
//  	return nil
//  	//return DoForward()
//  }
//
//  func DoForward() error {
//  	return nil
//  }

type forwardT struct {
	Self      *forwardT `cli:"c,config" usage:"config file" json:"-" parser:"jsonfile" dft:"cfg_forward.json"`
	Token     string    `cli:"t,token" usage:"The telegram bot token (mandatory)" dft:"$SHUTTLEBOT_TOKEN"`
	ChatID    []string  `cli:"i,id" usage:"The telegram ChatID(s) (without -) to forward to (mandatory)" dft:"$SHUTTLEBOT_CID"`
	Template  string    `cli:"template" usage:"Message template" dft:"messages.en.tmpl"`
	Command   string    `cli:"command" usage:"External command file" dft:"./commands.sh"`
	Daemonize bool      `cli:"D,daemonize" usage:"daemonize the service"`
}

var forwardDef = &cli.Command{
	Name: "forward",
	Desc: "forwards telegram messages to designated ChatID(s)",
	Text: "Usage:\n  shuttlebot forward --log_level debug --id $GROUP --token $TOKEN --template $TEMPLATE --command ./commands.sh",
	Argv: func() interface{} { t := new(forwardT); t.Self = t; return t },
	Fn:   forwardCLI,

	NumOption: cli.AtLeast(1),
}

////////////////////////////////////////////////////////////////////////////
// send

//  func sendCLI(ctx *cli.Context) error {
//  	rootArgv = ctx.RootArgv().(*rootT)
//  	argv := ctx.Argv().(*sendT)
//  	fmt.Printf("[send]:\n  %+v\n  %+v\n  %v\n", rootArgv, argv, ctx.Args())
//  	Opts.LogLevel, Opts.Verbose, Opts.Version, Opts.Verbose =
//  		rootArgv.LogLevel, rootArgv.Verbose, rootArgv.Version, rootArgv.Verbose.Value()
//  	return nil
//  	//return DoSend()
//  }
//
//  func DoSend() error {
//  	return nil
//  }

type sendT struct {
	Token  string   `cli:"t,token" usage:"The telegram bot token (mandatory)" dft:"$SHUTTLEBOT_TOKEN"`
	ChatID []string `cli:"i,id" usage:"The telegram ChatID(s) (without -) to forward to (mandatory)" dft:"$SHUTTLEBOT_CID"`
	File   string   `cli:"*f,file" usage:"The file spec to send (mandatory)"`
}

var sendDef = &cli.Command{
	Name: "send",
	Desc: "Send file to to the designated ChatID(s)",
	Text: "Usage:\n  shuttlebot send --token $TOKEN --id $GROUP -i $CHANNEL --file /path/to/file",
	Argv: func() interface{} { return new(sendT) },
	Fn:   sendCLI,

	NumOption: cli.AtLeast(1),
}
