package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	cdx "github.com/CycloneDX/cyclonedx-go"
	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/pkg/cyclonedx"
	"github.com/ocurity/dracon/pkg/putil"
)

var (
	readPath       string
	writePath      string
	depsdevBaseURL = "https://deps.dev"
)

type Version struct {
	Version                string        `json:"version,omitempty"`
	SymbolicVersions       []interface{} `json:"symbolicVersions,omitempty"`
	RefreshedAt            int           `json:"refreshedAt,omitempty"`
	IsDefault              bool          `json:"isDefault,omitempty"`
	Licenses               []string      `json:"licenses,omitempty"`
	DependentCount         int           `json:"dependentCount,omitempty"`
	DependentCountDirect   int           `json:"dependentCountDirect,omitempty"`
	DependentCountIndirect int           `json:"dependentCountIndirect,omitempty"`
	Links                  struct {
		Origins []string `json:"origins,omitempty"`
	} `json:"links,omitempty"`
	Projects []struct {
		Type        string `json:"type,omitempty"`
		Name        string `json:"name,omitempty"`
		ObservedAt  int    `json:"observedAt,omitempty"`
		Issues      int    `json:"issues,omitempty"`
		Forks       int    `json:"forks,omitempty"`
		Stars       int    `json:"stars,omitempty"`
		Description string `json:"description,omitempty"`
		License     string `json:"license,omitempty"`
		DisplayName string `json:"displayName,omitempty"`
		Link        string `json:"link,omitempty"`
		ScorecardV2 struct {
			Date string `json:"date,omitempty"`
			Repo struct {
				Name   string `json:"name,omitempty"`
				Commit string `json:"commit,omitempty"`
			} `json:"repo,omitempty"`
			Scorecard struct {
				Version string `json:"version,omitempty"`
				Commit  string `json:"commit,omitempty"`
			} `json:"scorecard,omitempty"`
			Check []struct {
				Name          string `json:"name,omitempty"`
				Documentation struct {
					Short string `json:"short,omitempty"`
					URL   string `json:"url,omitempty"`
				} `json:"documentation,omitempty"`
				Score   int           `json:"score,omitempty"`
				Reason  string        `json:"reason,omitempty"`
				Details []interface{} `json:"details,omitempty"`
			} `json:"check,omitempty"`
			Metadata []interface{} `json:"metadata,omitempty"`
			Score    float64       `json:"score,omitempty"`
		} `json:"scorecardV2,omitempty"`
	} `json:"projects,omitempty"`
	Advisories      []interface{} `json:"advisories,omitempty"`
	RelatedPackages struct{}      `json:"relatedPackages,omitempty"`
}
type Response struct {
	Package struct {
		System string `json:"system,omitempty"`
		Name   string `json:"name,omitempty"`
	} `json:"package,omitempty"`
	Owners         []interface{} `json:"owners,omitempty"`
	Version        Version       `json:"version,omitempty"`
	DefaultVersion string        `json:"defaultVersion,omitempty"`
}

func lookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func enrichIssue(i *v1.Issue) (*v1.EnrichedIssue, error) {
	enrichedIssue := v1.EnrichedIssue{}
	annotations := map[string]string{}
	bom, err := cyclonedx.FromDracon(i)
	if err != nil {
		return &enrichedIssue, err
	}
	var depsResp Response

	for index, component := range *bom.Components {
		licenses := cdx.Licenses{}
		if component.Type == cdx.ComponentTypeLibrary {
			if component.Licenses == nil {
				resp, err := http.Get(fmt.Sprintf("%s/_/s/go/p/%s/v/%s", depsdevBaseURL, url.QueryEscape(component.Name), url.QueryEscape(component.Version)))
				if err != nil {
					log.Println(err)
					continue
				}
				err = json.NewDecoder(resp.Body).Decode(&depsResp)
				if err != nil {
					log.Println(err)
					continue
				}
				for _, lic := range depsResp.Version.Licenses {
					log.Println(lic)

					bomLicense := cdx.License{
						Name: lic,
					}
					licenses = append(licenses, cdx.LicenseChoice{License: &bomLicense})
				}
				(*bom.Components)[index].Licenses = &licenses
				annotations["Enriched Licenses"] = "True"
			}
			// TODO(): enrich with vulnerability and scorecard info whenever a consumer supports showing arbitrary properties in components
		}
	}

	marshalled, err := json.Marshal(bom)
	if err != nil {
		return &enrichedIssue, err
	}
	originalIssue, err := cyclonedx.ToDracon(marshalled, "json")
	if err != nil {
		return &enrichedIssue, err
	}
	enrichedIssue = v1.EnrichedIssue{
		RawIssue:    originalIssue[0],
		Annotations: map[string]string{},
	}
	enrichedIssue.Annotations = annotations
	return &enrichedIssue, nil
}

func run() {
	res, err := putil.LoadTaggedToolResponse(readPath)
	if err != nil {
		log.Fatalf("could not load tool response from path %s , error:%v", readPath, err)
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
	flag.Parse()
	run()
}
