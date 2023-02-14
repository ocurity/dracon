genrule(
    name = "openapi_schema",
    srcs = ["//third_party/k8s/tektoncd/pipeline:v1beta1_swagger"],
    outs = ["openapi_schema.json"],
    cmd = "$TOOL",
    tools = [":generate_openapi_schema_tool"],
    visibility = ["//components/base/..."],
)

sh_binary(
    name = "generate_openapi_schema_tool",
    main = ".generate_openapi_schema_tool.sh",
    deps = [
        "//third_party/binary/stedolan/jq",
    ],
)
