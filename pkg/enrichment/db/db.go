package db

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/url"
	"path/filepath"
	"strings"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/lib/pq"
	v1 "github.com/ocurity/dracon/api/proto/v1"

	"github.com/jmoiron/sqlx"
)

// migrationsFS holds the SQL migration files as static assets.
//
//nolint:typecheck
//go:embed migrations/*.sql
var migrationsFS embed.FS

// EnrichDatabase represents the db methods that are used for the enricher.
type EnrichDatabase interface {
	GetIssueByHash(string) (*v1.EnrichedIssue, error)
	CreateIssue(context.Context, *v1.EnrichedIssue) error
	UpdateIssue(context.Context, *v1.EnrichedIssue) error
	DeleteIssueByHash(string) error
}

// DB  Database implements DB.
type DB struct {
	*sqlx.DB
}

// NewDB returns a new DB for the enricher.
func NewDB(connStr string) (*DB, error) {
	searchPath, err := getSchemaSearchPathFromConnStr(connStr)
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	migrationsConfig := &postgres.Config{}
	if searchPath != "" {
		_, err := db.Exec(fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s`,
			pq.QuoteIdentifier(searchPath)))
		if err != nil {
			return nil, err
		}

		migrationsConfig.SchemaName = searchPath
	}

	driver, err := postgres.WithInstance(db.DB, migrationsConfig)
	if err != nil {
		return nil, err
	}

	var assetNames []string
	if err := fs.WalkDir(migrationsFS, ".", func(path string, info fs.DirEntry, err error) error {
		if !info.IsDir() {
			assetNames = append(assetNames, info.Name())
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("could not walk migrations: %w", err)
	}

	resources := bindata.Resource(assetNames,
		func(name string) ([]byte, error) {
			return fs.ReadFile(migrationsFS, filepath.Join(".", name))
		})
	log.Printf("migrations discovered are %q", resources.Names)

	resourcesDriver, err := bindata.WithInstance(resources)
	if err != nil {
		return nil, fmt.Errorf("could not create migration bindata instance: %w", err)
	}

	m, err := migrate.NewWithInstance("go-bindata", resourcesDriver, "dracon", driver)
	if err != nil {
		return nil, err
	}

	dbVersion, isDBDirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return nil, fmt.Errorf("error getting migration db version: %w", err)
	}

	log.Printf("migrationVersion: %d, isDBDirty %v", dbVersion, isDBDirty)
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("could not migrate DB: %w", err)
	}

	return &DB{db}, nil
}

// getSchemaSearchPathFromConnStr extracts the database schema component from a
// PostgreSQL connection string; if no schema was specified, the empty string is
// returned.
func getSchemaSearchPathFromConnStr(connStr string) (string, error) {
	url, err := url.Parse(connStr)

	if err == nil && url.Scheme == "postgres" {
		return getSchemaSearchPathFromURL(url)
	}
	return getSchemaSearchPathFromKV(connStr)
}

// getSchemaSearchPathFromURL extracts the schema search path component from a
// PostgreSQL connection URL; if no search path is specified, the empty string
// is returned.
func getSchemaSearchPathFromURL(connURL *url.URL) (string, error) {
	path, found := connURL.Query()["search_path"]
	if !found {
		return "", nil
	}

	if len(path) == 0 {
		return "", nil
	} else if len(path) == 1 {
		return path[0], nil
	}
	return "", errors.New("Multiple search_paths defined in database connection DSN")
}

// getSchemaSearchPathFromKV extracts the schema search path component from a
// PostgreSQL keyword/value connection string; if no search path is specified,
// the empty string is returned.
func getSchemaSearchPathFromKV(kvStr string) (string, error) {
	var path string

	for _, pair := range strings.Fields(kvStr) {
		elems := strings.SplitN(pair, "=", 2)
		if elems[0] == "search_path" {
			if path != "" {
				return "", errors.New("Multiple search_paths defined in database connection DSN")
			}

			path = elems[1]
		}
	}

	return path, nil
}
