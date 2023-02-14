subinclude("///k8s//build/defs:k8s")

kustomized_config(
    name = "k8s",
    srcs = [
        "kustomization.yaml",
    ],
    deps = [
        "//api/openapi/tekton/tests/pipeline/dag/component-patch/base:k8s",
        "//api/openapi/tekton/tests/pipeline/dag/component-patch/components/my-component:k8s",
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
