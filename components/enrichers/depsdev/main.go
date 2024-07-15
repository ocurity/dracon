package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"

	cdx "github.com/CycloneDX/cyclonedx-go"
	packageurl "github.com/package-url/packageurl-go"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components/enrichers"
	"github.com/ocurity/dracon/pkg/cyclonedx"
)

const defaultAnnotation = "Enriched Licenses"

var (
	depsdevBaseURL     = "https://deps.dev"
	licensesInEvidence string
	scoreCardInfo      string
	annotation         string
)

// Check is a deps.dev ScoreCardV2 check
type Check struct {
	Name          string `json:"name,omitempty"`
	Documentation struct {
		Short string `json:"short,omitempty"`
		URL   string `json:"url,omitempty"`
	} `json:"documentation,omitempty"`
	Score   int           `json:"score,omitempty"`
	Reason  string        `json:"reason,omitempty"`
	Details []interface{} `json:"details,omitempty"`
}

// ScorecardV2 is a deps.dev ScoreCardV2 result
type ScorecardV2 struct {
	Date string `json:"date,omitempty"`
	Repo struct {
		Name   string `json:"name,omitempty"`
		Commit string `json:"commit,omitempty"`
	} `json:"repo,omitempty"`
	Scorecard struct {
		Version string `json:"version,omitempty"`
		Commit  string `json:"commit,omitempty"`
	} `json:"scorecard,omitempty"`
	Check    []Check       `json:"check,omitempty"`
	Metadata []interface{} `json:"metadata,omitempty"`
	Score    float64       `json:"score,omitempty"`
}

// Project is a deps.dev project
type Project struct {
	Type        string      `json:"type,omitempty"`
	Name        string      `json:"name,omitempty"`
	ObservedAt  int         `json:"observedAt,omitempty"`
	Issues      int         `json:"issues,omitempty"`
	Forks       int         `json:"forks,omitempty"`
	Stars       int         `json:"stars,omitempty"`
	Description string      `json:"description,omitempty"`
	License     string      `json:"license,omitempty"`
	DisplayName string      `json:"displayName,omitempty"`
	Link        string      `json:"link,omitempty"`
	ScorecardV2 ScorecardV2 `json:"scorecardV2,omitempty"`
}

// Version is a deps.dev version, main object in the response
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
	Projects        []Project     `json:"projects,omitempty"`
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

func makeURL(component cdx.Component, api bool) (string, error) {
	instance, err := packageurl.FromString(component.PackageURL)
	if err != nil {
		return "", err
	}
	ecosystem := ""
	version := url.QueryEscape(component.Version)
	switch instance.Type {
	case packageurl.TypeGolang:
		ecosystem += "/go"
		version = "v" + version
	case packageurl.TypePyPi:
		ecosystem += "/pypi"
	case packageurl.TypeMaven:
		ecosystem += "/maven"
	// case packageurl.TypeCargo:
	// 	ecosystem += "/cargo"
	case packageurl.TypeNPM:
		ecosystem += "/npm"
	case packageurl.TypeNuget:
		ecosystem += "/nuget"
	default:
		log.Println(instance.Namespace, "not supported by this enricher")
	}
	resultURL := ""
	if api {
		resultURL = fmt.Sprintf("%s/_/s%s/p/%s/v/%s", depsdevBaseURL, ecosystem, url.QueryEscape(component.Name), version)
	} else {
		resultURL = fmt.Sprintf("%s%s/p/%s/v/%s", depsdevBaseURL, ecosystem, url.QueryEscape(component.Name), version)
	}
	return resultURL, nil
}

func addDepsDevLink(component cdx.Component) (cdx.Component, error) {
	url, err := makeURL(component, false)
	if err != nil {
		return component, err
	}
	depsDevRef := cdx.ExternalReference{
		Type: cdx.ERTypeOther,
		URL:  url,
	}

	if component.ExternalReferences != nil && len(*component.ExternalReferences) > 0 {
		refs := append(*component.ExternalReferences, depsDevRef)
		component.ExternalReferences = &refs
	} else {
		refs := []cdx.ExternalReference{depsDevRef}
		component.ExternalReferences = &refs
	}

	return component, nil
}

func addDepsDevInfo(component cdx.Component, annotations map[string]string) (cdx.Component, map[string]string, error) {
	var depsResp Response
	licenses := cdx.Licenses{}
	url, err := makeURL(component, true)
	if err != nil {
		return component, annotations, err
	}
	resp, err := http.Get(url) // nolint: gosec, url get constructed above with a hardcoded domain and relatively trusted data
	if err != nil {
		return component, annotations, err
	}
	err = json.NewDecoder(resp.Body).Decode(&depsResp)
	if err != nil {
		return component, annotations, err
	}
	if len(depsResp.Version.Licenses) == 0 {
		log.Println("could not find license for component", component.Name)
	}

	for _, lic := range depsResp.Version.Licenses {
		licenseName := cdx.License{
			Name: lic,
		}
		licenses = append(licenses, cdx.LicenseChoice{License: &licenseName})
		log.Println("found license", lic, "for component", component.Name)
	}
	if scoreCardInfo == "true" {
		log.Println("adding scorecard info")
		for _, project := range depsResp.Version.Projects {
			if project.ScorecardV2.Date != "" && len(project.ScorecardV2.Check) != 0 && project.ScorecardV2.Score >= 0 {
				scoreCardInfo, err := json.MarshalIndent(project.ScorecardV2, "", "\t")
				if err != nil {
					log.Println("could not marshal score card information, err:", err)
					continue
				}
				properties := []cdx.Property{
					{
						Name:  "ScorecardScore",
						Value: fmt.Sprintf("%f", project.ScorecardV2.Score),
					},
					{
						Name:  "ScorecardInfo",
						Value: string(scoreCardInfo),
					},
				}
				props := append(*component.Properties, properties...)
				component.Properties = &props
			}
		}
	}
	if licensesInEvidence == "true" {
		log.Println("adding Licenses in the 'Evidence' field")
		evid := cdx.Evidence{
			Licenses: &licenses,
		}
		if component.Evidence == nil {
			component.Evidence = &evid
		} else {
			component.Evidence.Licenses = &licenses
		}
	} else {
		component.Licenses = &licenses
	}
	annotations[annotation] = "True"
	return component, annotations, nil
}

func enrichIssue(i *v1.Issue) (*v1.EnrichedIssue, error) {
	enrichedIssue := v1.EnrichedIssue{}
	annotations := map[string]string{}
	bom, err := cyclonedx.FromDracon(i)
	if err != nil {
		return &enrichedIssue, err
	}
	if bom == nil || bom.Components == nil {
		return &enrichedIssue, errors.New("bom does not have components")
	}
	newComponents := (*bom.Components)[:0]
	for _, component := range *bom.Components {
		newComp := component
		if component.Type == cdx.ComponentTypeLibrary {
			if component.Licenses == nil {
				newComp, annotations, err = addDepsDevInfo(component, annotations)
				if err != nil {
					log.Println(err)
					continue
				}
			}
			newComp, err = addDepsDevLink(newComp)
			if err != nil {
				log.Println(err)
				continue
			}
			// TODO(): enrich with vulnerability info whenever a consumer supports showing arbitrary properties in components
		}
		newComponents = append(newComponents, newComp)
	}
	bom.Components = &newComponents
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

func run() error {
	res, err := enrichers.LoadData()
	if err != nil {
		return err
	}
	if annotation == "" {
		annotation = defaultAnnotation
	}
	for _, r := range res {
		enrichedIssues := []*v1.EnrichedIssue{}
		for _, i := range r.GetIssues() {
			eI, err := enrichIssue(i)
			if err != nil {
				slog.Error(err.Error())
				continue
			}
			enrichedIssues = append(enrichedIssues, eI)
		}
		return enrichers.WriteData(&v1.EnrichedLaunchToolResponse{
			OriginalResults: r,
			Issues:          enrichedIssues,
		}, "deps-dev")
	}
	return nil
}

func main() {
	flag.StringVar(&annotation, "annotation", enrichers.LookupEnvOrString("ANNOTATION", defaultAnnotation), "what is the annotation this enricher will add to the issues, by default `Enriched Licenses`")
	flag.StringVar(&scoreCardInfo, "scoreCardInfo", enrichers.LookupEnvOrString("SCORECARD_INFO", "false"), "add security score card scan results from deps.dev to the components of the SBOM as properties")
	flag.StringVar(&licensesInEvidence, "licensesInEvidence", enrichers.LookupEnvOrString("LICENSES_IN_EVIDENCE", ""),
		`If this flag is provided and set to "true", the enricher will populate the 'evidence' CycloneDX field with license information instead of the license field.
	This means that the result conforms to the CycloneDX intention of providing accurate information when licensing information cannot be guaranteed to be accurate.
	However, no tools currently support reading license information from evidence.
	This is because deps.dev does not guarantee accurate licensing information for Go.
	Enable this switch if you need to provide SBOM information for regulatory reasons.`)
	if err := enrichers.ParseFlags(); err != nil {
		log.Fatal(err)
	}
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
