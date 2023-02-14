# Dracon Dependency Track Consumer

This producer imports SBOM results from Dracon producers into [owasp/dependency-track](https://owasp.org/www-project-dependency-track/). It ignores all other results as dependency-track does not do vulnerability management and Dracon does not have any VEX producers yet.

You can use this producer to generate or keep up to date SBOMs for your projects. 


## Testing without Dracon

You can run this producer outside of dracon for development with

``` bash
plz run //components/consumers/dependency-track:dependency-track -- -apiKey=<dependency track api key> -projectName=<name of the project we should upload the bom to> -url=<where to find dependency track> -in <where to find dracon results> --projectUUID=<matching uuid of the target dependency track project>
```
