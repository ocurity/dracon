# Improve Observability Via Intermediate Result Cache

## Introduction

Currently our components use a Tekton workspace to save all intermediate
results. This is bad for observability and also does not provide long term
attestations on what happened, when and how.

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
* all components are refactored to use those methods only
