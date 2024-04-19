# Policy Enricher

This enricher implements an MVP of the design outlined in
[the enricher design document](/docs/designs/policy-enricher.md).
In its current state it evaluates the `allow` rule of a single policy found
under the `policy_path` variable.

## Future steps

When/If there is interest around the policy enricher it should be modified to
work with an s3 compatible API and load policies more efficiently based on
either tool name or pipeline name.
