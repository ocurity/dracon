package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/jmespath/go-jmespath"
	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/pkg/putil"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const defaultAnnotation = "Reachability"

var (
	readPath   string
	writePath  string
	annotation string
	sliceFile  string
)

type regexes struct {
	purlPkg             *regexp.Regexp
	purlTrailingVersion *regexp.Regexp
	purlVersion         *regexp.Regexp
	filename            *regexp.Regexp
	fileLine            *regexp.Regexp
}

// checkReachable checks if the search term is reachable
func checkReachable(data ReachableSlice, search string, patterns regexes) (bool, error) {
	if search == "" {
		return false, nil
	}
	if result := checkReachablePurl(data, search, patterns); result == true {
		return true, nil
	}

	match := patterns.fileLine.FindStringSubmatch(search)
	if match != nil {
		file := match[1]
		lineMin, lineMax := getLnRange(match[2])
		return filterFlows(data, file, lineMin, lineMax), nil
	}
	return false, errors.New(fmt.Sprintf("Invalid search term: %s", search))
}

// checkReachablePurl checks if a given pkg:version is included in the reachables
func checkReachablePurl(data ReachableSlice, purl string, patterns regexes) bool {
	purls := enumerateReachablePurls(data, patterns)
	purl = strings.ToLower(purl)
	_, exists := purls[purl]
	if exists {
		return true
	}
	return false
}

func enrichIssue(i *v1.Issue, data ReachableSlice, patterns regexes) (*v1.EnrichedIssue, error) {
	enrichedIssue := v1.EnrichedIssue{}
	//annotations := map[string]string{}
	//targetType, newTarget := identifyTargetType(i.Target)
	//var result bool
	//var err error
	result, err := checkReachable(data, i.Target, patterns)
	if err != nil {
	}
	if result == true {

	}
	enrichedIssue = v1.EnrichedIssue{
		RawIssue:    i,
		Annotations: map[string]string{},
	}
	enrichedIssue.Annotations["reachable"] = strconv.FormatBool(result)
	return &enrichedIssue, err
}

// enumerateReachablePurls extracts all the reached purls from the reachables
func enumerateReachablePurls(data ReachableSlice, patterns regexes) map[string]string {
	rawPurls, _ := jmespath.Search("reachables[].purls[]", data)
	allPurls := make(map[string]string)
	for _, purl := range rawPurls.([]interface{}) {
		parsedPurls := parsePurl(purl.(string), patterns)
		for _, pp := range parsedPurls {
			allPurls[pp] = ""
		}
	}
	return allPurls
}

// filterFlows filters flows based on reachables, filename, and line numbers.
func filterFlows(data ReachableSlice, filename string, lnMin int, lnMax int) bool {
	for _, flows := range data.Reachables {
		for _, f := range flows.Flows {
			if f.LineNumber != 0 && !contains(lnMin, lnMax, f.LineNumber) {
				continue
			}
			if strings.HasSuffix(f.ParentFileName, filename) {
				return true
			}
		}
	}
	return false
}

// getLnRange extracts line numbers from a string and returns a tuple of (start, end).
func getLnRange(value string) (int, int) {
	if strings.Contains(value, "-") {
		values := strings.Split(value, "-")
		if len(values) == 2 {
			start, err1 := strconv.Atoi(values[0])
			end, err2 := strconv.Atoi(values[1])
			if err1 == nil && err2 == nil {
				return start, end
			}
		}
	} else {
		num, err := strconv.Atoi(value)
		if err == nil {
			return num, num
		}
	}
	log.Printf("Ignoring invalid line number: %s.", value)
	return 0, 0
}

// contains checks if an integer is in a slice.
func contains(ln int, lnMin int, lnMax int) bool {
	return (ln-lnMax)*(ln-lnMin) <= 0
}

// parsePurlPkgs extracts package and version variations from a purl.
func parsePurlPkgs(matches []string, pattern *regexp.Regexp) []string {
	// Creating a map to ensure uniqueness
	pkgSet := make(map[string]struct{})

	// Adding the packages
	pkgSet[matches[pattern.SubexpIndex("p1")]] = struct{}{}
	pkgSet[matches[pattern.SubexpIndex("p2")]] = struct{}{}

	// Converting the map to a slice and cleaning up the packages
	var pkgs []string
	for pkg := range pkgSet {
		pkg = strings.ReplaceAll(pkg, "pypi/", "")
		pkg = strings.ReplaceAll(pkg, "npm/", "")
		pkg = strings.ReplaceAll(pkg, "%40", "@")
		pkgs = append(pkgs, pkg)
	}

	return pkgs
}

// parsePurlVersions returns a list of version variations from a purl.
func parsePurlVersions(matches []string, pattern *regexp.Regexp) []string {
	var versions []string
	if len(matches) == 0 {
		return versions
	}

	// Creating a map to ensure uniqueness
	versionSet := make(map[string]struct{})

	// Assuming the named groups are in the match
	vers1 := matches[pattern.SubexpIndex("v1")]
	vers2 := matches[pattern.SubexpIndex("v2")]
	ext := matches[pattern.SubexpIndex("ext")]

	// Adding the basic versions
	versionSet[vers1] = struct{}{}
	versionSet[vers2] = struct{}{}

	// Adding the extended versions if ext exists
	if ext != "" {
		versionSet[vers1+ext] = struct{}{}
		versionSet[vers2+ext] = struct{}{}
	}

	// Converting the map to a slice
	for version := range versionSet {
		versions = append(versions, version)
	}

	return versions
}

// parsePurl returns a list of permutations of pkg:version from a purl.
func parsePurl(purl string, patterns regexes) []string {
	purl = patterns.purlTrailingVersion.ReplaceAllString(purl, "$1@")
	var result []string
	var pkgs []string
	var versions []string

	if match := patterns.purlVersion.FindStringSubmatch(purl); match != nil {
		versions = parsePurlVersions(match, patterns.purlVersion)
		match = nil
	}
	if match := patterns.purlPkg.FindStringSubmatch(purl); match != nil {
		pkgs = parsePurlPkgs(match, patterns.purlPkg)
	}

	for _, i := range pkgs {
		for _, j := range versions {
			result = append(result, i+":"+j)
		}
	}
	return removeDuplicates(result)
}

// removeDuplicates removes duplicate strings from a slice.
func removeDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	var result []string

	for v := range elements {
		if encountered[elements[v]] == true {
		} else {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}

func lookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

// readSlice loads the atom reachables slice
func readSlice(sliceFile string) (ReachableSlice, error) {
	// Define a variable to hold the JSON data
	var data ReachableSlice

	// Open the JSON file
	file, err := os.Open(sliceFile)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	if err != nil {
		//fmt.Println("Error opening file:", err)
		//return data, err
	}

	// Read the file content
	byteValue, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		//return data, err
	}

	// Unmarshal the JSON data into the variable
	if err := json.Unmarshal(byteValue, &data); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return data, err
	}

	return data, nil
}

type ReachableSlice struct {
	Reachables []struct {
		Flows []struct {
			ID                    int    `json:"id"`
			Label                 string `json:"label"`
			Name                  string `json:"name"`
			FullName              string `json:"fullName"`
			Signature             string `json:"signature"`
			IsExternal            bool   `json:"isExternal"`
			Code                  string `json:"code"`
			TypeFullName          string `json:"typeFullName"`
			ParentMethodName      string `json:"parentMethodName"`
			ParentMethodSignature string `json:"parentMethodSignature"`
			ParentFileName        string `json:"parentFileName"`
			ParentPackageName     string `json:"parentPackageName"`
			ParentClassName       string `json:"parentClassName"`
			LineNumber            int    `json:"lineNumber"`
			ColumnNumber          int    `json:"columnNumber"`
			Tags                  string `json:"tags"`
		} `json:"flows"`
		Purls []string `json:"purls"`
	} `json:"reachables"`
}

func run() {
	res, err := putil.LoadTaggedToolResponse(readPath)
	if err != nil {
		log.Fatalf("could not load tool response from path %s , error:%v", readPath, err)
	}
	//if annotation == "" {
	//	annotation = defaultAnnotation
	//}
	data, err := readSlice(sliceFile)
	patterns := regexes{
		purlPkg:             regexp.MustCompile(`(?P<p1>[^/:]+/(?P<p2>[^/]+))(?:(?:.|/)v\d+)?@`),
		purlTrailingVersion: regexp.MustCompile(`[./]v\d+@`),
		purlVersion:         regexp.MustCompile(`@(?P<v1>v?(?P<v2>[\d.]+){1,3})(?P<ext>[^?\s]+)?`),
		filename:            regexp.MustCompile(`[^/]+[^/]?`),
		fileLine:            regexp.MustCompile(`(?P<file>[^/]+):(?P<line>[\d-]+)`),
	}

	for _, r := range res {
		var enrichedIssues []*v1.EnrichedIssue
		for _, i := range r.GetIssues() {
			eI, err := enrichIssue(i, data, patterns)
			if err != nil {
				log.Println(err)
				continue
			}
			enrichedIssues = append(enrichedIssues, eI)
		}
		if len(enrichedIssues) > 0 {
			if err := putil.WriteEnrichedResults(r, enrichedIssues,
				filepath.Join(writePath, fmt.Sprintf("%s.reachability.enriched.pb", r.GetToolName())),
			); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Println("no enriched issues were created for", r.GetToolName())
		}
		if len(r.GetIssues()) > 0 {
			scanStartTime := r.GetScanInfo().GetScanStartTime().AsTime()
			if err := putil.WriteResults(
				r.GetToolName(),
				r.GetIssues(),
				filepath.Join(writePath, fmt.Sprintf("%s.raw.pb", r.GetToolName())),
				r.GetScanInfo().GetScanUuid(),
				scanStartTime,
				r.GetScanInfo().GetScanTags(),
			); err != nil {
				log.Fatalf("could not write results: %s", err)
			}
		}

	}
}

func main() {
	flag.StringVar(&readPath, "read_path", lookupEnvOrString("READ_PATH", ""), "where to find producer results")
	flag.StringVar(&writePath, "write_path", lookupEnvOrString("WRITE_PATH", ""), "where to put enriched results")
	flag.StringVar(&annotation, "annotation", lookupEnvOrString("ANNOTATION", defaultAnnotation), "what is the annotation this enricher will add to the issues, by default `Reachability`")
	flag.StringVar(&sliceFile, "atom_slice", lookupEnvOrString("ATOM_SLICE", ""), "location of the atom slice file")
	flag.Parse()
	run()
}
