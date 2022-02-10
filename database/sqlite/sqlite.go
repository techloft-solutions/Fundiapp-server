package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"sort"
	"time"

	_ "github.com/go-sql-driver/mysql"
	//_ "github.com/mattn/go-sqlite3"
)

const (
	statusPending  = "pending"
	statusCanceled = "canceled"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

//go:embed seeds/*.sql
var seedFS embed.FS

type Criteria struct {
	Needle   string
	Haystack string
}

type DB struct {
	db     *sql.DB
	ctx    context.Context // background context
	cancel func()          // cancel background context
	Now    func() time.Time
	// Datasource name.
	DSN string
}

// NewDB returns a new instance of DB associated with the given datasource name.
func NewDB(dsn string) *DB {
	db := &DB{
		DSN: dsn,
		Now: time.Now,
	}
	db.ctx, db.cancel = context.WithCancel(context.Background())
	return db
}

// Open opens the database connection.
func (db *DB) Open() (err error) {
	// Ensure a DSN is set before attempting to open the database.
	if db.DSN == "" {
		return fmt.Errorf("dsn required")
	}

	// Make the parent directory unless using an in-memory db.
	/*if db.DSN != ":memory:" {
		if err := os.MkdirAll(filepath.Dir(db.DSN), 0700); err != nil {
			return err
		}
	}*/

	// Connect to the database.
	if db.db, err = sql.Open("mysql", db.DSN); err != nil {
		return err
	}

	db.db.SetConnMaxLifetime(time.Minute * 5)
	db.db.SetMaxOpenConns(25)
	db.db.SetMaxIdleConns(25)

	// Enable WAL. SQLite performs better with the WAL  because it allows
	// multiple readers to operate while data is being written.
	//if _, err := db.db.Exec(`PRAGMA journal_mode = wal;`); err != nil {
	//	return fmt.Errorf("enable wal: %w", err)
	//}

	// Enable foreign key checks. For historical reasons, SQLite does not check
	// foreign key constraints by default... which is kinda insane. There's some
	// overhead on inserts to verify foreign key integrity but it's definitely
	// worth it.
	//if _, err := db.db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
	//	return fmt.Errorf("foreign keys pragma: %w", err)
	//}
	/*
		if err := db.migrate(); err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
		time.Sleep(1 * time.Second)

		if err := db.seed(); err != nil {
			log.Println("seeding error:", err)
		}
	*/
	return nil
}

// Tx wraps the SQL Tx object to provide a timestamp at the start of the transaction.
type Tx struct {
	*sql.Tx
	db  *DB
	now time.Time
}

func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Return wrapper Tx that includes the transaction start time.
	return &Tx{
		Tx:  tx,
		db:  db,
		now: db.Now().UTC().Truncate(time.Second),
	}, nil
}

// drop drops the database tables and resets the migrations table. This is used to drop the database and start over.
//
// This is a destructive operation and should only be used for testing.
func (db *DB) Drop() error {
	log.Println("dropping database tables")
	// Read drop files from our embedded file system.
	names, err := fs.Glob(migrationFS, "migrations/*_down.sql")
	if err != nil {
		return err
	}
	sort.Strings(names)

	// disable foreign key checks
	if _, err := db.db.Exec(`SET foreign_key_checks = 0;`); err == nil {
		log.Println("disable foreign key checks:")
	}

	// Loop over all migration files and execute them in order.
	for _, name := range names {
		if err := db.dropFile(name); err != nil {
			return fmt.Errorf("dropping error: name=%q err=%w", name, err)
		}
	}

	// enable foreign key checks
	if _, err := db.db.Exec(`SET foreign_key_checks = 1;`); err == nil {
		log.Println("enable foreign key checks:")
	}
	log.Println("tables dropped!")

	return nil
}

func (db *DB) dropFile(name string) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Read and execute drop file.
	if buf, err := fs.ReadFile(migrationFS, name); err != nil {
		return err
	} else if _, err := tx.Exec(string(buf)); err != nil {
		return err
	}
	return tx.Commit()
}

func (db *DB) dropQuery() error {
	query := `
	SELECT CONCAT('DROP TABLE IF EXISTS ', table_name, ';')
	FROM information_schema.tables
	WHERE table_schema = 'hudumaapp'
	ORDER BY table_name;
	`

	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Read and execute drop file.
	if _, err := tx.Exec(query); err != nil {
		return err
	}
	return tx.Commit()
}

// migrate sets up migration tracking and executes pending migration files.
//
// Migration files are embedded in the sqlite/migration folder and are executed
// in lexigraphical order.
//
// Once a migration is run, its name is stored in the 'migrations' table so it
// is not re-executed. Migrations run in a transaction to prevent partial
// migrations.
func (db *DB) Migrate() error {
	log.Println("migrating database...")
	/*
		if os.Getenv("APP_ENV") != "production" || os.Getenv("DB_RESET") == "true" {
			if err := db.drop(); err != nil {
				log.Println(err)
			}
		}

		time.Sleep(15 * time.Second)
	*/
	// Ensure the 'migrations' table exists so we don't duplicate migrations.
	if _, err := db.db.Exec(`CREATE TABLE IF NOT EXISTS migrations (name varchar(255) PRIMARY KEY);`); err != nil {
		return fmt.Errorf("cannot create migrations table: %w", err)
	}

	// Read migration files from our embedded file system.
	// This uses Go 1.16's 'embed' package.
	names, err := fs.Glob(migrationFS, "migrations/*_up.sql")
	if err != nil {
		return err
	}
	sort.Strings(names)

	// Loop over all migration files and execute them in order.
	for _, name := range names {
		if err := db.migrateFile(name); err != nil {
			return fmt.Errorf("migration error: name=%q err=%w", name, err)
		}
	}
	log.Println("migrations DONE!")
	return nil
}

// migrate runs a single migration file within a transaction. On success, the
// migration file name is saved to the "migrations" table to prevent re-running.
func (db *DB) migrateFile(name string) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Ensure migration has not already been run.
	var n int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM migrations WHERE name = ?`, name).Scan(&n); err != nil {
		return err
	} else if n != 0 {
		return nil // already run migration, skip
	}

	// Read and execute migration file.
	if buf, err := fs.ReadFile(migrationFS, name); err != nil {
		return err
	} else if _, err := tx.Exec(string(buf)); err != nil {
		return err
	}

	// Insert record into migrations to prevent re-running migration.
	if _, err := tx.Exec(`INSERT INTO migrations (name) VALUES (?)`, name); err != nil {
		return err
	}

	return tx.Commit()
}

func (db *DB) Seed() error {
	log.Println("seeding database...")

	names, err := fs.Glob(seedFS, "seeds/*_seed.sql")
	if err != nil {
		return err
	}
	sort.Strings(names)

	// Loop over all seed files and execute them in order.
	for _, name := range names {
		if err := db.seedFile(name); err != nil {
			return fmt.Errorf("seeding error: name=%q err=%w", name, err)
		}
	}
	log.Println("seeding DONE!")
	return nil
}

func (db *DB) seedFile(name string) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Read and execute seed file.
	if buf, err := fs.ReadFile(seedFS, name); err != nil {
		return err
	} else if _, err := tx.Exec(string(buf)); err != nil {
		return err
	}

	return tx.Commit()
}
