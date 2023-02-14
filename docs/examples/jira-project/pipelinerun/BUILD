# genrule(
#     name = "pipelinerun",
#     srcs = ["pipelinerun.tmpl.yaml"],
#     outs = ["pipelinerun.yaml"],
#     cmd = """
# if [ -z "${CONSUMER_JIRA_URL:-}" ]; then
#     echo "'CONSUMER_JIRA_URL' is not set."
#     exit 0
# fi
# if [ -z "${CONSUMER_JIRA_API_TOKEN:-}" ]; then
#     echo "'CONSUMER_JIRA_API_TOKEN' is not set."
#     exit 0
# fi
# if [ -z "${CONSUMER_JIRA_USER:-}" ]; then
#     echo "'CONSUMER_JIRA_USER' is not set."
#     exit 0
# fi
# mv $SRCS $OUTS
# sed -i "s#\\\\$CONSUMER_JIRA_URL#$CONSUMER_JIRA_URL#g" $OUTS
# sed -i "s#\\\\$CONSUMER_JIRA_API_TOKEN#$CONSUMER_JIRA_API_TOKEN#g" $OUTS
# sed -i "s#\\\\$CONSUMER_JIRA_USER#$CONSUMER_JIRA_USER#g" $OUTS
#     """,
#     pass_env = [
#         "CONSUMER_JIRA_URL",
#         "CONSUMER_JIRA_API_TOKEN",
#         "CONSUMER_JIRA_USER",
#     ],
#     visibility = ["//examples/pipelines/jira-project/..."],
# )
