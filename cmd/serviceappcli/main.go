package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	app "github.com/andrwkng/hudumaapp"
	"github.com/andrwkng/hudumaapp/config"
	"github.com/andrwkng/hudumaapp/database/sqlite"
	"github.com/andrwkng/hudumaapp/pkg/cmd"
	"github.com/go-sql-driver/mysql"
)

func main() {

	// Setup signal handlers.
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	// Execute program.
	//
	// If an ErrHelp error is returned then that means the user has used an "-h"
	// flag and the flag package will handle output. We just need exit.
	//
	// If we have an application error (wtf.Error) then we can just display the
	// message. If we have any other error, print the raw error message.
	var e *app.Error
	if err := Run(ctx, os.Args[1:]); err == flag.ErrHelp {
		os.Exit(1)
	} else if errors.As(err, &e) {
		fmt.Fprintln(os.Stderr, e.Message)
		os.Exit(1)
	} else if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

// Run executes the main program.
func Run(ctx context.Context, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config", err)
	}

	log.Println("Config:", cfg)

	dbCfg := mysql.Config{
		User:   cfg.DBUser,
		Net:    "tcp",
		Addr:   cfg.DBAddr,
		DBName: cfg.DBName,
		Passwd: cfg.DBPass,
		Params: nil,
	}

	db := sqlite.NewDB(dbCfg.FormatDSN())

	err = db.Open()
	if err != nil {
		log.Fatal(err)
	}

	// Shift off subcommand from the argument list, if available.
	var cmdName string
	if len(args) > 0 {
		cmdName, args = args[0], args[1:]
	}

	// Delegate subcommands to their own Run() methods.
	switch cmdName {
	case "db":
		//cli := &cli.DBCommand{}
		//cli.DB = db
		//return (&DBCommand{}).Run(ctx, args)
		cli := cmd.NewDBCommand(db)
		return cli.Run(ctx, args)
	default:
		return fmt.Errorf("serviceAapp %s: unknown command", cmdName)
	}
}
