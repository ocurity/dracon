package migrations

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/ocurity/dracon/pkg/db"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func newInspectSubCmd() *cobra.Command {
	inspectSubCmd := &cobra.Command{
		Use:     "inspect",
		Short:   "Check which migrations have been applied to the database.",
		GroupID: "Migrations",
		RunE:    entrypointWrapper(inspectMigrations),
	}

	return inspectSubCmd
}

func inspectMigrations(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("you need to provide a path to the migrations that will be applied")
	}
	dirFS := os.DirFS(args[0])

	dbURL, err := db.ParseConnectionStr(migrationsCmdConfig.connStr)
	if err != nil {
		return fmt.Errorf("could not parse connection string: %w", err)
	}

	dbConn, err := dbURL.Connect()
	if err != nil {
		return fmt.Errorf("could not connect to the database: %w", err)
	}

	migrations := db.Migrations{DB: dbConn, PGUrl: dbURL, MigrationsTable: migrationsCmdConfig.migrationsTable}
	latestMigration, isDirty, err := migrations.State(dirFS)
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("could not get state of database: %w", err)
	}

	table := tablewriter.NewWriter(os.Stdout) //cmd.OutOrStdout())
	table.SetHeader([]string{"Foo", "Bar"})
	table.Append([]string{"Latest Migration Version", fmt.Sprintf("%d", latestMigration)})
	table.Append([]string{"Has Failed Migrations", fmt.Sprintf("%v", isDirty)})
	table.Render()

	return nil
}
