package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	owners "github.com/hairyhenderson/go-codeowners"
	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/pkg/putil"
)

const defaultAnnotation = "Owner"

var (
	readPath       string
	writePath      string
	repoBasePath   string
	repoURL        string
	sparseCheckout string
	annotation     string
)

func lookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func enrichIssue(i *v1.Issue) (*v1.EnrichedIssue, error) {
	enrichedIssue := v1.EnrichedIssue{}
	annotations := map[string]string{}
	if i.GetCycloneDXSBOM() != "" {
		annotationsAdded := 0
		// shortcut, if there is a CycloneDX BOM then there is no target.
		// we get the url from the repoURL parameter and add the sparse checkout directories to it.
		targets := []string{}
		if sparseCheckout != "" {
			targets = strings.Split(sparseCheckout, ",")
		} else {
			targets = []string{"."}
		}
		for _, target := range targets {
			path := filepath.Join(repoBasePath, target)
			c, err := owners.FromFile()
			if err != nil {
				log.Println("could not instantiate owners for path", path, "err", err)
				continue
			}
			owners := c.Owners(".")
			for _, owner := range owners {
				annotations[fmt.Sprintf("Owner-%d", annotationsAdded)] = owner
				annotationsAdded++
			}
		}
	} else {
		target := strings.Split(i.GetTarget(), ":")
		annotationsAdded := 0
		if len(target) > 1 {
			c, err := owners.FromFile(filepath.Join(repoBasePath, target[0]))
			path := filepath.Join(repoBasePath, target[0])
			if err != nil {
				log.Println("could not instantiate owners for path", path, "err", err)
			}
			owners := c.Owners(".")
			for _, owner := range owners {
				annotations[fmt.Sprintf("Owner-%d", annotationsAdded)] = owner
				annotationsAdded++
			}
		}
	}
	enrichedIssue = v1.EnrichedIssue{
		RawIssue:    i,
		Annotations: annotations,
	}
	enrichedIssue.Annotations = annotations
	return &enrichedIssue, nil
}

func run() {
	res, err := putil.LoadTaggedToolResponse(readPath)
	if err != nil {
		log.Fatalf("could not load tool response from path %s , error:%v", readPath, err)
	}
	if annotation == "" {
		annotation = defaultAnnotation
	}
	for _, r := range res {
		enrichedIssues := []*v1.EnrichedIssue{}
		for _, i := range r.GetIssues() {
			eI, err := enrichIssue(i)
			if err != nil {
				log.Println(err)
				continue
			}
			enrichedIssues = append(enrichedIssues, eI)
		}
		if len(enrichedIssues) > 0 {
			if err := putil.WriteEnrichedResults(r, enrichedIssues,
				filepath.Join(writePath, fmt.Sprintf("%s.depsdev.enriched.pb", r.GetToolName())),
			); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Println("no enriched issues were created")
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
}

func main() {
	flag.StringVar(&readPath, "read_path", lookupEnvOrString("READ_PATH", ""), "where to find producer results")
	flag.StringVar(&writePath, "write_path", lookupEnvOrString("WRITE_PATH", ""), "where to put enriched results")
	flag.StringVar(&annotation, "annotation", lookupEnvOrString("ANNOTATION", defaultAnnotation), "what is the annotation this enricher will add to the issues, by default `Enriched Licenses`")
	flag.StringVar(&repoBasePath, "repoBasePath", lookupEnvOrString("REPO_BASE_PATH", ""), `the base path of the repository, this is most likely an internally set variable`)
	flag.StringVar(&repoURL, "repoURL", lookupEnvOrString("REPOSITORY_URL", ""), `the base path of the repository, this is most likely an internally set variable`)
	flag.StringVar(&sparseCheckout, "sparseCheckout", lookupEnvOrString("REPO_SPARSE_CHECKOUT", ""), `the sparse checkout directories settings passed to this repository`)
	flag.Parse()
	run()
}
