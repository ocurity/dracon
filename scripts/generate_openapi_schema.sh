#!/usr/bin/env bash

set -euo pipefail

elementIn () {
  local e match="$1"
  shift
  for e; do [[ "$e" == "$match" ]] && return 0; done
  return 1
}

source ./scripts/util.sh

if [ "$#" -lt 2 ]
then
    util::error "You need to provide 2 positional arguments for this tool: $0 <input.json> <output.json>"
    exit 1
fi

# only use .definitions
cp $1 $2

# add x-kubernetes-group-version-kind to all objects
crds=($(jq -r '.definitions | keys | .[]' $2))
for crd in "${crds[@]}"; do
    gvk_crd="${crd//v1alpha1/tekton.dev\/v1alpha1}"
    gvk_crd="${gvk_crd//v1beta1/tekton.dev\/v1beta1}"
    kind="$(echo $gvk_crd | rev | cut -d\. -f1 | rev)"
    version="$(echo $gvk_crd | rev | cut -d\. -f2 | rev | cut -d\/ -f2)"
    group="$(echo $gvk_crd | cut -d\/ -f1)"
    if [ -n "$group" ]; then
        jq \
            --arg group "$group" \
            --arg version "$version" \
            --arg kind "$kind" \
            --arg crd "$crd" \
        '.definitions[$crd] += { "x-kubernetes-group-version-kind": [{ "group": $group, "version": $version, "kind": $kind }] }' $2 > $2.new
        mv $2.new $2
    fi
done

# remove v1alpha1 crds
for crd in "${crds[@]}"; do
    if [[ $crd == v1alpha1* ]]; then
        echo "removing $crd"
        jq --arg crd "$crd" 'del(.definitions[$crd])' $2 > $2.new
        mv $2.new $2
    fi
done
crds=($(jq -r '.definitions | keys | .[]' $2))

# strip \$refs which don't reference in the same schema
all_refs=($(jq -r 'paths | join("][")' $2 | grep "\$ref$" | while read line; do echo "[${line}]"; done | sed -e 's|\[|\["|g' | sed -e 's|\]|"\]|g'))
for ref in "${all_refs[@]}"; do
    ref_value=$(jq -r ".$ref" $2 | sed 's|#/definitions/||g')
    if ! elementIn "${ref_value}" "${crds[@]}"; then
        echo "removing: ${ref} (${ref_value} does not exist)"
        jq "del(.$ref)" $2 > $2.new
        mv $2.new $2
    fi
done

# add x-kubernetes-patch-merge-key and x-kubernetes-patch-strategy
# allow the Pipeline.spec.tasks list to be merged
jq '.definitions["v1beta1.PipelineSpec"].properties.tasks += { "x-kubernetes-patch-strategy": "merge" }' $2 > $2.new && mv $2.new $2
jq '.definitions["v1beta1.PipelineSpec"].properties.tasks += { "x-kubernetes-patch-merge-key": "name" }' $2 > $2.new && mv $2.new $2

# allow the Pipeline.spec.params list to be merged
jq '.definitions["v1beta1.PipelineSpec"].properties.params += { "x-kubernetes-patch-strategy": "merge" }' $2 > $2.new && mv $2.new $2
jq '.definitions["v1beta1.PipelineSpec"].properties.params += { "x-kubernetes-patch-merge-key": "name" }' $2 > $2.new && mv $2.new $2

# allow the Pipeline.spec.tasks[].params to be merged
jq '.definitions["v1beta1.PipelineTask"].properties.params += { "x-kubernetes-list-type": "map" }' $2 > $2.new && mv $2.new $2
jq '.definitions["v1beta1.PipelineTask"].properties.params += { "x-kubernetes-patch-strategy": "merge" }' $2 > $2.new && mv $2.new $2
jq '.definitions["v1beta1.PipelineTask"].properties.params += { "x-kubernetes-patch-merge-key": "name" }' $2 > $2.new && mv $2.new $2

# butcher the Pipeline.spec.tasks[].params[].value so that they can be merged additively
jq '.definitions["v1beta1.Param"].properties.value += { "type": "array", "items": { "type": "string" } }' $2 > $2.new && mv $2.new $2
jq 'del(.definitions["v1beta1.Param"].properties.value["$ref"])' $2 > $2.new && mv $2.new $2
jq 'del(.definitions["v1beta1.Param"].properties.value.default)' $2 > $2.new && mv $2.new $2
jq '.definitions["v1beta1.Param"].properties.value += { "x-kubernetes-patch-strategy": "merge" }' $2 > $2.new && mv $2.new $2
jq '.definitions["v1beta1.Param"].properties.value += { "x-kubernetes-list-type": "map" }' $2 > $2.new && mv $2.new $2


# allow the Task.spec.steps list to be merged
jq '.definitions["v1beta1.TaskSpec"].properties.steps += { "x-kubernetes-patch-strategy": "merge" }' $2 > $2.new && mv $2.new $2
jq '.definitions["v1beta1.TaskSpec"].properties.steps += { "x-kubernetes-patch-merge-key": "name" }' $2 > $2.new && mv $2.new $2

# allow the Task.spec.params list to be merged
jq '.definitions["v1beta1.TaskSpec"].properties.params += { "x-kubernetes-patch-strategy": "merge" }' $2 > $2.new && mv $2.new $2
jq '.definitions["v1beta1.TaskSpec"].properties.params += { "x-kubernetes-patch-merge-key": "name" }' $2 > $2.new && mv $2.new $2

# allow the Task.spec.volumes list to be merged
jq '.definitions["v1beta1.TaskSpec"].properties.volumes += { "x-kubernetes-patch-strategy": "merge" }' $2 > $2.new && mv $2.new $2
jq '.definitions["v1beta1.TaskSpec"].properties.volumes += { "x-kubernetes-patch-merge-key": "name" }' $2 > $2.new && mv $2.new $2

# allow the Task.spec.workspaces list to be merged
jq '.definitions["v1beta1.TaskSpec"].properties.workspaces += { "x-kubernetes-patch-strategy": "merge" }' $2 > $2.new && mv $2.new $2
jq '.definitions["v1beta1.TaskSpec"].properties.workspaces += { "x-kubernetes-patch-merge-key": "name" }' $2 > $2.new && mv $2.new $2

# allow the Task.spec.results list to be merged
jq '.definitions["v1beta1.TaskSpec"].properties.results += { "x-kubernetes-patch-strategy": "merge" }' $2 > $2.new && mv $2.new $2
jq '.definitions["v1beta1.TaskSpec"].properties.results += { "x-kubernetes-patch-merge-key": "name" }' $2 > $2.new && mv $2.new $2
