# Dracon Dependency Track Consumer

This producer imports SBOM results from Dracon producers into [owasp/dependency-track](https://owasp.org/www-project-dependency-track/). It ignores all other results as dependency-track does not do vulnerability management and Dracon does not have any VEX producers yet.

You can use this consumer to generate or keep up to date SBOMs for your projects.

This consumer recognises the annotation from the codeowners enricher and will add a project tag with each username found in the codeowners annotations.
The tag format is `Owner:<username>`

## Testing without Dracon

You can run this producer outside of dracon for development with

``` bash
plz run //components/consumers/dependency-track:dependency-track -- -apiKey=<dependency track api key> -projectName=<name of the project we should upload the bom to> -url=<where to find dependency track> -in <where to find dracon results> --projectUUID=<matching uuid of the target dependency track project>
```
