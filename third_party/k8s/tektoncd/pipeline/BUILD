subinclude("///shell//build_defs:shell")

subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:k8s",
    "///k8s//build/defs:helm",
)

VERSION = "0.44.0"

remote_file(
    name = "v1beta1_swagger",
    url = f"https://raw.githubusercontent.com/tektoncd/pipeline/v{VERSION}/pkg/apis/pipeline/v1beta1/swagger.json",
    visibility = ["//api/openapi/tekton/..."],
)

remote_file(
    name = "_release",
    url = f"https://storage.googleapis.com/tekton-releases/pipeline/previous/v{VERSION}/release.yaml",
)

# use this script to get the list of images to mirror for the current Tekton
# Release (for the `images_to_mirror` list below).
sh_cmd(
    name = "images_to_mirror",
    # we have to skip the mcr.microsoft.com/powershell container as it doesn't
    # exist for linux.
    cmd = """
    grep -oE '[^ "]*@sha256:[^ "]*' $(out_location :_release) \\
    | grep -v "mcr.microsoft.com/powershell" \\
    | sed 's/^/"/g' \\
    | sed 's/$/"/g' \\
    | sort -u \\
    | sed 's/:.*@/@/g' \\
    | sed '\\\\$!s/\\\\$/,/'
    """,
    data = [":_release"],
)

images_to_mirror = [
    "cgr.dev/chainguard/busybox@sha256:19f02276bf8dbdd62f069b922f10c65262cc34b710eea26ff928129a736be791",
    "gcr.io/google.com/cloudsdktool/cloud-sdk@sha256:27b2c22bf259d9bc1a291e99c63791ba0c27a04d2db0a43241ba0f1f20f4067f",
    "gcr.io/tekton-releases/github.com/tektoncd/pipeline/cmd/controller@sha256:51e4c16fc4d0b18912106ff7fe37cc24c0667de6af468cdd1a9e9dc174039de1",
    "gcr.io/tekton-releases/github.com/tektoncd/pipeline/cmd/entrypoint@sha256:0be7d0cd52ba04baa136f987c74d80d171e3c1c86d47db8368825cfbe4d8df0b",
    "gcr.io/tekton-releases/github.com/tektoncd/pipeline/cmd/git-init@sha256:baa4110837ef91c3ec26a00f3c5d980fef174b97934cf2fc9115092ec78d35ee",
    "gcr.io/tekton-releases/github.com/tektoncd/pipeline/cmd/imagedigestexporter@sha256:b65389daa312b9c1bfe9cc1d2c9066a877676034b4264b3576e3ccf8843b961d",
    "gcr.io/tekton-releases/github.com/tektoncd/pipeline/cmd/nop@sha256:779508dcf27d658148d532d59ff505d5afc45c6aa07111d12169a90be6a6a9d8",
    "gcr.io/tekton-releases/github.com/tektoncd/pipeline/cmd/pullrequest-init@sha256:c80af31a67027b9e9dc9b710856ecda8137348bb82aefb1dd59949039ea17fc2",
    "gcr.io/tekton-releases/github.com/tektoncd/pipeline/cmd/resolvers@sha256:9e156710dc7730487640faa798a7f70e32f15bd31ec9f50cab2c7f5e154a3dfe",
    "gcr.io/tekton-releases/github.com/tektoncd/pipeline/cmd/sidecarlogresults@sha256:4f3ac70c281f60de5634e3c8d4e17210a8d64aa6ae2d933739dbfc3cb674b8d0",
    "gcr.io/tekton-releases/github.com/tektoncd/pipeline/cmd/webhook@sha256:f9649cd6efac896defc9da687d0c5e6005e3c2c2e323aed3653917b29c257a5e",
    "gcr.io/tekton-releases/github.com/tektoncd/pipeline/cmd/workingdirinit@sha256:58d1145c65711eda6198c5b8d3a0d36955f0a56d732c961d29cab4381cbce644",
]

for image_to_mirror in images_to_mirror:
    repo = image_to_mirror.split("@")[0]
    digest = image_to_mirror.split("@")[1]

    buildkit_image_mirror(
        name = basename(repo),
        digest = digest,
        repo = repo,
        tags = [],
    )

kustomized_config(
    name = "release_no_namespace",
    srcs = [
        "kustomization.yaml",
        ":_release",
    ],
)

helm_chart(
    name = "pipeline",
    images = [":" + basename(image.split("@")[0]) for image in images_to_mirror],
    src_templates = [":release_no_namespace"],
    visibility = ["//build/k8s/tektoncd/pipeline/..."],
)
