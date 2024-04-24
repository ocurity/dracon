## Writing Consumers

A consumer is a program that parses the Dracon compatible outputs and pushes them into arbitrary destinations. The Dracon compatible outputs from from _producers_ and _enrichers_.

---

Consumers can be written in any language that supports protobufs. We currently have examples in Golang and Python. They are all structured in the same way:

1. Parse program arguments:
   1. `in`: the dracon compatible outputs location.
   2. `raw`: whether or not to use enriched results.
2. Parse all dracon compatible output files the `in` location.
3. Do arbitrary logic with issues.
4. Create a Tekton Task `task.yaml` with `.metadata.labels["v1.dracon.ocurity.com/component"] = consumer`.

### Consumer API

For convenience, there are helper functions in the `./consumers` pkg/module for Golang/Python.

See the godoc for more information.
