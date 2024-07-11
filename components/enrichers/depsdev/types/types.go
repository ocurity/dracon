package types

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
