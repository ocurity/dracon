# deps.dev Enricher

*WARNING*: as stated [here](https://github.com/ocurity/dracon/pull/15#discussion_r1125438946)
_Because Go licenses must be detected, there's always a chance of the detected license being wrong_. The licenses detected for GO might not be 100% accurate. Please ensure you validate the results if you intent to use them in a legally binding way.

This enricher implements a rudimentary client for Google's [Open Source Insights](https://deps.dev) project.
It *ONLY* works with issues containing SBOMs and at its current version adds missing licensing information to each package in the SBOM.



## Future steps

When/If there is interest around this enricher we can add more information retrieval such as Security Advisories.