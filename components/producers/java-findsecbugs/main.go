package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	v1 "github.com/ocurity/dracon/api/proto/v1"

	"github.com/ocurity/dracon/components/producers"
	"github.com/ocurity/dracon/components/producers/java-findsecbugs/types"
	"github.com/ocurity/dracon/pkg/sarif"
)

func loadXML(filename string) ([]byte, error) {
	xmlFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()
	return io.ReadAll(xmlFile)
}

func readXML(xmlFile []byte) []*v1.Issue {
	/**
	Reads a file containing spotbugs XML results
	and converts the results in the "SECURITY" category
	into an array Dracon issues
	*/

	output := []*v1.Issue{}
	var bugs types.BugCollection
	if len(xmlFile) == 0 {
		return output
	}
	err := xml.Unmarshal(xmlFile, &bugs)
	if err != nil {
		log.Fatal("could not unmarshal findsecbugs output", err)
	}
	for _, instance := range bugs.BugInstance {

		// parse standalone SourceLine elements
		for _, line := range instance.SourceLine {
			output = append(output, parseLine(instance, line))
		}
		// parse SourceLines in Field elements
		for _, field := range instance.Field {
			for _, line := range field.SourceLine {
				output = append(output, parseLine(instance, line))
			}
		}
		// parse SourceLines in Method elements
		for _, method := range instance.Method {
			for _, line := range method.SourceLine {
				output = append(output, parseLine(instance, line))
			}
		}
		// parse SourceLines in Class elements
		for _, cls := range instance.Class {
			for _, line := range cls.SourceLine {
				output = append(output, parseLine(instance, line))
			}
		}

	}
	return output
}

func parseLine(instance types.BugInstance, sourceLine types.SourceLine) *v1.Issue {
	return &v1.Issue{
		Target:      producers.GetFileTarget(sourceLine.Sourcepath, sourceLine.Start, sourceLine.End),
		Type:        instance.Type,
		Severity:    normalizeRank(instance.Rank),
		Cvss:        0.0,
		Confidence:  v1.Confidence(v1.Confidence_value[fmt.Sprintf("CONFIDENCE_%s", "MEDIUM")]),
		Description: instance.LongMessage,
		Title:       instance.ShortMessage,
	}
}

func normalizeRank(rank string) v1.Severity {
	/*
			Normalizes the rank according to the following table
			Scariest: ranked between 1 & 4.
		Scary: ranked between 5 & 9.
		Troubling: ranked between 10 & 14.
		Of concern: ranked between 15 & 20.
	*/
	intRank, err := strconv.ParseInt(rank, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	switch {
	case 1 < intRank && intRank < 4:
		return v1.Severity_SEVERITY_CRITICAL
	case 5 < intRank && intRank < 9:
		return v1.Severity_SEVERITY_HIGH
	case 10 < intRank && intRank < 14:
		return v1.Severity_SEVERITY_MEDIUM
	case 15 < intRank && intRank < 20:
		return v1.Severity_SEVERITY_LOW
	}
	return v1.Severity_SEVERITY_INFO
}

// Sarif is the switch that tells us that findsecbugs output is in sarif format.
var Sarif bool

func main() {
	flag.BoolVar(&Sarif, "sarifOut", false, "Output is in sarif format}")

	if err := producers.ParseFlags(); err != nil {
		log.Fatal(err)
	}

	issues := []*v1.Issue{}
	if Sarif {
		var sarifResults []*sarif.DraconIssueCollection
		inFile, err := producers.ReadInFile()
		if err != nil {
			log.Fatal(err)
		}
		sarifResults, err = sarif.ToDracon(string(inFile))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(len(sarifResults))
		for _, result := range sarifResults {
			if result.ToolName != "SpotBugs" {
				log.Printf("Toolname from Sarif results is not 'SpotBugs' it is %s instead\n", result.ToolName)
			}
			issues = append(issues, result.Issues...)
		}
	} else {
		xmlByteVal, _ := loadXML(producers.InResults)
		issues = readXML(xmlByteVal)
	}
	if err := producers.WriteDraconOut(
		"spotbugs",
		issues,
	); err != nil {
		log.Fatal(err)
	}
}
