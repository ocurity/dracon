subinclude("//build/defs:dracon")

dracon_pipeline(
    name = "e2e-project",
    components = [
        "//components/base:k8s",
        "//components/consumers/mongodb:k8s",
        "//components/consumers/elasticsearch:k8s",
        
        "//components/enrichers/aggregator:k8s",
        "//components/enrichers/deduplication:k8s",
        "//components/enrichers/policy:k8s",
        
        "//components/producers/aggregator:k8s",

        "//components/producers/zaproxy:k8s",
        "//components/producers/typescript-eslint:k8s",
        "//components/producers/typescript-npm-audit:k8s",
        "//components/producers/testsslsh:k8s",
        "//components/producers/terraform-tfsec:k8s",
        "//components/producers/semgrep:k8s",
        "//components/producers/ossf-scorecard:k8s",
        "//components/producers/kics:k8s",
        "//components/producers/docker-trivy:k8s",
        "//components/producers/python-bandit:k8s",
        "//components/producers/python-pip-safety:k8s",
        "//components/producers/golang-gosec:k8s",
        "//components/producers/golang-nancy:k8s",
        
        "//components/sources/git:k8s",
    ],
    kube_context = "//build/k8s/k3d:dracon",
    kustomization_yaml = "kustomization.yaml",
    pipelinerun = "//examples/pipelines/e2e-project/pipelinerun:pipelinerun",
)

