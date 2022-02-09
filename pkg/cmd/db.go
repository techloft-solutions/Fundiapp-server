package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/andrwkng/hudumaapp/database/sqlite"
)

type DBCommand struct {
	DB *sqlite.DB
}

func NewDBCommand(db *sqlite.DB) *DBCommand {
	return &DBCommand{
		DB: db,
	}
}

func (d *DBCommand) Run(ctx context.Context, args []string) error {
	var cmd string
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		cmd, args = args[0], args[1:]
	}

	switch cmd {
	case "migrate":
		return d.migrate()
	case "seed":
		return d.seed()
	default:
		return fmt.Errorf("ServiceApp cli %s: unknown command", cmd)
	}
}

func (d *DBCommand) migrate() error {
	log.Println("migrate")
	return d.DB.Migrate()
}

func (d *DBCommand) seed() error {
	log.Println("seed")
	return d.DB.Seed()
}
