// Package main of the cdxgen producer parses the CycloneDX output of cdxgen and
// create a singular Dracon issue from it
package main

import (
	"log"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components/producers"
	"github.com/ocurity/dracon/pkg/cyclonedx"
)

func main() {
	if err := producers.ParseFlags(); err != nil {
		log.Fatal(err)
	}
	var results []*v1.Issue
	inFile, err := producers.ReadInFile()
	if err != nil {
		log.Fatal("could not load file err:%s", err)
	}
	results, err = handleCycloneDX(inFile)
	if err != nil {
		log.Fatalf("could not parse cyclonedx document err:%s", err)
	}
	if err := producers.WriteDraconOut(
		"cdxgen", results,
	); err != nil {
		log.Fatal("could not write dracon out err:%s", err)
	}
}

func handleCycloneDX(inFile []byte) ([]*v1.Issue, error) {
	return cyclonedx.ToDracon(inFile, "json")
}
