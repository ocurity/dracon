package main

import (
	"encoding/json"
	"testing"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components/producers/typescript-eslint/types"

	"github.com/stretchr/testify/assert"
)

const exampleOutput = `[{"filePath":"/foo/bar/js/types/File.ts",
"messages":[
	{"ruleId":"@typescript-eslint/explicit-module-boundary-types",
	"severity":1,
	"message":"Missing return type on function.",
	"line":21,
	"column":46,
	"nodeType":"ArrowFunctionExpression",
	"messageId":"missingReturnType",
	"endLine":160,"endColumn":104},

	{"ruleId":"jsdoc/require-jsdoc","severity":1,"message":"Missing JSDoc comment.",
	"line":29,"column":1,"nodeType":"FunctionDeclaration",
	"messageId":"missingJsDoc","endColumn":null,
	"fix":{"range":[890,890],"text":"/**\n *\n */\n"}}],

	"errorCount":0,"warningCount":2,
	"fixableErrorCount":0,
	"fixableWarningCount":1,
	"source":"import {withPrefix} from 'gatsby';\nimport React from 'react';\n\nimport {NavigationNode} from '../../../types';\n\nimport {PrimaryMenuLink} from '../../atoms/Links/PrimaryMenuLink';\n\nimport {SidebarNavigationSecondaryLinks} from './SidebarNavigationSecondaryLinks';\nimport {markStringForWrappingInNav} from './utils';\n\ninterface SidebarNavigationPrimaryLinksProps {\n  navNodes: NavigationNode[];\n  route: string;\n}\n\nexport const SidebarNavigationPrimaryLinks = ({navNodes, route}: SidebarNavigationPrimaryLinksProps) => (\n  <ol>\n    {navNodes.map(({link, title, children}) => (\n      <li key={link}>\n        <PrimaryMenuLink active={isActive(withPrefix(link), route)} route={encodeURI(withPrefix(link))}>\n          {markStringForWrappingInNav(title) || ''}\n        </PrimaryMenuLink>\n        <SidebarNavigationSecondaryLinks navNodes={children} route={route} />\n      </li>\n    ))}\n  </ol>\n);\n\nfunction isActive(link: string, route: string) {\n  const nodeIsActive = link === route?.split('#')[0];\n  const nodeChildIsActive = route?.startsWith(link);\n\n  return nodeIsActive || nodeChildIsActive;\n}\n",

	"usedDeprecatedRules":[{"ruleId":"jsx-a11y/accessible-emoji","replacedBy":[]},{"ruleId":"jsx-a11y/label-has-for","replacedBy":[]}	]
	}]`

func TestParseIssues(t *testing.T) {
	var results []types.ESLintIssue
	err := json.Unmarshal([]byte(exampleOutput), &results)
	assert.Nil(t, err)
	issues := parseIssues(results)

	expectedIssue := &v1.Issue{
		Target:      "/foo/bar/js/types/File.ts:21:46",
		Type:        "@typescript-eslint/explicit-module-boundary-types",
		Title:       "@typescript-eslint/explicit-module-boundary-types",
		Severity:    v1.Severity_SEVERITY_MEDIUM,
		Cvss:        0.0,
		Confidence:  v1.Confidence_CONFIDENCE_MEDIUM,
		Description: "Missing return type on function.",
	}
	issue2 := &v1.Issue{
		Target:      "/foo/bar/js/types/File.ts:29:1",
		Type:        "jsdoc/require-jsdoc",
		Title:       "jsdoc/require-jsdoc",
		Severity:    v1.Severity_SEVERITY_MEDIUM,
		Cvss:        0.0,
		Confidence:  v1.Confidence_CONFIDENCE_MEDIUM,
		Description: "Missing JSDoc comment.",
	}
	assert.Equal(t, expectedIssue, issues[0])
	assert.Equal(t, issue2, issues[1])
}
