package types

// Position represents where in the file the finding is located.
type Position struct {
	Col  int `json:"col"`
	Line int `json:"line"`
}

// Extra contains extra info needed for semgrep issue.
type Extra struct {
	Message  string   `json:"message"`
	Metavars Metavars `json:"metavars"`
	Metadata Metadata `json:"metadata"`
	Severity string   `json:"severity"`
	Lines    string   `json:"lines"`
}

// SemgrepIssue represents a semgrep issue.
type SemgrepIssue struct {
	CheckID string   `json:"check_id"`
	Path    string   `json:"path"`
	Start   Position `json:"start"`
	End     Position `json:"end"`
	Extra   Extra    `json:"extra"`
}

// SemgrepResults represents a series of semgrep issues.
type SemgrepResults struct {
	Results []SemgrepIssue `'json:"results"`
}

// Metavars currently is empty but could represent more metavariables for semgrep.
type Metavars struct{}

// Metadata currently is empty, however, could represent semgrep issue metadata going forward.
type Metadata struct{}
