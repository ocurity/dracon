package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components/producers/terraform-tfsec/types"
)

func TestParseOut(t *testing.T) {
	var results types.TfSecOut
	err := json.Unmarshal([]byte(exampleOutput), &results)
	if err != nil {
		t.Logf(err.Error())
		t.Fail()
	}
	issues := parseOut(results)
	expectedIssues := []*v1.Issue{
		{
			Target:      "/home/foobar/go/pkg/mod/github.com/aquasecurity/tfsec@v1.28.1/_examples/main.tf:41-41",
			Type:        "aws-api-gateway-use-secure-tls-policy",
			Title:       "API Gateway domain name uses outdated SSL/TLS protocols.",
			Severity:    4,
			Confidence:  3,
			Description: "{\"rule_id\":\"AVD-AWS-0005\",\"long_id\":\"aws-api-gateway-use-secure-tls-policy\",\"rule_description\":\"API Gateway domain name uses outdated SSL/TLS protocols.\",\"rule_provider\":\"aws\",\"rule_service\":\"api-gateway\",\"impact\":\"Outdated SSL policies increase exposure to known vulnerabilities\",\"resolution\":\"Use the most modern TLS/SSL policies available\",\"links\":[\"https://aquasecurity.github.io/tfsec/latest/checks/aws/api-gateway/use-secure-tls-policy/\",\"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_domain_name#security_policy\"],\"description\":\"Domain name is configured with an outdated TLS policy.\",\"severity\":\"HIGH\",\"warning\":false,\"status\":0,\"resource\":\"aws_api_gateway_domain_name.outdated_security_policy\",\"location\":{\"filename\":\"/home/foobar/go/pkg/mod/github.com/aquasecurity/tfsec@v1.28.1/_examples/main.tf\",\"start_line\":41,\"end_line\":41}}",
		},
		{
			Target:      "/home/foobar/go/pkg/mod/github.com/aquasecurity/tfsec@v1.28.1/_examples/main.tf:37-37",
			Type:        "aws-api-gateway-use-secure-tls-policy",
			Title:       "API Gateway domain name uses outdated SSL/TLS protocols.",
			Severity:    4,
			Confidence:  3,
			Description: "{\"rule_id\":\"AVD-AWS-0005\",\"long_id\":\"aws-api-gateway-use-secure-tls-policy\",\"rule_description\":\"API Gateway domain name uses outdated SSL/TLS protocols.\",\"rule_provider\":\"aws\",\"rule_service\":\"api-gateway\",\"impact\":\"Outdated SSL policies increase exposure to known vulnerabilities\",\"resolution\":\"Use the most modern TLS/SSL policies available\",\"links\":[\"https://aquasecurity.github.io/tfsec/latest/checks/aws/api-gateway/use-secure-tls-policy/\",\"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_domain_name#security_policy\"],\"description\":\"Domain name is configured with an outdated TLS policy.\",\"severity\":\"HIGH\",\"warning\":false,\"status\":0,\"resource\":\"aws_api_gateway_domain_name.empty_security_policy\",\"location\":{\"filename\":\"/home/foobar/go/pkg/mod/github.com/aquasecurity/tfsec@v1.28.1/_examples/main.tf\",\"start_line\":37,\"end_line\":37}}",
		},
	}

	found := 0
	assert.Equal(t, 2, len(issues))
	for _, issue := range issues {
		singleMatch := 0
		for _, expected := range expectedIssues {
			if expected.Target == issue.Target {
				singleMatch++
				found++
				assert.Equal(t, singleMatch, 1) // assert no duplicates
				assert.EqualValues(t, expected.Type, issue.Type)
				assert.EqualValues(t, expected.Title, issue.Title)
				assert.EqualValues(t, expected.Severity, issue.Severity)
				assert.EqualValues(t, expected.Cvss, issue.Cvss)
				assert.EqualValues(t, expected.Confidence, issue.Confidence)
				assert.EqualValues(t, expected.Description, issue.Description)
			}
		}
	}
	assert.Equal(t, found, len(issues)) // assert everything has been found
}

var exampleOutput = `
{
	"results": [
		{
			"rule_id": "AVD-AWS-0005",
			"long_id": "aws-api-gateway-use-secure-tls-policy",
			
			"rule_description": "API Gateway domain name uses outdated SSL/TLS protocols.",
			
			"rule_provider": "aws",
			"rule_service": "api-gateway",
			"impact": "Outdated SSL policies increase exposure to known vulnerabilities",
			"resolution": "Use the most modern TLS/SSL policies available",
			"links": [
				"https://aquasecurity.github.io/tfsec/latest/checks/aws/api-gateway/use-secure-tls-policy/",
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_domain_name#security_policy"
			],
			"description": "Domain name is configured with an outdated TLS policy.",
			"severity": "HIGH",
			"warning": false,
			"status": 0,
			"resource": "aws_api_gateway_domain_name.outdated_security_policy",
			"location": {
				"filename": "/home/foobar/go/pkg/mod/github.com/aquasecurity/tfsec@v1.28.1/_examples/main.tf",
				"start_line": 41,
				"end_line": 41
			}
		},
		{
			"rule_id": "AVD-AWS-0005",
			"long_id": "aws-api-gateway-use-secure-tls-policy",
			"rule_description": "API Gateway domain name uses outdated SSL/TLS protocols.",
			"rule_provider": "aws",
			"rule_service": "api-gateway",
			"impact": "Outdated SSL policies increase exposure to known vulnerabilities",
			"resolution": "Use the most modern TLS/SSL policies available",
			"links": [
				"https://aquasecurity.github.io/tfsec/latest/checks/aws/api-gateway/use-secure-tls-policy/",
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_domain_name#security_policy"
			],
			"description": "Domain name is configured with an outdated TLS policy.",
			"severity": "HIGH",
			"warning": false,
			"status": 0,
			"resource": "aws_api_gateway_domain_name.empty_security_policy",
			"location": {
				"filename": "/home/foobar/go/pkg/mod/github.com/aquasecurity/tfsec@v1.28.1/_examples/main.tf",
				"start_line": 37,
				"end_line": 37
			}
		}]}`
