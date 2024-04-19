# Dracon CDXGEN Producer  

This producer runs [CycloneDX/cdxgen](https://github.com/CycloneDX/cdxgen) against the specified filesystem or image.
It then parses the results into the Dracon format and exits.

## Testing without Dracon

You can run this producer outside of dracon for development with

``` bash
plz run //components/producers/cdxgen:cdxgen -- -in <any cyclonedx sbom document> -out ./cdxgen.pb 
```

cdxgen can be run as a docker image by pulling `ghcr.io/cyclonedx/cdxgen`

## SBOM mode

The producer will output a `LaunchToolResponse` containing a single issue which will have its `CycloneDXSBOM` field populated with the output from cdxgen.
