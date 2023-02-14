package main

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	v1 "github.com/ocurity/dracon/api/proto/v1"

	"github.com/ocurity/dracon/pkg/enrichment"
	"github.com/ocurity/dracon/pkg/enrichment/db"
	"github.com/ocurity/dracon/pkg/putil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	connStr   string
	readPath  string
	writePath string
)

var rootCmd = &cobra.Command{
	Use:   "enricher",
	Short: "enricher",
	Long:  "tool to enrich issues against a database",
	RunE: func(cmd *cobra.Command, args []string) error {
		connStr = viper.GetString("db_connection")
		db, err := db.NewDB(connStr)
		if err != nil {
			return err
		}
		readPath = viper.GetString("read_path")
		res, err := putil.LoadTaggedToolResponse(readPath)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Loaded %d tagged tool responses\n", len(res))
		writePath = viper.GetString("write_path")
		for _, r := range res {
			enrichedIssues := []*v1.EnrichedIssue{}
			for _, i := range r.GetIssues() {
				eI, err := enrichment.EnrichIssue(db, i)
				if err != nil {
					log.Printf("error enriching issue %s, err: %#v\n", i.Uuid, err)
					continue
				}
				enrichedIssues = append(enrichedIssues, eI)
				log.Printf("enriched issue '%s'", eI.GetRawIssue().GetUuid())
			}
			if len(enrichedIssues) > 0 {
				if err := putil.WriteEnrichedResults(r, enrichedIssues,
					filepath.Join(writePath, fmt.Sprintf("%s.enriched.pb", r.GetToolName())),
				); err != nil {
					return err
				}
			}
			if len(r.GetIssues()) > 0 {
				scanStartTime := r.GetScanInfo().GetScanStartTime().AsTime()
				if err := putil.WriteResults(
					r.GetToolName(),
					r.GetIssues(),
					filepath.Join(writePath, fmt.Sprintf("%s.raw.pb", r.GetToolName())),
					r.GetScanInfo().GetScanUuid(),
					scanStartTime.Format(time.RFC3339),
				); err != nil {
					log.Fatalf("could not write results: %s", err)
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.Flags().StringVar(&connStr, "db_connection", "", "the database connection DSN")
	rootCmd.Flags().StringVar(&readPath, "read_path", "", "the path to read LaunchToolResponses from")
	rootCmd.Flags().StringVar(&writePath, "write_path", "", "the path to write enriched results to")
	if err := viper.BindPFlag("db_connection", rootCmd.Flags().Lookup("db_connection")); err != nil {
		log.Fatalf("could not bind db_connection flag: %s", err)
	}
	if err := viper.BindPFlag("read_path", rootCmd.Flags().Lookup("read_path")); err != nil {
		log.Fatalf("could not bind read_path flag: %s", err)
	}
	if err := viper.BindPFlag("write_path", rootCmd.Flags().Lookup("write_path")); err != nil {
		log.Fatalf("could not bind write_path flag: %s", err)
	}
	viper.SetEnvPrefix("enricher")
	viper.AutomaticEnv()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
