#!/usr/bin/env bash
# This script waits for pipelineruns to finish and exits depending on if all pipelineruns were successful or not.
set -Eeuo pipefail

source "//build/util"

YQ_BIN="//third_party/binary/mikefarah/yq"

function all_pipelineruns_succeeded {
    mapfile -t pipelinerun_statuses < \
        <(kubectl get pipelineruns.tekton.dev --all-namespaces -oyaml \
            | "$YQ_BIN" e -N ".items[] \
                | .metadata.namespace + \"/\" + .metadata.name + \",\" \
                + .status.conditions[].reason + \",\" \
                + .status.conditions[].message + \",\"
            " -
        )
    mapfile -t pipelinerun_number < <(kubectl get pipelineruns.tekton.dev --all-namespaces | wc -l)
    
    let "num=$pipelinerun_number - 1"

    if [[ ${#pipelinerun_statuses[@]} -ne ${num} ]]; then
        util::info "status numbers retrieved ${#pipelinerun_statuses[@]} not equal to number of pipelineruns in the api ${num} this usually means the server is struggling, will wait"
        return 2
    fi

    has_pipeline_running=false
    has_pipeline_error=false
    mapfile -t pipelinerun_namespace_names < \
        <(printf '%s\n' "${pipelinerun_taskrun_statuses[@]}" | grep "." | cut -f1 -d, | sort -u)

    util::info "waiting for ${#pipelinerun_statuses[@]} PipelineRun(s) to complete..."

    for pipelinerun_status in "${pipelinerun_statuses[@]}"; do
        namespace="$(echo "$pipelinerun_status" | cut -f1 -d, | cut -f1 -d/)"
        name="$(echo "$pipelinerun_status" | cut -f1 -d, | cut -f2 -d/)"
        status="$(echo "$pipelinerun_status" | cut -f2 -d,)"
        pipelinerun_message="$(echo "$pipelinerun_status" | cut -f3 -d,)"

        case "$status" in
            "Running")
                mapfile -t taskrun_ids < \
                    <(kubectl \
                        --namespace "$namespace" \
                        get pipelineruns.tekton.dev -oyaml \
                        | "$YQ_BIN" e -N ".status.childReferences[].name" -
                    )

                mapfile -t taskrun_statuses < \
                    <(kubectl \
                        --namespace "$namespace" \
                        get taskruns.tekton.dev -oyaml \
                        "${taskrun_ids[@]}" \
                        | "$YQ_BIN" e -N ".items[] \
                            | .metadata.namespace + \"/\" + .metadata.name + \",\" \
                            + .status.conditions[].message + \",\" \
                            + .status.conditions[].status + \",\"
                        " -
                    )

                util::info "$name: $status, $pipelinerun_message"
                for taskrun_status in "${taskrun_statuses[@]}"; do
                    util::info "    TaskRun: $taskrun_status"
                done
                has_pipeline_running=true
            ;;
            "Succeeded")
                util::success "PipelineRun: $name: $status, $pipelinerun_message"
            ;;
            "null")
                util::info "PipelineRun: $name: $status, $pipelinerun_message"
                has_pipeline_running=true
                continue
            ;;
            "")
                util::info "PipelineRun: $name: $status, $pipelinerun_message"
                has_pipeline_running=true
                continue
            ;;
            *)
                mapfile -t taskrun_ids < \
                    <(kubectl \
                        --namespace "$namespace" \
                        get pipelineruns.tekton.dev -oyaml \
                        | "$YQ_BIN" e -N ".status.childReferences[].name" -
                    )

                mapfile -t taskrun_statuses < \
                    <(kubectl \
                        --namespace "$namespace" \
                        get taskruns.tekton.dev -oyaml \
                        "${taskrun_ids[@]}" \
                        | "$YQ_BIN" e -N ".items[] \
                            | .metadata.namespace + \"/\" + .metadata.name + \",\" \
                            + .status.conditions[].message + \",\" \
                            + .status.conditions[].status + \",\"
                        " -
                    )

                util::error "PipelineRun: $name: $status, $pipelinerun_message"
                for taskrun_status in "${taskrun_statuses[@]}"; do
                    type="$(echo "$taskrun_status" | cut -f3 -d,)"
                    if [[ "$type" != "True" ]]; then
                        taskrun_namespace="$(echo "$taskrun_status" | cut -f1 -d, | cut -f1 -d/)"
                        taskrun_name="$(echo "$taskrun_status" | cut -f1 -d, | cut -f2 -d/)"
                        taskrun_message="$(echo "$taskrun_status" | cut -f2 -d,)"

                        # print logs for pod associated with task run.
                        if [[ ! -z "$(echo "$taskrun_name" | xargs)" ]]; then
                            util::error "    TaskRun: $taskrun_name: $taskrun_message"
                            pod_name="$(kubectl \
                                --namespace "$taskrun_namespace" \
                                get taskruns.tekton.dev \
                                "$taskrun_name" -oyaml \
                                | "$YQ_BIN" e -N ".status.podName" -
                            )"
                            if [[ ! -z "$(echo "$pod_name" | xargs)" ]]; then
                                util::info "    $pod_name Logs:"
                                kubectl --namespace "$taskrun_namespace" logs \
                                    --all-containers "$pod_name" \
                                    | sed 's/^/        /g' -
                            fi
                        fi
                    fi
                done
                # # print all pod logs for pipelinerun
                # mapfile -t pipelinerun_pods_namespace_names < \
                #     <(kubectl get pod --all-namespaces --selector="tekton.dev/pipelineRun=$name" -o=jsonpath='{range .items[*]}{.metadata.namespace}{","}{.metadata.name}{"\n"}{end}')
                # for pipeline_pod_namespace_name in "${pipelinerun_pods_namespace_names[@]}"; do
                #     namespace="$(echo "$pipeline_pod_namespace_name" | cut -f1 -d,)"
                #     pod_name="$(echo "$pipeline_pod_namespace_name" | cut -f2 -d,)"
                #     logs="$(kubectl -n "$namespace" logs "$pod_name" 2>&1 --all-containers)"
                #     util::error "Pod: $namespace/$pod_name logs:"
                #     echo "$logs"
                # done
                has_pipeline_error=true
        esac
    done
    echo ""

    if [ "$has_pipeline_error" = true ]; then
        return 1
    fi

    if [ "$has_pipeline_running" = true ]; then
        return 2
    fi

    return 0
}



time_limit_secs=360000
sleep_interval=5
intervals=$((time_limit_secs/sleep_interval))
attempts=0

util::info "waiting for all pipelines to succeed within ${time_limit_secs}s"

set +x
until all_pipelineruns_succeeded; do
    ec="$?"
    if [ "$ec" == 1 ]; then
        exit 1
    fi
    if [ $attempts -eq $intervals ]; then
        util::error "timed out"
        exit 1
    fi
    attempts=$((attempts + 1))
    sleep $sleep_interval
done
set -e

util::success "all pipelineruns completed successfully"
