package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"testing"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/pkg/testutil"

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
	f, err := testutil.CreateFile("checkmarx_tests_vuln_code", code)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(f.Name())
	exampleOutput := fmt.Sprintf(checkmarxout, f.Name(), f.Name(), f.Name(), f.Name(), f.Name(), f.Name(), f.Name(), f.Name(), f.Name(), f.Name(), f.Name(), f.Name(), f.Name())
	var results Flaws
	err = xml.Unmarshal([]byte(exampleOutput), &results)
	assert.Nil(t, err)

	issues, err := parseIssues(&results)
	assert.Nil(t, err)

	desc0 := DraconDescription{
		OriginalIssueDescription: results.Flaw[0].IssueDescription,
	}
	d0, _ := json.Marshal(desc0)

	desc1 := DraconDescription{
		OriginalIssueDescription: results.Flaw[1].IssueDescription,
	}
	d1, _ := json.Marshal(desc1)
	expectedIssues := []*v1.Issue{
		{
			Target:         fmt.Sprintf("%s:2", f.Name()),
			ContextSegment: &code,
			Type:           "209",
			Title:          "Recurrent - High - WebGoat",
			Severity:       v1.Severity_SEVERITY_HIGH,
			Confidence:     v1.Confidence_CONFIDENCE_UNSPECIFIED,
			Source:         "165072:WebgoatMay5:WebGoat",
			Description:    string(d0),
		},
		{
			Target:         fmt.Sprintf("%s:2", f.Name()),
			Type:           "210",
			Title:          "Recurrent - High - WebGoat",
			Severity:       v1.Severity_SEVERITY_HIGH,
			Cvss:           0,
			Confidence:     v1.Confidence_CONFIDENCE_UNSPECIFIED,
			Source:         "165072:WebgoatMay5:WebGoat",
			Description:    string(d1),
			ContextSegment: &code,
		},
	}
	assert.Equal(t, expectedIssues, issues)
}

var checkmarxout = `
<?xml version="1.0" encoding="UTF-8"?>

<flaws>
    <metaData
        appID="165072" appName="WebgoatMay5" componentName="WebGoat"
        date="Monday, January 6, 2020 5:00:34 AM" releaseName="CX" sourceName="Checkmarx"
        sourceDesc="Checkmarx" />

    <flaw>
        <id>165072-WebGoat-Checkmarx-1539227441</id>
        <status>Recurrent</status>
        <issueDescription>
            - File: %s, Line:2, Column:
            22, Name:t, Code: } catch (Throwable t)
            - File: %s, Line:3, Column: 4,
            Name:t, Code: t.printStackTrace();
            - File: %s, Line:4, Column:
            21, Name:printStackTrace, Code: t.printStackTrace();</issueDescription>
        <remediationDesc />
        <exploitDesc />
        <issueRecommendation />
        <componentName>WebGoat</componentName>
        <module />
        <apiName />
        <vulnerabilityType>209</vulnerabilityType>
        <classification />
        <severity>High</severity>
        <fileName>%s</fileName>
        <lineNumber>2</lineNumber>
        <srcContext>WebGoatProd/JavaSource/org/owasp/webgoat/</srcContext>
        <defectInfo />
        <notes />
        <trace />
        <callerName />
        <findingCodeRegion />
        <dateFirstOccurrence />
        <issueBornDate />
        <issueName />
        <kBReference />
        <cVSSScore />
        <relatedExploitRange />
        <attackComplexity />
        <levelofAuthenticationNeeded />
        <confidentialityImpact />
        <integrityImpact />
        <availabilityImpact />
        <collateralDamagePotential />
        <targetDistribution />
        <confidentialityRequirement />
        <integrityRequirement />
        <availabilityRequirement />
        <availabilityofExploit />
        <typeofFixAvailable />
        <levelofVerificationthatVulnerabilityExist />
        <cVSSEquation />
    </flaw>
  
	<flaw>
	<id>165072-WebGoat-Checkmarx--366645893</id>
	<status>Recurrent</status>
	<issueDescription>
		- File: %s, Line:2, Column:
		23, Name:thr, Code: } catch (Throwable thr)
		- File: %s, Line:3, Column: 5,
		Name:thr, Code: thr.printStackTrace();
		- File: %s, Line:4, Column:
		24, Name:printStackTrace, Code: thr.printStackTrace();</issueDescription>
	<remediationDesc />
	<exploitDesc />
	<issueRecommendation />
	<componentName>WebGoat</componentName>
	<module />
	<apiName />
	<vulnerabilityType>210</vulnerabilityType>
	<classification />
	<severity>High</severity>
	<fileName>%s</fileName>
	<lineNumber>2</lineNumber>
	<srcContext>WebGoatProd/JavaSource/org/owasp/webgoat/</srcContext>
	<defectInfo />
	<notes />
	<trace />
	<callerName />
	<findingCodeRegion />
	<dateFirstOccurrence />
	<issueBornDate />
	<issueName />
	<kBReference />
	<cVSSScore />
	<relatedExploitRange />
	<attackComplexity />
	<levelofAuthenticationNeeded />
	<confidentialityImpact />
	<integrityImpact />
	<availabilityImpact />
	<collateralDamagePotential />
	<targetDistribution />
	<confidentialityRequirement />
	<integrityRequirement />
	<availabilityRequirement />
	<availabilityofExploit />
	<typeofFixAvailable />
	<levelofVerificationthatVulnerabilityExist />
	<cVSSEquation />
</flaw>

    </flaws>`
