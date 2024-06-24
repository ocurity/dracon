# Improve Observability Via Intermediate Result Cache

## Introduction

Currently our components use a Tekton workspace to save all intermediate
results. This is simple implementation wise but:

* It's not great  for observability
* It does not provide long term attestations on what happened, when and how
* If a step fails it also prevents the step from being retriggered
  (future feature)

### Suggested solution

Slightly refactor `putil` to allow for reading from and writing to more urls
than `file://` (start with `postgresql` since we have an enrichment db)

## Design

### SDK changes

`putil/load.go` and `putil/write.go` get new methods that:

* can resolve a read-from or write-from url to supported intermediate results
  locations (e.g. `kafka://<endpoint>/topic` or `postgres://` or
  `file:///absolute/path`)
* can use a `postgresql` database and write intermediate, `LaunchToolResponse`
  and `EnrichedLaunchToolResponse` results to a table
* SDK is refactored to use those methods only

#### `WriteResults` method

Method `putil.WriteResults` changes instead of `outFile` to take
`writeTarget string` if the target is a postgres connection string, then it
calls `writeResultsToDB`. Alternatively, if the method is a filesystem path or
starts with `file://` then it calls `writeResultsToDisk` which functions like
the current WriteResults method.

#### Method `writeResultsToDB`

A new method is created
`writeResultsToDB(string, LaunchToolResponse|EnrichedLaunchToolResponse)error`
which given a DB connection string and a Raw or Enriched LaunchToolResponse,
appends the marshalled protobuf either to the `LaunchToolResponse` or
`EnrichedLaunchToolResponse` table.

#### `LoadTaggedToolResponse` method

The method `putil.LoadTaggedToolResponse` changes.
Instead of accepting `inPath string` it accepts `inTarget string`.
Similar to WriteResults depending on the `inTarget` format it either calls
`readRawResultsFromDB` or `ReadRawResultsFromDrive` which does what
`LoadTaggedToolResponse` currently does.

#### Method `readRawResultsFromDB`

The new method `readRawResultsFromDB(connectionSTR scanID string)([]*LaunchToolResponse, error)`
does a `SELECT * from rawResults WHERE scanID={};` and recreates the
`LaunchToolResponse`s based on tool and scan start time information.

#### Method `readEnrichedResultsFromDB`

The new method `readEnrichedResultsFromDB(connectionSTR scanID string)([]*LaunchToolResponse, error)`
does a `SELECT * from EnrichedResults WHERE scanID={};` and recreates the
`EnrichedLaunchToolResponse`s based on tool and scan start time information.

#### `*Aggregated` methods

The current `load.go` file contains also method `LoadEnrichedNonAggregatedToolResponse`
This is another case of a slightly different `Select` query

#### Database Tables

The `intermediate-results` database needs new tables that reflect a
flattened version of `issue` and `enrichedIssue` the tables look like this:

##### Table(s)

This is temporary information which we expect to delete soon-ish after creation.
Therefore it makes sense to create a single table. The table contains a
flattened version of each Issue as follows:

```
Issue{
id : string # <-- unique Issue ID
scan_uuid: string
scan_start_time: timestamp
scan_tags: string # a Join("<>",scan_tags)
tool_name: string
cwe: string # a Join("<>",cwe)
context_segment: bigString
cyclone_d_x_s_b_o_m: bigString
uuid: string  # <-- what Dracon components thing is the UUID of this issue, different than the ID because of enrichers and aggregation
cve: string # <-- a join("<>", cve)
source: string
description: string # bigString
confidence: string 
cvss: string
severity: string
title: string
type: string
target: string
first_seen: timestamp
count: int
false_positive: bool
updated_at: timestamp
hash: string
annotations: string # a Join("<>",annotations)
enriched: bool
}
```

### Component Changes

All our components take `--in` and `--out` args, they don't process the args and
pass them transparently to `putil` so no component changes required
