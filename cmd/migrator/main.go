package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	//var storagePath, migrationsPath, migrationsTable string
	var (
	migrationsPath  string
	migrationsTable string
	host            string
	port            string
	user            string
	password        string
	dbname          string
)

	// flag.StringVar(&storagePath, "storage-path", "", "path to storage")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")

	flag.StringVar(&host, "host", "localhost", "")
	flag.StringVar(&port, "port", "5433", "")
	flag.StringVar(&user, "user", "postgres_auth", "")
	flag.StringVar(&password, "password", "1234", "")
	flag.StringVar(&dbname, "dbname", "sso", "")

	flag.Parse()
	// if storagePath == "" {
	// 	panic("storage-path is required")
	// }
	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable&x-migrations-table=%s",
			user, password, host, port, dbname, migrationsTable,
		),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}

		panic(err)
	}

	fmt.Println("migrations apllied successfully")
}
