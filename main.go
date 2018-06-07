//go:generate go-bindata -pkg $GOPACKAGE -prefix ./ -o bindata.go ./messages.tmpl

package main

import (
	"fmt"
	"os"
	"time"

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
