// package main of the multitool deduplication enricher finds and tags issues found by multiple tools in this scan

package main

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"time"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/package-url/packageurl-go"

	"github.com/ocurity/dracon/pkg/putil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// for every tool response
//  for every issue
// 		if issue target is purl
// 			func handle purls(sca findings)
// 			if the issue has a CVE attached
// 				if cveMap[cve] != nil
// 					mark dup, attach which other tool found it
// 				else
//	 				map the cve to the issue
// 			else if the issue does not have a CVE attached
//  TODO match if the tool is SCA, and match output title and remote references to some CRE, then try to match to CREs of other tools
// 		if issue target is code (file/system/path:line or file/system/path:line-range)
// 			func handle sast
// 				if the issue has a CWE
// 					if cweMap[cwe] != nil
// 						if cweMap[cwe] target and line (range)  includes the line and targer of current issue
// 							mark dup, attach which other tool found it
// 						else
// 							append to cweMap
// 					else
// 						cweMap[cwe] = []v1.Issue{issue}
//				else if the issue does not have a CWE
// TODO match if the tool is SAST, match title and description (and content of remote references) to CRE, try to match to CREs of other tools

const AnnotationStr = "Found By Other Tools"

var (
	connStr   string
	readPath  string
	writePath string
	cve       map[string][]tooledIssue = map[string][]tooledIssue{}
	// cwe                 map[int32][]tooledIssue  = map[string][]tooledIssue{}
	// cre                 map[string][]tooledIssue = map[string][]tooledIssue{}
	filesystemTargetReg = `(?P<path>^(.+)\/([^\/]+)/?)(?P<range>(?P<lineFrom>:\d+)(?P<lineTo> ?- ? \d+)?)?$`
)

type tooledIssue struct {
	issue         *v1.Issue
	producingTool string
}

// type lineRange struct {
// 	filePath  string
// 	lineStart int
// 	lineEnd   int
// }

// func convertTargetToLineRange(target string) lineRange {
// 	result := lineRange{}
// 	compiledRegexp := regexp.MustCompile(filesystemTargetReg)
// 	matches := compiledRegexp.FindStringSubmatch(target)
// 	result.filePath = matches[compiledRegexp.SubexpIndex("path")]
// 	lineStart := strings.TrimSpace(strings.ReplaceAll(matches[compiledRegexp.SubexpIndex("lineFrom")], ":", ""))
// 	if lineStart != "" {
// 		ls, err := strconv.Atoi(lineStart)
// 		if err != nil {
// 			log.Println("Could not recognise line start for issue target", target)
// 		}
// 		result.lineStart = ls
// 	} else {
// 		fmt.Printf("%#v\n", matches)
// 	}
// 	lineEnd := strings.TrimSpace(strings.ReplaceAll(matches[compiledRegexp.SubexpIndex("lineTo")], "-", ""))
// 	if lineEnd != "" {
// 		le, err := strconv.Atoi(lineEnd)
// 		if err != nil {
// 			log.Println("Could not recognise line end for issue target", target)
// 		}
// 		result.lineEnd = le
// 	} else {
// 		fmt.Printf("%#v\n", matches)
// 	}
// 	return result
// }

// func handleFilesystem(issue *v1.Issue, issueTool string) string {
// 	target := convertTargetToLineRange(issue.Target)
// 	if len(issue.Cwe) > 0 {
// 		for _, cweNum := range issue.Cwe {
// 			cweIssues, ok := cwe[cweNum]
// 			if ok {
// 				for _, cweIssue := range cweIssues {
// 					matchedTarget := convertTargetToLineRange(cweIssue.issue.Target)
// 					if matchedTarget.filePath == target.filePath { // if we are on the same file
// 						if (matchedTarget.lineStart <= target.lineStart && // and line range for existing vuln is larger than line range for new vuln
// 							matchedTarget.lineEnd >= target.lineEnd) ||
// 							(matchedTarget.lineStart >= target.lineStart && // or if the line range for the new vuln contains the line range for the existing vuln
// 								matchedTarget.lineEnd <= target.lineEnd) {
// 							return cweIssue.producingTool
// 						}
// 					}
// 				}
// 			} else {
// 				cwe[cweNum] = []tooledIssue{{issue: issue, producingTool: issueTool}}
// 			}
// 		}
// 	}
// 	// @TODO(spyros) match if the tool is SAST, match title and description (and content of remote references) to CRE, try to match to CREs of other tools
// 	return ""
// }

func handlePurl(issue *v1.Issue, issueTool string) string {
	if issue.Cve != "" {
		cveIssues, ok := cve[issue.Cve]
		if ok {
			return cveIssues[0].producingTool
		}
		cve[issue.Cve] = []tooledIssue{{
			issue:         issue,
			producingTool: issueTool,
		}}
	}
	//@TODO(spyros): match if the tool is SCA, and match output title and remote references to some CRE, then try to match to CREs of other tools
	return ""
}

func enrichIssue(issue *v1.Issue, foundTool string) v1.EnrichedIssue {
	matchFilesystem := regexp.MustCompile(filesystemTargetReg)
	_, err := packageurl.FromString(issue.Target)
	if err == nil {
		if tool := handlePurl(issue, foundTool); tool != "" {
			if foundTool != tool {
				return v1.EnrichedIssue{
					RawIssue:    issue,
					Annotations: map[string]string{AnnotationStr: "true"},
				}
			} else {
				log.Println("Duplicate CVE found for tool", tool, "and CVE", issue.Cve)
				return v1.EnrichedIssue{
					RawIssue: issue,
				}
			}
		} else {
			fmt.Println("registered new purl", issue.Target)
			return v1.EnrichedIssue{
				RawIssue: issue,
			}
		}
	} else if matchFilesystem.Match([]byte(issue.Target)) {
		return v1.EnrichedIssue{
			RawIssue: issue,
		}
		// SAST not supported yet!!!
		// 	if tool := handleFilesystem(issue, foundTool); tool != "" {
		// 		if foundTool != tool {
		// 			return v1.EnrichedIssue{
		// 				RawIssue:    issue,
		// 				Annotations: map[string]string{AnnotationStr: "true"},
		// 			}
		// 		} else {
		// 			log.Println("Duplicate CWE found for tool", tool, "and CWEs", issue.Cwe)
		// 			return v1.EnrichedIssue{
		// 				RawIssue: issue,
		// 			}
		// 		}
		// 	}
	}
	log.Println("Target", issue.Target, "is not supported, at this time this enricher only supports filesystem targets or purls")
	return v1.EnrichedIssue{
		RawIssue: issue,
	}
}

var rootCmd = &cobra.Command{
	Use:   "enricher",
	Short: "enricher",
	Long:  "tool to enrich issues against a database",
	RunE: func(cmd *cobra.Command, args []string) error {
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
				eI := enrichIssue(i, r.ToolName)
				enrichedIssues = append(enrichedIssues, &eI)
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
					r.GetScanInfo().GetScanTags(),
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
