// Package main of the npm-audit component transforms the output of either npm full-audit or npm quick-audit to dracon
package main

import (
	"errors"
	"flag"
	"log"

	"github.com/ocurity/dracon/components/producers"
	atypes "github.com/ocurity/dracon/components/producers/typescript-npm-audit/types"
	"github.com/ocurity/dracon/components/producers/typescript-npm-audit/types/npmfullaudit"
	"github.com/ocurity/dracon/components/producers/typescript-npm-audit/types/npmquickaudit"
)

// producer flags.
var (
	PackagePath string
)

func inFileToReport(inFile []byte) (atypes.Report, error) {
	reportConstructors := []func([]byte) (atypes.Report, error){
		npmquickaudit.NewReport,
		npmfullaudit.NewReport,
	}

	for _, constructor := range reportConstructors {
		report, err := constructor(inFile)

		//nolint:errorlint
		switch err.(type) {
		case nil:
			return report, nil
		case *atypes.ParsingError, *atypes.FormatError:
			// Ignore parsing and incorrect format errors from constructors -
			// we'll just attempt again with the next one
		default:
			// Any other errors returned by a constructor are likely fatal
			return nil, err
		}
	}

	return nil, errors.New("input file is not a supported npm audit report format")
}

func main() {
	flag.StringVar(&PackagePath, "package-path", "", "Path to the package.json file corresponding to this audit report; will be prepended to vulnerable dependency names in issue reports if specified")

	if err := producers.ParseFlags(); err != nil {
		log.Fatal(err)
	}

	inFile, err := producers.ReadInFile()
	if err != nil {
		log.Fatal(err)
	}

	report, err := inFileToReport(inFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Parsed input file as %s\n", report.Type())

	report.SetPackagePath(PackagePath)

	if err := producers.WriteDraconOut(
		"npm-audit",
		report.AsIssues(),
	); err != nil {
		log.Fatal(err)
	}
}
