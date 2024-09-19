package main

import (
	"log"
	"log/slog"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components/producers"
	"github.com/ocurity/dracon/pkg/sarif"
)

func main() {

	if err := producers.ParseFlags(); err != nil {
		log.Fatal(err)
	}
	producers.Append = true

	inFile, err := producers.ReadInFile()
	if err != nil {
		log.Fatal(err)
	}
	results, err := processInput(string(inFile))
	if err != nil {
		log.Fatal(err)
	}
	if err := writeOutput(results); err != nil {
		log.Fatal(err)
	}
}

func writeOutput(results map[string][]*v1.Issue) error {
	for _, issues := range results {
		slog.Info(
			"appending",
			slog.Int("issues", len(issues)),
			slog.String("tool", "snuk"),
		)
		if err := producers.WriteDraconOut(
			"snyk",
			issues,
		); err != nil {
			slog.Error("error writing dracon out for the snyk tool", "err", err)
		}
	}
	return nil
}

func processInput(input string) (map[string][]*v1.Issue, error) {
	issues, err := sarif.ToDracon(string(input))
	if err != nil {
		return nil, err
	}
	results := map[string][]*v1.Issue{}
	for _, output := range issues {
		results[output.ToolName] = append(results[output.ToolName], output.Issues...)
	}
	return results, nil
}
