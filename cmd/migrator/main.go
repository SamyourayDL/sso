package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"

	//driver to make migrations SQLite 3
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	//driver to get migrations from file
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationsPath, migrationsTable string
	flag.StringVar(&storagePath, "storage-path", "", "path to storage file")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations files")
	flag.StringVar(&storagePath, "migrations-table", "migrations", "migrations table name")

	flag.Parse()

	if storagePath == "" {
		panic("storage-path is required")
	}
	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
	}
}
