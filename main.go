package main

import (
	"fmt"
	"github.com/coreos/bbolt"
	"os"
	"os/signal"
	"runtime"
	"shopaholic/cmd"
	"shopaholic/store/engine"
	"shopaholic/store/service"
	"syscall"
	"time"

	log "github.com/go-pkgz/lgr"
	"github.com/jessevdk/go-flags"
)

// Opts with all cli commands and flags
type Opts struct {
	cmd.UserCreateCommand        `command:"user:create"`
	cmd.UserListCommand          `command:"user:list"`
	cmd.TransactionCreateCommand `command:"transaction:create"`
	cmd.TransactionListCommand   `command:"transaction:list"`

	Currency   string `long:"currency" env:"CURRENCY" default:"usd" description:"money currency"`
	DBFilename string `long:"dbfilename" env:"DBFILENAME" default:"shopaholic.db" description:"database filename"`

	Dbg bool `long:"dbg" env:"DEBUG" description:"debug mode"`
}

var revision = "unknown"

func main() {
	fmt.Printf("shopaholic %s\n", revision)

	var opts Opts
	p := flags.NewParser(&opts, flags.Default)

	p.CommandHandler = func(command flags.Commander, args []string) error {
		setupLog(opts.Dbg)
		// commands implements CommonOptionsCommander to allow passing set of extra options defined for all commands
		c := command.(cmd.Commander)
		c.SetCommon(cmd.CommonOpts{
			Currency: opts.Currency,
			Store:    *initDataStore(opts.DBFilename),
		})
		err := c.Execute(args)
		if err != nil {
			log.Printf("[ERROR] failed with %+v", err)
		}
		return err
	}

	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}

func initDataStore(dbFilename string) *service.DataStore {
	b, err := engine.NewBoltDB(bbolt.Options{Timeout: 30 * time.Second}, dbFilename)
	if err != nil {
		log.Printf("[ERROR] can not initiate DB %s", dbFilename)
		os.Exit(0)
	}

	return &service.DataStore{Interface: b}
}

func setupLog(dbg bool) {
	if dbg {
		log.Setup(log.Debug, log.CallerFile, log.Msec, log.LevelBraces)
		return
	}
	log.Setup(log.Msec, log.LevelBraces, log.CallerPkg, log.CallerIgnore("logger"))
}

// getDump reads runtime stack and returns as a string
func getDump() string {
	maxSize := 5 * 1024 * 1024
	stacktrace := make([]byte, maxSize)
	length := runtime.Stack(stacktrace, true)
	if length > maxSize {
		length = maxSize
	}
	return string(stacktrace[:length])
}

func init() {
	sigChan := make(chan os.Signal)
	go func() {
		for range sigChan {
			log.Printf("[INFO] SIGQUIT detected, dump:\n%s", getDump())
		}
	}()
	signal.Notify(sigChan, syscall.SIGQUIT)
}
