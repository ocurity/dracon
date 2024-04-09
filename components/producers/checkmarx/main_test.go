package main

import (
	"encoding/json"
	"encoding/xml"
	"testing"

	v1 "github.com/ocurity/dracon/api/proto/v1"

	"github.com/stretchr/testify/assert"
)

func TestParseIssues(t *testing.T) {
	exampleOutput := checkmarxout
	var results Flaws
	err := xml.Unmarshal([]byte(exampleOutput), &results)
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
			Target:      "/some/target:2",
			Type:        "209",
			Title:       "Recurrent - High - WebGoat",
			Severity:    v1.Severity_SEVERITY_HIGH,
			Confidence:  v1.Confidence_CONFIDENCE_UNSPECIFIED,
			Source:      "165072:WebgoatMay5:WebGoat",
			Description: string(d0),
		},
		{
			Target:      "/some/target:2",
			Type:        "210",
			Title:       "Recurrent - High - WebGoat",
			Severity:    v1.Severity_SEVERITY_HIGH,
			Cvss:        0,
			Confidence:  v1.Confidence_CONFIDENCE_UNSPECIFIED,
			Source:      "165072:WebgoatMay5:WebGoat",
			Description: string(d1),
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
            - File: /some/target, Line:2, Column:
            22, Name:t, Code: } catch (Throwable t)
            - File: /some/target, Line:3, Column: 4,
            Name:t, Code: t.printStackTrace();
            - File: /some/target, Line:4, Column:
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
        <fileName>/some/target</fileName>
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
		- File: /some/target, Line:2, Column:
		23, Name:thr, Code: } catch (Throwable thr)
		- File: /some/target, Line:3, Column: 5,
		Name:thr, Code: thr.printStackTrace();
		- File: /some/target, Line:4, Column:
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
	<fileName>/some/target</fileName>
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
