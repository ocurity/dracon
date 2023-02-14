package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components/producers/kics/types"
)

func TestParseOut(t *testing.T) {
	var results types.KICSOut
	err := json.Unmarshal([]byte(exampleOutput), &results)
	if err != nil {
		t.Logf(err.Error())
		t.Fail()
	}
	issues := parseOut(results)
	expectedIssues := []*v1.Issue{
		{
			Target:      "../../code/assets/queries/ansible/azure/sql_server_predictable_admin_account_name/test/positive.yaml:10",
			Type:        "MissingAttribute",
			Title:       "Insecure Configurations azure_rm_sqlserver Create (or update) SQL Server2",
			Severity:    4,
			Description: "{\"query_name\":\"AD Admin Not Configured For SQL Server\",\"query_id\":\"b176e927-bbe2-44a6-a9c3-041417137e5f\",\"query_url\":\"https://docs.ansible.com/ansible/latest/collections/azure/azcollection/azure_rm_sqlserver_module.html#parameter-ad_user\",\"severity\":\"HIGH\",\"platform\":\"Ansible\",\"cloud_provider\":\"AZURE\",\"category\":\"Insecure Configurations\",\"description\":\"The Active Directory Administrator is not configured for a SQL server\",\"description_id\":\"afa96f09\",\"files\":[{\"file_name\":\"../../code/assets/queries/ansible/azure/sql_server_predictable_admin_account_name/test/positive.yaml\",\"similarity_id\":\"819786035c5ea3943e303532e747fabac9f8beaddddfd6ecb863382481c50c0c\",\"line\":10,\"resource_type\":\"azure_rm_sqlserver\",\"resource_name\":\"Create (or update) SQL Server2\",\"issue_type\":\"MissingAttribute\",\"search_key\":\"name={{Create (or update) SQL Server2}}.{{azure_rm_sqlserver}}\",\"search_line\":0,\"search_value\":\"\",\"expected_value\":\"azure_rm_sqlserver.ad_user should be defined\",\"actual_value\":\"azure_rm_sqlserver.ad_user is undefined\"}]}",
		},
		{
			Target:      "../../code/assets/queries/ansible/azure/ad_admin_not_configured_for_sql_server/test/positive.yaml:3",
			Type:        "MissingAttribute",
			Title:       "Insecure Configurations azure_rm_sqlserver Create (or update) SQL Server",
			Severity:    4,
			Description: "{\"query_name\":\"AD Admin Not Configured For SQL Server\",\"query_id\":\"b176e927-bbe2-44a6-a9c3-041417137e5f\",\"query_url\":\"https://docs.ansible.com/ansible/latest/collections/azure/azcollection/azure_rm_sqlserver_module.html#parameter-ad_user\",\"severity\":\"HIGH\",\"platform\":\"Ansible\",\"cloud_provider\":\"AZURE\",\"category\":\"Insecure Configurations\",\"description\":\"The Active Directory Administrator is not configured for a SQL server\",\"description_id\":\"afa96f09\",\"files\":[{\"file_name\":\"../../code/assets/queries/ansible/azure/ad_admin_not_configured_for_sql_server/test/positive.yaml\",\"similarity_id\":\"e2acb4c0c4afe11e1908953bb768bab35e0ca53d4d20bc6af02f1c5350cbe587\",\"line\":3,\"resource_type\":\"azure_rm_sqlserver\",\"resource_name\":\"Create (or update) SQL Server\",\"issue_type\":\"MissingAttribute\",\"search_key\":\"name={{Create (or update) SQL Server}}.{{azure_rm_sqlserver}}\",\"search_line\":0,\"search_value\":\"\",\"expected_value\":\"azure_rm_sqlserver.ad_user should be defined\",\"actual_value\":\"azure_rm_sqlserver.ad_user is undefined\"}]}",
		},
		{
			Target:      "../../code/assets/queries/terraform/azure/small_msql_server_audit_retention/test/positive.tf:54",
			Type:        "MissingAttribute",
			Title:       "Insecure Configurations azurerm_sql_server sqlserver",
			Severity:    4,
			Description: "{\"query_name\":\"AD Admin Not Configured For SQL Server\",\"query_id\":\"a3a055d2-9a2e-4cc9-b9fb-12850a1a3a4b\",\"query_url\":\"https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/sql_active_directory_administrator\",\"severity\":\"HIGH\",\"platform\":\"Terraform\",\"cloud_provider\":\"AZURE\",\"category\":\"Insecure Configurations\",\"description\":\"The Active Directory Administrator is not configured for a SQL server\",\"description_id\":\"bccbda19\",\"files\":[{\"file_name\":\"../../code/assets/queries/terraform/azure/small_msql_server_audit_retention/test/positive.tf\",\"similarity_id\":\"a34bdc92c802268f18d3d56d43da1744d74d2dcdb8de3ca74ea2a92471eb7f5e\",\"line\":54,\"resource_type\":\"azurerm_sql_server\",\"resource_name\":\"sqlserver\",\"issue_type\":\"MissingAttribute\",\"search_key\":\"azurerm_sql_server[positive4]\",\"search_line\":0,\"search_value\":\"\",\"expected_value\":\"A 'azurerm_sql_active_directory_administrator' should be defined for 'azurerm_sql_server[positive4]'\",\"actual_value\":\"A 'azurerm_sql_active_directory_administrator' is not defined for 'azurerm_sql_server[positive4]'\"}]}",
		},
		{
			Target:      "../../code/assets/queries/terraform/azure/sql_server_predictable_admin_account_name/test/positive.tf:35",
			Type:        "MissingAttribute",
			Title:       "Insecure Configurations azurerm_sql_server mssqlserver",
			Severity:    4,
			Description: "{\"query_name\":\"AD Admin Not Configured For SQL Server\",\"query_id\":\"a3a055d2-9a2e-4cc9-b9fb-12850a1a3a4b\",\"query_url\":\"https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/sql_active_directory_administrator\",\"severity\":\"HIGH\",\"platform\":\"Terraform\",\"cloud_provider\":\"AZURE\",\"category\":\"Insecure Configurations\",\"description\":\"The Active Directory Administrator is not configured for a SQL server\",\"description_id\":\"bccbda19\",\"files\":[{\"file_name\":\"../../code/assets/queries/terraform/azure/sql_server_predictable_admin_account_name/test/positive.tf\",\"similarity_id\":\"86e39f354b76020e5e4d3fa85469a87509ffff482eaada3567df0f873b1642de\",\"line\":35,\"resource_type\":\"azurerm_sql_server\",\"resource_name\":\"mssqlserver\",\"issue_type\":\"MissingAttribute\",\"search_key\":\"azurerm_sql_server[positive4]\",\"search_line\":0,\"search_value\":\"\",\"expected_value\":\"A 'azurerm_sql_active_directory_administrator' should be defined for 'azurerm_sql_server[positive4]'\",\"actual_value\":\"A 'azurerm_sql_active_directory_administrator' is not defined for 'azurerm_sql_server[positive4]'\"}]}",
		},
	}

	found := 0
	assert.Equal(t, len(expectedIssues), len(issues))
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
	"kics_version": "v1.6.5",
	"files_scanned": 6871,
	"lines_scanned": 238272,
	"files_parsed": 6858,
	"lines_parsed": 238633,
	"files_failed_to_scan": 0,
	"queries_total": 2441,
	"queries_failed_to_execute": 0,
	"queries_failed_to_compute_similarity_id": 0,
	"scan_id": "console",
	"severity_counters": {
		"HIGH": 15397,
		"INFO": 6131,
		"LOW": 10852,
		"MEDIUM": 36628,
		"TRACE": 0
	},
	"total_counter": 69008,
	"total_bom_resources": 0,
	"start": "2022-11-30T16:37:35.819129687Z",
	"end": "2022-11-30T16:43:36.55869707Z",
	"paths": [
		"/code"
	],
	"queries": [
		{
			"query_name": "AD Admin Not Configured For SQL Server",
			"query_id": "b176e927-bbe2-44a6-a9c3-041417137e5f",
			"query_url": "https://docs.ansible.com/ansible/latest/collections/azure/azcollection/azure_rm_sqlserver_module.html#parameter-ad_user",
			"severity": "HIGH",
			"platform": "Ansible",
			"cloud_provider": "AZURE",
			"category": "Insecure Configurations",
			"description": "The Active Directory Administrator is not configured for a SQL server",
			"description_id": "afa96f09",
			"files": [
				{
					"file_name": "../../code/assets/queries/ansible/azure/sql_server_predictable_admin_account_name/test/positive.yaml",
					"similarity_id": "819786035c5ea3943e303532e747fabac9f8beaddddfd6ecb863382481c50c0c",
					"line": 10,
					"resource_type": "azure_rm_sqlserver",
					"resource_name": "Create (or update) SQL Server2",
					"issue_type": "MissingAttribute",
					"search_key": "name={{Create (or update) SQL Server2}}.{{azure_rm_sqlserver}}",
					"search_line": 0,
					"search_value": "",
					"expected_value": "azure_rm_sqlserver.ad_user should be defined",
					"actual_value": "azure_rm_sqlserver.ad_user is undefined"
				},
				{
					"file_name": "../../code/assets/queries/ansible/azure/ad_admin_not_configured_for_sql_server/test/positive.yaml",
					"similarity_id": "e2acb4c0c4afe11e1908953bb768bab35e0ca53d4d20bc6af02f1c5350cbe587",
					"line": 3,
					"resource_type": "azure_rm_sqlserver",
					"resource_name": "Create (or update) SQL Server",
					"issue_type": "MissingAttribute",
					"search_key": "name={{Create (or update) SQL Server}}.{{azure_rm_sqlserver}}",
					"search_line": 0,
					"search_value": "",
					"expected_value": "azure_rm_sqlserver.ad_user should be defined",
					"actual_value": "azure_rm_sqlserver.ad_user is undefined"
				}
			]
		},
		{
			"query_name": "AD Admin Not Configured For SQL Server",
			"query_id": "a3a055d2-9a2e-4cc9-b9fb-12850a1a3a4b",
			"query_url": "https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/sql_active_directory_administrator",
			"severity": "HIGH",
			"platform": "Terraform",
			"cloud_provider": "AZURE",
			"category": "Insecure Configurations",
			"description": "The Active Directory Administrator is not configured for a SQL server",
			"description_id": "bccbda19",
			"files": [
				{
					"file_name": "../../code/assets/queries/terraform/azure/small_msql_server_audit_retention/test/positive.tf",
					"similarity_id": "a34bdc92c802268f18d3d56d43da1744d74d2dcdb8de3ca74ea2a92471eb7f5e",
					"line": 54,
					"resource_type": "azurerm_sql_server",
					"resource_name": "sqlserver",
					"issue_type": "MissingAttribute",
					"search_key": "azurerm_sql_server[positive4]",
					"search_line": 0,
					"search_value": "",
					"expected_value": "A 'azurerm_sql_active_directory_administrator' should be defined for 'azurerm_sql_server[positive4]'",
					"actual_value": "A 'azurerm_sql_active_directory_administrator' is not defined for 'azurerm_sql_server[positive4]'"
				},
				{
					"file_name": "../../code/assets/queries/terraform/azure/sql_server_predictable_admin_account_name/test/positive.tf",
					"similarity_id": "86e39f354b76020e5e4d3fa85469a87509ffff482eaada3567df0f873b1642de",
					"line": 35,
					"resource_type": "azurerm_sql_server",
					"resource_name": "mssqlserver",
					"issue_type": "MissingAttribute",
					"search_key": "azurerm_sql_server[positive4]",
					"search_line": 0,
					"search_value": "",
					"expected_value": "A 'azurerm_sql_active_directory_administrator' should be defined for 'azurerm_sql_server[positive4]'",
					"actual_value": "A 'azurerm_sql_active_directory_administrator' is not defined for 'azurerm_sql_server[positive4]'"
				}
			]}]}`
