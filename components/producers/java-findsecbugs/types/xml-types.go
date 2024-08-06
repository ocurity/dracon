package types

import "encoding/xml"

type jar struct {
	XMLName xml.Name `xml:"Jar"`
}

type project struct {
	XMLName xml.Name `xml:"Project"`
	Jar     *jar     `xml:"Jar"`
}

type method struct {
	XMLName    xml.Name     `xml:"Method"`
	Classname  string       `xml:"classname,attr"`
	Name       string       `xml:"name,attr"`
	Signature  string       `xml:"signature,attr"`
	IsStatic   string       `xml:"isStatic,attr"`
	SourceLine []SourceLine `xml:"SourceLine"`
}
type SourceLine struct {
	XMLName       xml.Name `xml:"SourceLine"`
	Classname     string   `xml:"classname,attr"`
	Start         int      `xml:"start,attr"`
	End           int      `xml:"end,attr"`
	StartBytecode string   `xml:"startBytecode,attr"`
	EndBytecode   string   `xml:"endBytecode,attr"`
	Sourcefile    string   `xml:"sourcefile,attr"`
	Sourcepath    string   `xml:"sourcepath,attr"`
	Role          string   `xml:"role,attr"`
}
type class struct {
	XMLName    xml.Name     `xml:"Class"`
	Classname  string       `xml:"classname,attr"`
	Role       string       `xml:"role,attr"`
	SourceLine []SourceLine `xml:"SourceLine"`
}
type field struct {
	XMLName    xml.Name     `xml:"Field"`
	Classname  string       `xml:"classname,attr"`
	SourceLine []SourceLine `xml:"SourceLine"`
}
type BugInstance struct {
	XMLName      xml.Name     `xml:"BugInstance"`
	Class        []class      `xml:"Class"`
	Method       []method     `xml:"Method"`
	SourceLine   []SourceLine `xml:"SourceLine"`
	Field        []field      `xml:"Field"`
	LongMessage  string       `xml:"LongMessage"`
	ShortMessage string       `xml:"ShortMessage"`
	Type         string       `xml:"type,attr"`
	Priority     string       `xml:"priority,attr"`
	Rank         string       `xml:"rank,attr"`
	Abbrev       string       `xml:"abbrev,attr"`
	Category     string       `xml:"category,attr"`
}
type BugCollection struct {
	XMLName     xml.Name      `xml:"BugCollection"`
	Project     *project      `xml:"Project"`
	BugInstance []BugInstance `xml:"BugInstance"`
}
