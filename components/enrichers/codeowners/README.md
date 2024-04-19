# CodeOwners Enricher

This enricher scans the cloned source for
[CODEOWNERS](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners)
files, For each finding, it adds the following annotation.
"Owner-$incremental_number:$the_username_of_the_owner>"
