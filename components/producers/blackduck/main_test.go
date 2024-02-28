package main

import (
	"encoding/json"
	"testing"

	v1 "github.com/ocurity/dracon/api/proto/v1"

	"github.com/stretchr/testify/assert"
)

func TestParseIssues(t *testing.T) {
	var results BlackduckOut
	err := json.Unmarshal([]byte(BDOut), &results)
	assert.Nil(t, err)

	issues, err := parseIssues(&results)
	assert.Nil(t, err)
	// cwe0 := []int32{400}
	// cwe1 := []int32{601, 918}
	expectedIssue := []*v1.Issue{
		{
			Target:      "pkg:maven/org.apache.tomcat/tomcat-annotations-api@9.0.81",
			Type:        "CVE-2023-46589",
			Title:       "Improper Input Validation vulnerability in Apache Tomcat.Tomcat from 11.0.0-M1 through 11.0.0-M10, from 10.1.0-M1 through 10.1.15, from 9.0.0-M1 through 9.0.82 and from 8.5.0 through 8.5.95 did not correctly parse HTTP trailer headers. A trailer header that exceeded the header size limit could cause Tomcat to treat a single \nrequest as multiple requests leading to the possibility of request \nsmuggling when behind a reverse proxy.\n\nUsers are recommended to upgrade to version 11.0.0-M11\u00a0onwards, 10.1.16 onwards, 9.0.83 onwards or 8.5.96 onwards, which fix the issue.\n\n",
			Severity:    4,
			Cvss:        7.5,
			Description: "Improper Input Validation vulnerability in Apache Tomcat.Tomcat from 11.0.0-M1 through 11.0.0-M10, from 10.1.0-M1 through 10.1.15, from 9.0.0-M1 through 9.0.82 and from 8.5.0 through 8.5.95 did not correctly parse HTTP trailer headers. A trailer header that exceeded the header size limit could cause Tomcat to treat a single \nrequest as multiple requests leading to the possibility of request \nsmuggling when behind a reverse proxy.\n\nUsers are recommended to upgrade to version 11.0.0-M11\u00a0onwards, 10.1.16 onwards, 9.0.83 onwards or 8.5.96 onwards, which fix the issue.\n\n\nSolution Available: false\nWorkaround Available: false\nExploit Available: false\nOriginal Description: Improper Input Validation vulnerability in Apache Tomcat.Tomcat from 11.0.0-M1 through 11.0.0-M10, from 10.1.0-M1 through 10.1.15, from 9.0.0-M1 through 9.0.82 and from 8.5.0 through 8.5.95 did not correctly parse HTTP trailer headers. A trailer header that exceeded the header size limit could cause Tomcat to treat a single \nrequest as multiple requests leading to the possibility of request \nsmuggling when behind a reverse proxy.\n\nUsers are recommended to upgrade to version 11.0.0-M11\u00a0onwards, 10.1.16 onwards, 9.0.83 onwards or 8.5.96 onwards, which fix the issue.\n\n",
			Cve:         "CVE-2023-46589",
			Cwe:         []int32{444},
		},
		{
			Target:      "pkg:maven/org.apache.kafka/kafka-clients@3.1.2",
			Type:        "BDSA-2023-0235",
			Title:       "Apache Kafka is vulnerable to remote code execution (RCE) due to insecure handling of untrusted serialized data sent using the Kafka REST API. This could allow an attacker to supply crafted data to execute Java deserialization gadget chains on the Kafka connect server.",
			Severity:    4,
			Cvss:        7.9,
			Description: "Apache Kafka is vulnerable to remote code execution (RCE) due to insecure handling of untrusted serialized data sent using the Kafka REST API. This could allow an attacker to supply crafted data to execute Java deserialization gadget chains on the Kafka connect server.\nSolution Available: true\nWorkaround Available: false\nExploit Available: true\nOriginal Description: Apache Kafka is vulnerable to remote code execution (RCE) due to insecure handling of untrusted serialized data sent using the Kafka REST API. This could allow an attacker to supply crafted data to execute Java deserialization gadget chains on the Kafka connect server.",
			Cve:         "CVE-2023-25194",
			Cwe:         []int32{502},
		},
	}
	assert.Equal(t, expectedIssue, issues)
}

var BDOut = `{
    "hubUILink": "https://foo.opensourcehub31.foo..net/api/projects/eec0e9a7-d2e5-41e7-97b2-f91c46df9685/versions/0e7a2872-230f-4eee-a9f0-b18e1c50c2d3/components",
    "hubAPILink": "https://foo.opensourcehub31.foo..net/api/projects/eec0e9a7-d2e5-41e7-97b2-f91c46df9685/versions/0e7a2872-230f-4eee-a9f0-b18e1c50c2d3/components",
    "versionId": "0e7a2872-230f-4eee-a9f0-b18e1c50c2d3",
    "totalComponentsFound": 30,
    "matchedFlag": true,
    "headers": [
        "ComponentName",
        "ComponentVersion",
        "Critical",
        "High",
        "Medium",
        "Low",
        "Operational Risk",
        "License Risk",
        "License"
    ],
    "summary": {
        "securityRisksSeverityCritical": 1,
        "securityRisksSeverityHigh": 9,
        "securityRisksSeverityMedium": 1,
        "securityRisksSeverityLow": 0,
        "licenseRisksSeverityHigh": 0,
        "licenseRisksSeverityMedium": 4,
        "licenseRisksSeverityLow": 0,
        "operationRisksSeverityHigh": 4,
        "operationRisksSeverityMedium": 3,
        "operationRisksSeverityLow": 22,
        "quarantineEligibleFlag": true,
        "scanStartTime": "2024-01-18T17:07:01.596Z",
        "scanCompleteTime": "2024-01-18T17:07:13.284Z"
    },
    "scanId": "ed8590ff-1d07-4a91-8453-d0a8abd792a3",
    "appid": "171845",
    "appname": "dtliveupdateconsumer",
    "releaseid": "LatestProductionDeployedScan",
    "projectId": "eec0e9a7-d2e5-41e7-97b2-f91c46df9685",
    "data": [
		{
            "componentName": "Apache Tomcat",
            "componentVersion": "9.0.81",
            "releasedOn": "2023-10-10T01:10:06.000Z",
            "critical": 0,
            "high": 1,
            "medium": 0,
            "low": 0,
            "operationalRisk": "LOW",
            "licenseRisk": "",
            "license": "Apache License 2.0",
            "ownership": "OPEN_SOURCE",
            "kburl": "https://foo.opensourcehub31.blah.net/api/components/5a7e1c49-9a98-4393-b4e0-8011122bbe2f/versions/fcaa4dea-eaed-48f0-85af-feb369a58945",
            "versionid": "fcaa4dea-eaed-48f0-85af-feb369a58945",
            "ctcid": "",
            "checksum": "",
            "componentid": "5a7e1c49-9a98-4393-b4e0-8011122bbe2f",
            "matchedFiles": {
                "totalCount": 1,
                "items": [
                    {
                        "filePath": {
                            "path": "/javax/",
                            "archiveContext": "extractedContent/dtwebhook/webhook/dtliveupdateconsumer.jar!/BOOT-INF/lib/jakarta.annotation-api-1.3.5.jar!/",
                            "compositePathContext": "/javax/#extractedContent/dtwebhook/webhook/dtliveupdateconsumer.jar!/BOOT-INF/lib/jakarta.annotation-api-1.3.5.jar!/",
                            "fileName": "javax"
                        },
                        "usages": [
                            "DYNAMICALLY_LINKED"
                        ],
                        "_meta": {
                            "links": [
                                {
                                    "rel": "codelocations",
                                    "href": "https://foo.opensourcehub31.blah.net/api/codelocations/c10e650e-5257-44f2-a0e7-ad5694c5dcf8"
                                }
                            ]
                        }
                    }
                ],
                "appliedFilters": [],
                "_meta": {}
            },
            "vulnerabilities": {
                "totalCount": 1,
                "items": [
                    {
                        "id": "CVE-2023-46589",
                        "summary": "Improper Input Validation vulnerability in Apache Tomcat.Tomcat from 11.0.0-M1 through 11.0.0-M10, from 10.1.0-M1 through 10.1.15, from 9.0.0-M1 through 9.0.82 and from 8.5.0 through 8.5.95 did not correctly parse HTTP trailer headers. A trailer header that exceeded the header size limit could cause Tomcat to treat a single \nrequest as multiple requests leading to the possibility of request \nsmuggling when behind a reverse proxy.\n\nUsers are recommended to upgrade to version 11.0.0-M11\u00a0onwards, 10.1.16 onwards, 9.0.83 onwards or 8.5.96 onwards, which fix the issue.\n\n",
                        "publishedDate": "2023-11-28T16:15:06.943Z",
                        "lastModified": "2024-01-05T11:15:09.847Z",
                        "source": "NVD",
                        "remediationStatus": "NEW",
                        "createdAt": "2024-01-18T17:07:20.083Z",
                        "updatedAt": "2024-01-18T17:07:20.083Z",
                        "createdBy": {
                            "userName": "eb_bd_fid",
                            "firstName": "eb_bd_fid",
                            "lastName": "eb_bd_fid",
                            "user": "https://foo.opensourcehub31.blah.net/api/users/c5c8c47c-32c6-424e-b199-f9c154bf26fa"
                        },
                        "updatedBy": {
                            "userName": "eb_bd_fid",
                            "firstName": "eb_bd_fid",
                            "lastName": "eb_bd_fid",
                            "user": "https://foo.opensourcehub31.blah.net/api/users/c5c8c47c-32c6-424e-b199-f9c154bf26fa"
                        },
                        "cweIds": [
                            "CWE-444"
                        ],
                        "cvss2": {
                            "baseScore": 0.0,
                            "impactSubscore": 0.0,
                            "exploitabilitySubscore": 1.2,
                            "accessVector": "LOCAL",
                            "accessComplexity": "HIGH",
                            "authentication": "MULTIPLE",
                            "confidentialityImpact": "NONE",
                            "availabilityImpact": "NONE",
                            "source": "NVD",
                            "severity": "LOW",
                            "integrityImpact": "NONE",
                            "vector": "(AV:L/AC:H/Au:M/C:N/I:N/A:N)",
                            "overallScore": 0.0
                        },
                        "cvss3": {
                            "baseScore": 7.5,
                            "impactSubscore": 3.6,
                            "exploitabilitySubscore": 3.9,
                            "attackVector": "NETWORK",
                            "attackComplexity": "LOW",
                            "confidentialityImpact": "NONE",
                            "integrityImpact": "HIGH",
                            "availabilityImpact": "NONE",
                            "privilegesRequired": "NONE",
                            "scope": "UNCHANGED",
                            "userInteraction": "NONE",
                            "source": "NVD",
                            "severity": "HIGH",
                            "vector": "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:H/A:N",
                            "overallScore": 7.5
                        },
                        "useCvss3": true,
                        "overallScore": 7.5,
                        "solutionAvailable": false,
                        "workaroundAvailable": false,
                        "exploitAvailable": false,
                        "bdsaTags": [],
                        "_meta": {
                            "allow": [
                                "GET"
                            ],
                            "href": "https://foo.opensourcehub31.blah.net/api/projects/eec0e9a7-d2e5-41e7-97b2-f91c46df9685/versions/0e7a2872-230f-4eee-a9f0-b18e1c50c2d3/components/5a7e1c49-9a98-4393-b4e0-8011122bbe2f/versions/fcaa4dea-eaed-48f0-85af-feb369a58945/origins/82f84e5d-1ac9-4812-9159-701a9fb07c3d/vulnerabilities",
                            "links": [
                                {
                                    "rel": "vulnerability",
                                    "href": "https://foo.opensourcehub31.blah.net/api/vulnerabilities/CVE-2023-46589"
                                },
                                {
                                    "rel": "related-vulnerability",
                                    "href": "https://foo.opensourcehub31.blah.net/api/vulnerabilities/BDSA-2023-3298",
                                    "label": "BDSA-2023-3298"
                                },
                                {
                                    "rel": "vulnerability-with-remediation",
                                    "href": "https://foo.opensourcehub31.blah.net/api/projects/eec0e9a7-d2e5-41e7-97b2-f91c46df9685/versions/0e7a2872-230f-4eee-a9f0-b18e1c50c2d3/components/5a7e1c49-9a98-4393-b4e0-8011122bbe2f/versions/fcaa4dea-eaed-48f0-85af-feb369a58945/origins/82f84e5d-1ac9-4812-9159-701a9fb07c3d/vulnerabilities/CVE-2023-46589/remediation"
                                }
                            ]
                        },
                        "name": "CVE-2023-46589",
                        "description": "Improper Input Validation vulnerability in Apache Tomcat.Tomcat from 11.0.0-M1 through 11.0.0-M10, from 10.1.0-M1 through 10.1.15, from 9.0.0-M1 through 9.0.82 and from 8.5.0 through 8.5.95 did not correctly parse HTTP trailer headers. A trailer header that exceeded the header size limit could cause Tomcat to treat a single \nrequest as multiple requests leading to the possibility of request \nsmuggling when behind a reverse proxy.\n\nUsers are recommended to upgrade to version 11.0.0-M11\u00a0onwards, 10.1.16 onwards, 9.0.83 onwards or 8.5.96 onwards, which fix the issue.\n\n"
                    }
                ],
                "appliedFilters": [],
                "_meta": {
                    "allow": [
                        "GET"
                    ],
                    "href": "https://foo.opensourcehub31.blah.net/api/projects/eec0e9a7-d2e5-41e7-97b2-f91c46df9685/versions/0e7a2872-230f-4eee-a9f0-b18e1c50c2d3/components/5a7e1c49-9a98-4393-b4e0-8011122bbe2f/versions/fcaa4dea-eaed-48f0-85af-feb369a58945/origins/82f84e5d-1ac9-4812-9159-701a9fb07c3d/vulnerabilities",
                    "links": [
                        {
                            "rel": "static-filter",
                            "href": "https://foo.opensourcehub31.blah.net/api/projects/eec0e9a7-d2e5-41e7-97b2-f91c46df9685/versions/0e7a2872-230f-4eee-a9f0-b18e1c50c2d3/components/5a7e1c49-9a98-4393-b4e0-8011122bbe2f/versions/fcaa4dea-eaed-48f0-85af-feb369a58945/origins/82f84e5d-1ac9-4812-9159-701a9fb07c3d/vulnerability-filters?filterKey=exploitAvailable",
                            "name": "exploitAvailable",
                            "label": "Exploit"
                        },
                        {
                            "rel": "range-filter",
                            "href": "https://foo.opensourcehub31.blah.net/api/projects/eec0e9a7-d2e5-41e7-97b2-f91c46df9685/versions/0e7a2872-230f-4eee-a9f0-b18e1c50c2d3/components/5a7e1c49-9a98-4393-b4e0-8011122bbe2f/versions/fcaa4dea-eaed-48f0-85af-feb369a58945/origins/82f84e5d-1ac9-4812-9159-701a9fb07c3d/vulnerability-filters?filterKey=securityRiskScore",
                            "name": "securityRiskScore",
                            "label": "Overall Score"
                        },
                        {
                            "rel": "quick-filter",
                            "href": "https://foo.opensourcehub31.blah.net/api/projects/eec0e9a7-d2e5-41e7-97b2-f91c46df9685/versions/0e7a2872-230f-4eee-a9f0-b18e1c50c2d3/components/5a7e1c49-9a98-4393-b4e0-8011122bbe2f/versions/fcaa4dea-eaed-48f0-85af-feb369a58945/origins/82f84e5d-1ac9-4812-9159-701a9fb07c3d/vulnerability-filters?filterKey=reachable",
                            "name": "reachable",
                            "label": "Reachable"
                        },
                        {
                            "rel": "static-filter",
                            "href": "https://foo.opensourcehub31.blah.net/api/projects/eec0e9a7-d2e5-41e7-97b2-f91c46df9685/versions/0e7a2872-230f-4eee-a9f0-b18e1c50c2d3/components/5a7e1c49-9a98-4393-b4e0-8011122bbe2f/versions/fcaa4dea-eaed-48f0-85af-feb369a58945/origins/82f84e5d-1ac9-4812-9159-701a9fb07c3d/vulnerability-filters?filterKey=remediationType",
                            "name": "remediationType",
                            "label": "Remediation Status"
                        },
                        {
                            "rel": "static-filter",
                            "href": "https://foo.opensourcehub31.blah.net/api/projects/eec0e9a7-d2e5-41e7-97b2-f91c46df9685/versions/0e7a2872-230f-4eee-a9f0-b18e1c50c2d3/components/5a7e1c49-9a98-4393-b4e0-8011122bbe2f/versions/fcaa4dea-eaed-48f0-85af-feb369a58945/origins/82f84e5d-1ac9-4812-9159-701a9fb07c3d/vulnerability-filters?filterKey=solutionAvailable",
                            "name": "solutionAvailable",
                            "label": "Solution"
                        },
                        {
                            "rel": "static-filter",
                            "href": "https://foo.opensourcehub31.blah.net/api/projects/eec0e9a7-d2e5-41e7-97b2-f91c46df9685/versions/0e7a2872-230f-4eee-a9f0-b18e1c50c2d3/components/5a7e1c49-9a98-4393-b4e0-8011122bbe2f/versions/fcaa4dea-eaed-48f0-85af-feb369a58945/origins/82f84e5d-1ac9-4812-9159-701a9fb07c3d/vulnerability-filters?filterKey=workaroundAvailable",
                            "name": "workaroundAvailable",
                            "label": "Workaround"
                        }
                    ]
                }
            },
            "policyviolations": {},
            "origins": {
                "name": "9.0.81",
                "origin": "https://foo.opensourcehub31.blah.net/api/components/5a7e1c49-9a98-4393-b4e0-8011122bbe2f/versions/fcaa4dea-eaed-48f0-85af-feb369a58945/origins/82f84e5d-1ac9-4812-9159-701a9fb07c3d",
                "externalNamespace": "maven",
                "externalId": "org.apache.tomcat:tomcat-annotations-api:9.0.81",
                "externalNamespaceDistribution": false,
                "_meta": {
                    "allow": [],
                    "links": [
                        {
                            "rel": "origin",
                            "href": "https://foo.opensourcehub31.blah.net/api/components/5a7e1c49-9a98-4393-b4e0-8011122bbe2f/versions/fcaa4dea-eaed-48f0-85af-feb369a58945/origins/82f84e5d-1ac9-4812-9159-701a9fb07c3d"
                        },
                        {
                            "rel": "matched-files",
                            "href": "https://foo.opensourcehub31.blah.net/api/projects/eec0e9a7-d2e5-41e7-97b2-f91c46df9685/versions/0e7a2872-230f-4eee-a9f0-b18e1c50c2d3/components/5a7e1c49-9a98-4393-b4e0-8011122bbe2f/versions/fcaa4dea-eaed-48f0-85af-feb369a58945/origins/82f84e5d-1ac9-4812-9159-701a9fb07c3d/matched-files"
                        },
                        {
                            "rel": "upgrade-guidance",
                            "href": "https://foo.opensourcehub31.blah.net/api/components/5a7e1c49-9a98-4393-b4e0-8011122bbe2f/versions/fcaa4dea-eaed-48f0-85af-feb369a58945/origins/82f84e5d-1ac9-4812-9159-701a9fb07c3d/upgrade-guidance"
                        },
                        {
                            "rel": "component-origin-copyrights",
                            "href": "https://foo.opensourcehub31.blah.net/api/components/5a7e1c49-9a98-4393-b4e0-8011122bbe2f/versions/fcaa4dea-eaed-48f0-85af-feb369a58945/origins/82f84e5d-1ac9-4812-9159-701a9fb07c3d/copyrights"
                        }
                    ]
                },
                "originName": "maven",
                "originId": "org.apache.tomcat:tomcat-annotations-api:9.0.81"
            }
        },{
            "componentName": "Apache Kafka",
            "componentVersion": "3.1.2",
            "releasedOn": "2022-09-09T20:08:36.000Z",
            "critical": 0,
            "high": 1,
            "medium": 0,
            "low": 0,
            "operationalRisk": "LOW",
            "licenseRisk": "",
            "license": "Apache License 2.0",
            "ownership": "OPEN_SOURCE",
            "kburl": "https://foo.opensourcehub31.foo..net/api/components/05248305-719c-4b7d-a693-0b1a7992b4ec/versions/b54827b4-ee6f-4204-9f13-f3504e064efa",
            "versionid": "b54827b4-ee6f-4204-9f13-f3504e064efa",
            "ctcid": "",
            "checksum": "",
            "componentid": "05248305-719c-4b7d-a693-0b1a7992b4ec",
            "matchedFiles": {
                "totalCount": 1,
                "items": [
                    {
                        "filePath": {
                            "path": "/BOOT-INF/lib/kafka-clients-3.1.2.jar",
                            "archiveContext": "extractedContent/dtwebhook/webhook/dtliveupdateconsumer.jar!/",
                            "compositePathContext": "/BOOT-INF/lib/kafka-clients-3.1.2.jar#extractedContent/dtwebhook/webhook/dtliveupdateconsumer.jar!/",
                            "fileName": "kafka-clients-3.1.2.jar"
                        },
                        "usages": [
                            "DYNAMICALLY_LINKED"
                        ],
                        "_meta": {
                            "links": [
                                {
                                    "rel": "codelocations",
                                    "href": "https://foo.opensourcehub31.foo..net/api/codelocations/c10e650e-5257-44f2-a0e7-ad5694c5dcf8"
                                }
                            ]
                        }
                    }
                ],
                "appliedFilters": [],
                "_meta": {}
            },
            "vulnerabilities": {
                "totalCount": 1,
                "items": [
                    {
                        "id": "BDSA-2023-0235",
                        "summary": "Apache Kafka is vulnerable to remote code execution (RCE) due to insecure handling of untrusted serialized data sent using the Kafka REST API. This could allow an attacker to supply crafted data to execute Java deserialization gadget chains on the Kafka connect server.",
                        "publishedDate": "2023-02-09T12:35:59.049Z",
                        "lastModified": "2023-02-09T12:35:59.037Z",
                        "source": "BDSA",
                        "remediationStatus": "NEW",
                        "createdAt": "2024-01-18T17:07:20.083Z",
                        "updatedAt": "2024-01-18T17:07:20.083Z",
                        "createdBy": {
                            "userName": "eb_bd_fid",
                            "firstName": "eb_bd_fid",
                            "lastName": "eb_bd_fid",
                            "user": "https://foo.opensourcehub31.foo..net/api/users/c5c8c47c-32c6-424e-b199-f9c154bf26fa"
                        },
                        "updatedBy": {
                            "userName": "eb_bd_fid",
                            "firstName": "eb_bd_fid",
                            "lastName": "eb_bd_fid",
                            "user": "https://foo.opensourcehub31.foo..net/api/users/c5c8c47c-32c6-424e-b199-f9c154bf26fa"
                        },
                        "cweIds": [
                            "CWE-502"
                        ],
                        "cvss2": {
                            "baseScore": 6.5,
                            "impactSubscore": 6.4,
                            "exploitabilitySubscore": 8.0,
                            "accessVector": "NETWORK",
                            "accessComplexity": "LOW",
                            "authentication": "SINGLE",
                            "confidentialityImpact": "PARTIAL",
                            "availabilityImpact": "PARTIAL",
                            "temporalMetrics": {
                                "exploitability": "PROOF_OF_CONCEPT",
                                "remediationLevel": "OFFICIAL_FIX",
                                "reportConfidence": "CONFIRMED",
                                "score": 5.1
                            },
                            "source": "BDSA",
                            "severity": "MEDIUM",
                            "integrityImpact": "PARTIAL",
                            "vector": "(AV:N/AC:L/Au:S/C:P/I:P/A:P/E:POC/RL:OF/RC:C)",
                            "overallScore": 5.1
                        },
                        "cvss3": {
                            "baseScore": 8.8,
                            "impactSubscore": 5.9,
                            "exploitabilitySubscore": 2.8,
                            "attackVector": "NETWORK",
                            "attackComplexity": "LOW",
                            "confidentialityImpact": "HIGH",
                            "integrityImpact": "HIGH",
                            "availabilityImpact": "HIGH",
                            "privilegesRequired": "LOW",
                            "scope": "UNCHANGED",
                            "userInteraction": "NONE",
                            "source": "BDSA",
                            "severity": "HIGH",
                            "temporalMetrics": {
                                "exploitCodeMaturity": "PROOF_OF_CONCEPT",
                                "remediationLevel": "OFFICIAL_FIX",
                                "reportConfidence": "CONFIRMED",
                                "score": 7.9
                            },
                            "vector": "CVSS:3.1/AV:N/AC:L/PR:L/UI:N/S:U/C:H/I:H/A:H/E:P/RL:O/RC:C",
                            "overallScore": 7.9
                        },
                        "useCvss3": true,
                        "overallScore": 7.9,
                        "solutionAvailable": true,
                        "workaroundAvailable": false,
                        "exploitAvailable": true,
                        "bdsaTags": [
                            "RCE"
                        ],
                        "_meta": {
                            "allow": [
                                "GET"
                            ],
                            "href": "https://foo.opensourcehub31.foo..net/api/projects/eec0e9a7-d2e5-41e7-97b2-f91c46df9685/versions/0e7a2872-230f-4eee-a9f0-b18e1c50c2d3/components/05248305-719c-4b7d-a693-0b1a7992b4ec/versions/b54827b4-ee6f-4204-9f13-f3504e064efa/origins/abc9553b-a18b-4709-8914-fc6b8e7ba6c2/vulnerabilities",
                            "links": [
                                {
                                    "rel": "vulnerability",
                                    "href": "https://foo.opensourcehub31.foo..net/api/vulnerabilities/BDSA-2023-0235"
                                },
                                {
                                    "rel": "unmatched-related-vulnerability",
                                    "href": "https://foo.opensourcehub31.foo..net/api/vulnerabilities/CVE-2023-25194",
                                    "label": "CVE-2023-25194"
                                },
                                {
                                    "rel": "vulnerability-with-remediation",
                                    "href": "https://foo.opensourcehub31.foo..net/api/projects/eec0e9a7-d2e5-41e7-97b2-f91c46df9685/versions/0e7a2872-230f-4eee-a9f0-b18e1c50c2d3/components/05248305-719c-4b7d-a693-0b1a7992b4ec/versions/b54827b4-ee6f-4204-9f13-f3504e064efa/origins/abc9553b-a18b-4709-8914-fc6b8e7ba6c2/vulnerabilities/BDSA-2023-0235/remediation"
                                }
                            ]
                        },
                        "name": "BDSA-2023-0235",
                        "description": "Apache Kafka is vulnerable to remote code execution (RCE) due to insecure handling of untrusted serialized data sent using the Kafka REST API. This could allow an attacker to supply crafted data to execute Java deserialization gadget chains on the Kafka connect server."
                    }
                ],
                "appliedFilters": [],
                "_meta": {
                    "allow": [
                        "GET"
                    ],
                    "href": "https://foo.opensourcehub31.foo..net/api/projects/eec0e9a7-d2e5-41e7-97b2-f91c46df9685/versions/0e7a2872-230f-4eee-a9f0-b18e1c50c2d3/components/05248305-719c-4b7d-a693-0b1a7992b4ec/versions/b54827b4-ee6f-4204-9f13-f3504e064efa/origins/abc9553b-a18b-4709-8914-fc6b8e7ba6c2/vulnerabilities",
                    "links": [
                        {
                            "rel": "static-filter",
                            "href": "https://foo.opensourcehub31.foo..net/api/projects/eec0e9a7-d2e5-41e7-97b2-f91c46df9685/versions/0e7a2872-230f-4eee-a9f0-b18e1c50c2d3/components/05248305-719c-4b7d-a693-0b1a7992b4ec/versions/b54827b4-ee6f-4204-9f13-f3504e064efa/origins/abc9553b-a18b-4709-8914-fc6b8e7ba6c2/vulnerability-filters?filterKey=exploitAvailable",
                            "name": "exploitAvailable",
                            "label": "Exploit"
                        },
                        {
                            "rel": "range-filter",
                            "href": "https://foo.opensourcehub31.foo..net/api/projects/eec0e9a7-d2e5-41e7-97b2-f91c46df9685/versions/0e7a2872-230f-4eee-a9f0-b18e1c50c2d3/components/05248305-719c-4b7d-a693-0b1a7992b4ec/versions/b54827b4-ee6f-4204-9f13-f3504e064efa/origins/abc9553b-a18b-4709-8914-fc6b8e7ba6c2/vulnerability-filters?filterKey=securityRiskScore",
                            "name": "securityRiskScore",
                            "label": "Overall Score"
                        },
                        {
                            "rel": "quick-filter",
                            "href": "https://foo.opensourcehub31.foo..net/api/projects/eec0e9a7-d2e5-41e7-97b2-f91c46df9685/versions/0e7a2872-230f-4eee-a9f0-b18e1c50c2d3/components/05248305-719c-4b7d-a693-0b1a7992b4ec/versions/b54827b4-ee6f-4204-9f13-f3504e064efa/origins/abc9553b-a18b-4709-8914-fc6b8e7ba6c2/vulnerability-filters?filterKey=reachable",
                            "name": "reachable",
                            "label": "Reachable"
                        },
                        {
                            "rel": "static-filter",
                            "href": "https://foo.opensourcehub31.foo..net/api/projects/eec0e9a7-d2e5-41e7-97b2-f91c46df9685/versions/0e7a2872-230f-4eee-a9f0-b18e1c50c2d3/components/05248305-719c-4b7d-a693-0b1a7992b4ec/versions/b54827b4-ee6f-4204-9f13-f3504e064efa/origins/abc9553b-a18b-4709-8914-fc6b8e7ba6c2/vulnerability-filters?filterKey=remediationType",
                            "name": "remediationType",
                            "label": "Remediation Status"
                        },
                        {
                            "rel": "static-filter",
                            "href": "https://foo.opensourcehub31.foo..net/api/projects/eec0e9a7-d2e5-41e7-97b2-f91c46df9685/versions/0e7a2872-230f-4eee-a9f0-b18e1c50c2d3/components/05248305-719c-4b7d-a693-0b1a7992b4ec/versions/b54827b4-ee6f-4204-9f13-f3504e064efa/origins/abc9553b-a18b-4709-8914-fc6b8e7ba6c2/vulnerability-filters?filterKey=solutionAvailable",
                            "name": "solutionAvailable",
                            "label": "Solution"
                        },
                        {
                            "rel": "static-filter",
                            "href": "https://foo.opensourcehub31.foo..net/api/projects/eec0e9a7-d2e5-41e7-97b2-f91c46df9685/versions/0e7a2872-230f-4eee-a9f0-b18e1c50c2d3/components/05248305-719c-4b7d-a693-0b1a7992b4ec/versions/b54827b4-ee6f-4204-9f13-f3504e064efa/origins/abc9553b-a18b-4709-8914-fc6b8e7ba6c2/vulnerability-filters?filterKey=workaroundAvailable",
                            "name": "workaroundAvailable",
                            "label": "Workaround"
                        }
                    ]
                }
            },
            "policyviolations": {},
            "origins": {
                "name": "3.1.2",
                "origin": "https://foo.opensourcehub31.foo..net/api/components/05248305-719c-4b7d-a693-0b1a7992b4ec/versions/b54827b4-ee6f-4204-9f13-f3504e064efa/origins/abc9553b-a18b-4709-8914-fc6b8e7ba6c2",
                "externalNamespace": "maven",
                "externalId": "org.apache.kafka:kafka-clients:3.1.2",
                "externalNamespaceDistribution": false,
                "_meta": {
                    "allow": [],
                    "links": [
                        {
                            "rel": "origin",
                            "href": "https://foo.opensourcehub31.foo..net/api/components/05248305-719c-4b7d-a693-0b1a7992b4ec/versions/b54827b4-ee6f-4204-9f13-f3504e064efa/origins/abc9553b-a18b-4709-8914-fc6b8e7ba6c2"
                        },
                        {
                            "rel": "matched-files",
                            "href": "https://foo.opensourcehub31.foo..net/api/projects/eec0e9a7-d2e5-41e7-97b2-f91c46df9685/versions/0e7a2872-230f-4eee-a9f0-b18e1c50c2d3/components/05248305-719c-4b7d-a693-0b1a7992b4ec/versions/b54827b4-ee6f-4204-9f13-f3504e064efa/origins/abc9553b-a18b-4709-8914-fc6b8e7ba6c2/matched-files"
                        },
                        {
                            "rel": "upgrade-guidance",
                            "href": "https://foo.opensourcehub31.foo..net/api/components/05248305-719c-4b7d-a693-0b1a7992b4ec/versions/b54827b4-ee6f-4204-9f13-f3504e064efa/origins/abc9553b-a18b-4709-8914-fc6b8e7ba6c2/upgrade-guidance"
                        },
                        {
                            "rel": "component-origin-copyrights",
                            "href": "https://foo.opensourcehub31.foo..net/api/components/05248305-719c-4b7d-a693-0b1a7992b4ec/versions/b54827b4-ee6f-4204-9f13-f3504e064efa/origins/abc9553b-a18b-4709-8914-fc6b8e7ba6c2/copyrights"
                        }
                    ]
                },
                "originName": "maven",
                "originId": "org.apache.kafka:kafka-clients:3.1.2"
            }
        },
        {
            "componentName": "Apache Log4j",
            "componentVersion": "2.17.2",
            "releasedOn": "2022-02-27T18:35:56.000Z",
            "critical": 0,
            "high": 0,
            "medium": 0,
            "low": 0,
            "operationalRisk": "LOW",
            "licenseRisk": "",
            "license": "Apache License 2.0",
            "ownership": "OPEN_SOURCE",
            "kburl": "https://foo.opensourcehub31.foo..net/api/components/7460c937-f013-4c3a-bdf3-ace04cfd0304/versions/86f3974b-d17a-4bc7-8592-61a618a0e140",
            "versionid": "86f3974b-d17a-4bc7-8592-61a618a0e140",
            "ctcid": "",
            "checksum": "",
            "componentid": "7460c937-f013-4c3a-bdf3-ace04cfd0304",
            "matchedFiles": {},
            "vulnerabilities": {},
            "policyviolations": {},
            "origins": {
                "name": "2.17.2",
                "origin": "https://foo.opensourcehub31.foo..net/api/components/7460c937-f013-4c3a-bdf3-ace04cfd0304/versions/86f3974b-d17a-4bc7-8592-61a618a0e140/origins/f6b76814-4d31-459b-9eb8-56802c29a16b",
                "externalNamespace": "maven",
                "externalId": "org.apache.logging.log4j:log4j-core:2.17.2",
                "externalNamespaceDistribution": false,
                "_meta": {
                    "allow": [],
                    "links": [
                        {
                            "rel": "origin",
                            "href": "https://foo.opensourcehub31.foo..net/api/components/7460c937-f013-4c3a-bdf3-ace04cfd0304/versions/86f3974b-d17a-4bc7-8592-61a618a0e140/origins/f6b76814-4d31-459b-9eb8-56802c29a16b"
                        },
                        {
                            "rel": "matched-files",
                            "href": "https://foo.opensourcehub31.foo..net/api/projects/eec0e9a7-d2e5-41e7-97b2-f91c46df9685/versions/0e7a2872-230f-4eee-a9f0-b18e1c50c2d3/components/7460c937-f013-4c3a-bdf3-ace04cfd0304/versions/86f3974b-d17a-4bc7-8592-61a618a0e140/origins/f6b76814-4d31-459b-9eb8-56802c29a16b/matched-files"
                        },
                        {
                            "rel": "upgrade-guidance",
                            "href": "https://foo.opensourcehub31.foo..net/api/components/7460c937-f013-4c3a-bdf3-ace04cfd0304/versions/86f3974b-d17a-4bc7-8592-61a618a0e140/origins/f6b76814-4d31-459b-9eb8-56802c29a16b/upgrade-guidance"
                        },
                        {
                            "rel": "component-origin-copyrights",
                            "href": "https://foo.opensourcehub31.foo..net/api/components/7460c937-f013-4c3a-bdf3-ace04cfd0304/versions/86f3974b-d17a-4bc7-8592-61a618a0e140/origins/f6b76814-4d31-459b-9eb8-56802c29a16b/copyrights"
                        }
                    ]
                },
                "originName": "maven",
                "originId": "org.apache.logging.log4j:log4j-core:2.17.2"
            }
        }
    ]
}`
