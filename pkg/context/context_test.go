package context

import (
	"fmt"
	"os"
	"strings"
	"testing"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/stretchr/testify/assert"
)

func TestExtractCodeLineRange(t *testing.T) {
	file, err := os.CreateTemp("", "dracon_context_test")
	if err != nil {
		t.Errorf("could not setup tests for context pkg, could not create temporary files")
	}
	defer os.Remove(file.Name())
	if err := os.WriteFile(file.Name(), []byte(code), os.ModeAppend); err != nil {
		t.Errorf("could not setup tests for context pk, could not write temporary file")
	}
	issue := v1.Issue{
		Target:      fmt.Sprintf("%s:%d-%d", file.Name(), 15, 18),
		Type:        "id:985",
		Title:       "SQLI",
		Severity:    v1.Severity_SEVERITY_HIGH,
		Cvss:        9.00,
		Confidence:  v1.Confidence_CONFIDENCE_INFO,
		Description: "found sqli",
		Source:      "",
	}
	codeRange, err := ExtractCode(&issue)
	assert.Nil(t, err)
	assert.Equal(t, strings.Join(strings.Split(code, "\n")[15-DefaultLineRange:18+DefaultLineRange], "\n"), codeRange)
}

func TestExtractCodeLineRangeLessThanDefault(t *testing.T) {
	file, err := os.CreateTemp("", "dracon_context_test")
	if err != nil {
		t.Errorf("could not setup tests for context pkg, could not create temporary files")
	}
	defer os.Remove(file.Name())
	if err := os.WriteFile(file.Name(), []byte(code), os.ModeAppend); err != nil {
		t.Errorf("could not setup tests for context pk, could not write temporary file")
	}
	issue := v1.Issue{
		Target:      fmt.Sprintf("%s:%d-%d", file.Name(), 3, 18),
		Type:        "id:985",
		Title:       "SQLI",
		Severity:    v1.Severity_SEVERITY_HIGH,
		Cvss:        9.00,
		Confidence:  v1.Confidence_CONFIDENCE_INFO,
		Description: "found sqli",
		Source:      "",
	}
	codeRange, err := ExtractCode(&issue)
	assert.Nil(t, err)
	assert.Equal(t, strings.Join(strings.Split(code, "\n")[:18+DefaultLineRange], "\n"), codeRange)
}


func TestExtractCodeLine(t *testing.T) {
	file, err := os.CreateTemp("", "dracon_context_test")
	if err != nil {
		t.Errorf("could not setup tests for context pkg, could not create temporary files")
	}
	defer os.Remove(file.Name())
	if err := os.WriteFile(file.Name(), []byte(code), os.ModeAppend); err != nil {
		t.Errorf("could not setup tests for context pk, could not write temporary file")
	}
	issue := v1.Issue{
		Target:      fmt.Sprintf("%s:%d", file.Name(), 17),
		Type:        "id:985",
		Title:       "SQLI",
		Severity:    v1.Severity_SEVERITY_HIGH,
		Cvss:        9.00,
		Confidence:  v1.Confidence_CONFIDENCE_INFO,
		Description: "found sqli",
		Source:      "",
	}
	codeRange, err := ExtractCode(&issue)
	assert.Nil(t, err)
	assert.Equal(t, strings.Join(strings.Split(code, "\n")[17-DefaultLineRange:17+DefaultLineRange], "\n"), codeRange)
}

func TestExtractCodeInvalidTarget(t *testing.T) {
	// target is ip, url or file that does not exist

	issue := v1.Issue{
		Target:      "/foo/bar:15",
		Type:        "id:985",
		Title:       "SQLI",
		Severity:    v1.Severity_SEVERITY_HIGH,
		Cvss:        9.00,
		Confidence:  v1.Confidence_CONFIDENCE_INFO,
		Description: "found sqli",
		Source:      "",
	}
	_, err := ExtractCode(&issue)
	assert.NotNil(t, err)

	issue.Target = "192.168.1.1"
	_, err = ExtractCode(&issue)
	assert.NotNil(t, err)

	issue.Target = "https://www.example.com?a=9-2"
	_, err = ExtractCode(&issue)
	assert.NotNil(t, err)
}

const code = `from typing import Optional, NamedTuple
from aiopg.connection import Connection

class Student(NamedTuple):
    id: int
    name: str

    @classmethod
    def from_raw(cls, raw: tuple):
        return cls(*raw) if raw else None

    @staticmethod
    async def get(conn: Connection, id_: int):
        async with conn.cursor() as cur:
            await cur.execute(
                'SELECT id, name FROM students WHERE id = %s',
                (id_,),
            )
            r = await cur.fetchone()
            return Student.from_raw(r)

    @staticmethod
    async def get_many(conn: Connection, limit: Optional[int] = None,
                       offset: Optional[int] = None):
        q = 'SELECT id, name FROM students'
        params = {}
        if limit is not None:
            q += ' LIMIT + %(limit)s '
            params['limit'] = limit
        if offset is not None:
            q += ' OFFSET + %(offset)s '
            params['offset'] = offset
        async with conn.cursor() as cur:
            await cur.execute(q, params)
            results = await cur.fetchall()
            return [Student.from_raw(r) for r in results]

    @staticmethod
    async def create(conn: Connection, name: str):
        q = ("INSERT INTO students (name) "
             "VALUES ('%(name)s')" % {'name': name})
        async with conn.cursor() as cur:
            await cur.execute(q)

`
