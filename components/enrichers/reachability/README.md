# Reachability Enricher

This enricher takes findings and checks if the target is reachable with a 
reachables slice produced by
[appthreat/atom](https://github.com/appthreat/atom).

For each finding, it adds the following annotation.
"Reachable:<true/false>"
