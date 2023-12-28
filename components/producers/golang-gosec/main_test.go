package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	v1 "github.com/ocurity/dracon/api/proto/v1"

	"github.com/stretchr/testify/assert"
)

var code = `func GetProducts(ctx context.Context, db *sql.DB, category string) ([]Product, error) {
	rows, err := db.QueryContext(ctx, "SELECT * FROM product WHERE category='"+category+"'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.Id, &product.Name, &product.Category, &product.Price); err != nil {
			return nil, err
		}`

func TestParseIssues(t *testing.T) {
	file, err := os.CreateTemp("", "dracon_context_test")
	if err != nil {
		t.Errorf("could not setup tests for context pkg, could not create temporary files")
	}
	defer os.Remove(file.Name())
	if err := os.WriteFile(file.Name(), []byte(code), os.ModeAppend); err != nil {
		t.Errorf("could not setup tests for context pk, could not write temporary file")
	}

	exampleOutput := fmt.Sprintf(`
{
	"Issues": [
		{
			"severity": "MEDIUM",
			"confidence": "HIGH",
			"rule_id": "G304",
			"details": "Potential file inclusion via variable",
			"file": "%s",
			"code": "ioutil.ReadFile(path)",
			"line": "2",
			"column": "44"
		}
	],
	"Stats": {
		"files": 1,
		"lines": 60,
		"nosec": 0,
		"found": 1
	}
}`, file.Name())
	var results GoSecOut
	err = json.Unmarshal([]byte(exampleOutput), &results)
	assert.Nil(t, err)

	issues, err := parseIssues(&results)
	assert.Nil(t, err)
	expectedIssue := &v1.Issue{
		Target:         fmt.Sprintf("%s:2", file.Name()),
		Type:           "G304",
		Title:          "Potential file inclusion via variable",
		Severity:       v1.Severity_SEVERITY_MEDIUM,
		Cvss:           0.0,
		Confidence:     v1.Confidence_CONFIDENCE_HIGH,
		Description:    "ioutil.ReadFile(path)",
		ContextSegment: &code,
	}

	assert.Equal(t, []*v1.Issue{expectedIssue}, issues)
}
