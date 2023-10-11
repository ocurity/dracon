subinclude("///k8s//build/defs:k8s")

subinclude("///shell//build_defs:shell")

kustomized_config(
    name = "k8s",
    srcs = [
        "kustomization.yaml",
        "pipeline.yaml",
        "task.yaml",
        ":openapi_schema",
        ":transformer",
    ],
    images = [
        "//third_party/docker:busybox",
    ],
    visibility = ["PUBLIC"],
)

gentest(
    name = "openapi_schema_test",
    data = {
        "CURRENT": ["openapi_schema.json"],
        "EXPECTED": ["//api/openapi/tekton:openapi_schema"],
    },
    exit_on_error = True,
    no_test_output = True,
    test_cmd = """
output=$($TOOLS $DATA_CURRENT $DATA_EXPECTED)
if [[ $output ]]; then
    echo "$output"
    echo "to fix, use:
    \\`\\`\\`
    $ plz run //$PKG:sync_openapi_schema
    \\`\\`\\`
    "
    exit 1
fi
    """,
    test_tools = ["//third_party/binary/josephburnett/jd"],
)

export_file(
    name = "transformer",
    src = "transformer.yaml",
    visibility = ["//api/openapi/tekton/..."],
)

export_file(
    name = "openapi_schema",
    src = "openapi_schema.json",
    visibility = ["//api/openapi/tekton/..."],
)

sh_cmd(
    name = "sync_openapi_schema",
    cmd = """
    set -Eeuo pipefail
    cp \\
        $(out_location //api/openapi/tekton:openapi_schema) \\
        $PKG/openapi_schema.json
    echo "copied $(out_location //api/openapi/tekton:openapi_schema) to $PKG/openapi_schema.json"
    """,
    data = ["//api/openapi/tekton:openapi_schema"],
    shell = "/usr/bin/env bash",
)
