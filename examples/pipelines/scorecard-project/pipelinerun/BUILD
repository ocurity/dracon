genrule(
    name = "pipelinerun",
    srcs = ["pipelinerun.tmpl.yaml"],
    outs = ["pipelinerun.yaml"],
    cmd = """
if [ -z "${GITHUB_TOKEN:-}" ]; then
    echo "'GITHUB_TOKEN' is not set."
    exit 1
fi
mv $SRCS $OUTS
sed -i "s#\\\\$GITHUB_TOKEN#$GITHUB_TOKEN#g" $OUTS
    """,
    pass_env = ["GITHUB_TOKEN"],
    visibility = ["//examples/pipelines/scorecard-project/..."],
)
