package purl

import (
	"fmt"
	"regexp"
	"strings"
)

// Parser allows to extract information from purls - https://github.com/package-url/purl-spec.
type Parser struct {
	matcherPurlPkg             *regexp.Regexp
	matcherPurlTrailingVersion *regexp.Regexp
	matcherPurlVersion         *regexp.Regexp
}

func NewParser() (*Parser, error) {
	purlPkg, err := regexp.Compile(`(?P<p1>[^/:]+/(?P<p2>[^/]+))(?:(?:.|/)v\d+)?@`)
	if err != nil {
		return nil, fmt.Errorf("failed to compile purl pkg regex: %w", err)
	}
	purlTrailingVersion, err := regexp.Compile(`[./]v\d+@`)
	if err != nil {
		return nil, fmt.Errorf("failed to compile purl trailing version regex: %w", err)
	}
	purlVersion, err := regexp.Compile(`@(?P<v1>v?(?P<v2>[\d.]+){1,3})(?P<ext>[^?\s]+)?`)
	if err != nil {
		return nil, fmt.Errorf("failed to compile purl version regex: %w", err)
	}

	return &Parser{
		matcherPurlPkg:             purlPkg,
		matcherPurlTrailingVersion: purlTrailingVersion,
		matcherPurlVersion:         purlVersion,
	}, nil
}

// ParsePurl extracts pkg:version matches from the supplied purl.
func (p *Parser) ParsePurl(purl string) ([]string, error) {
	purl = p.matcherPurlTrailingVersion.ReplaceAllString(purl, "$1@")

	var (
		result   []string
		pkgs     []string
		versions []string
	)

	if match := p.matcherPurlVersion.FindStringSubmatch(purl); len(match) > 0 {
		versions = p.parsePurlVersions(match)
	}

	if match := p.matcherPurlPkg.FindStringSubmatch(purl); len(match) > 0 {
		pkgs = p.parsePurlPkgs(match)
	}

	for _, pkg := range pkgs {
		for _, version := range versions {
			result = append(result, fmt.Sprintf("%s:%s", pkg, version))
		}
	}

	return p.removeDuplicates(result), nil
}

func (p *Parser) parsePurlVersions(matches []string) []string {
	if len(matches) == 0 {
		return make([]string, 0)
	}

	var (
		pattern  = p.matcherPurlVersion
		versions []string
		// Creating a map to ensure uniqueness
		versionSet = make(map[string]struct{})

		// Assuming the named groups are in the match
		vers1 = matches[pattern.SubexpIndex("v1")]
		vers2 = matches[pattern.SubexpIndex("v2")]
		ext   = matches[pattern.SubexpIndex("ext")]
	)

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

func (p *Parser) parsePurlPkgs(matches []string) []string {
	var (
		pattern = p.matcherPurlPkg
		// Creating a map to ensure uniqueness
		pkgSet         = make(map[string]struct{})
		pkgs           []string
		pkgStrReplacer = strings.NewReplacer(
			// replaces "pypi/" with "".
			"pypi/", "",
			// replaces "npm/" with "".
			"npm/", "",
			// replaces "%40/" with "@".
			"%40", "@",
		)
	)

	// Adding the packages
	pkgSet[matches[pattern.SubexpIndex("p1")]] = struct{}{}
	pkgSet[matches[pattern.SubexpIndex("p2")]] = struct{}{}

	// Converting the map to a slice and cleaning up the packages
	for pkg := range pkgSet {
		pkgs = append(pkgs, pkgStrReplacer.Replace(pkg))
	}

	return pkgs
}

func (p *Parser) removeDuplicates(matches []string) []string {
	var (
		result      []string
		encountered = make(map[string]struct{})
	)

	for match := range matches {
		_, ok := encountered[matches[match]]
		if ok {
			continue
		}
		encountered[matches[match]] = struct{}{}
		result = append(result, matches[match])
	}

	return result
}
