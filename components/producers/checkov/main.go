package main

import (
	"flag"
	"log"

	"github.com/go-errors/errors"

	draconv1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components/producers"
	"github.com/ocurity/dracon/pkg/cyclonedx"
	"github.com/ocurity/dracon/pkg/sarif"
)

// the CycloneDX target override
var target string

func main() {
	flag.StringVar(&target, "target", "", "The target being scanned, this will override the CycloneDX target and is useful for cases where you scan iac or a dockerfile for an application that you know it's purl")

	if err := producers.ParseFlags(); err != nil {
		log.Fatal(err)
	}

	inFile, err := producers.ReadInFile()
	if err != nil {
		log.Fatal(err)
	}
	if err := run(inFile, target); err != nil {
		log.Fatal(err)
	}
}

func run(inFile []byte, target string) error {
	sarifResults, sarifErr := handleSarif(inFile)
	cyclondxResults, cyclonedxErr := handleCycloneDX(inFile, target)
	var issues []*draconv1.Issue
	if sarifErr == nil {
		issues = sarifResults
	} else if cyclonedxErr == nil {
		issues = cyclondxResults
	} else {
		return errors.Errorf("Could not parse input file as neither Sarif nor CycloneDX sarif error: %v, cyclonedx error: %v", sarifErr, cyclonedxErr)
	}
	return producers.WriteDraconOut(
		"checkov",
		issues,
	)
}

func handleSarif(inFile []byte) ([]*draconv1.Issue, error) {
	var sarifResults []*sarif.DraconIssueCollection
	var draconResults []*draconv1.Issue
	sarifResults, err := sarif.ToDracon(string(inFile))
	if err != nil {
		return draconResults, err
	}
	for _, result := range sarifResults {
		draconResults = append(draconResults, result.Issues...)
	}
	return draconResults, nil
}

func handleCycloneDX(inFile []byte, target string) ([]*draconv1.Issue, error) {
	return cyclonedx.ToDracon(inFile, "json", target)
}
