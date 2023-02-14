subinclude("///k8s//build/defs:k8s")

kustomized_config(
    name = "k8s",
    srcs = [
        "kustomization.yaml",
        "my-pipeline.yaml",
        "my-task.yaml",
        "patch-pipeline.yaml",
        ":api",
    ],
)

filegroup(
    name = "api",
    srcs = [
        "//components/base:openapi_schema",
        "//components/base:transformer",
    ],
)

gentest(
    name = "test",
    data = {
        "actual": [":k8s"],
        "expected": ["expected.yaml"],
    },
    no_test_output = True,
    test_cmd = """
    diff <($TOOL -P e '.' "$DATA_ACTUAL") <($TOOL -P e '.' "$DATA_EXPECTED")
    """,
    test_tools = ["//third_party/binary/mikefarah/yq"],
)
