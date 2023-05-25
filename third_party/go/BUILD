subinclude("///go//build_defs:go")

subinclude("///shell//build_defs:shell")

go_toolchain(
    name = "toolchain",
    version = "1.20.4",
)

sh_cmd(
    name = "gen_modules",
    cmd = """
    out="\\\\$(./pleasew query reporoot)/third_party/go/modules.build_defs"
    echo "GO_MODULES = {" > "\\\\$out.new"

    ./pleasew run ///go//tools/please_go:bootstrap -- get -m go.mod \\
    | sed 's#go_repo(module=##g' \\
    | sed 's#, version=#: #g' \\
    | sed 's#)#,#g' \\
    | sed 's#^#  #g' \\
    >> "\\\\$out.new"

    echo "}" >> "\\\\$out.new"
    mv "\\\\$out.new" "\\\\$out"
    """,
)

export_file(
    name = "modules",
    src = "modules.build_defs",
)

subinclude(":modules")

for mod, ver in GO_MODULES.items():
    name = mod.replace("/", "_")
    download = None
    install = []

    if mod == "github.com/onsi/ginkgo":
        download = go_mod_download(
            name = tag(name, "dl"),
            module = mod,
            strip = ["integration"],
            version = ver,
        )

    if mod == "google.golang.org/protobuf":
        install = [
            "proto",
            "compiler/protogen",
            "internal/version",
            "cmd/protoc-gen-go/internal_gengo",
            "reflect/protoreflect",
            "runtime/protoimpl",
            "types/known/timestamppb",
            "encoding/protojson",
        ]

    if mod == "github.com/apache/arrow/go/v12":
        download = go_mod_download(
            name = tag(name, "dl"),
            module = mod,
            strip = [
                # "internal/utils/_lib/arch.h",
                # "internal/utils/_lib/min_max_avx2_amd64.s",
                # "internal/utils/_lib/min_max_neon.s",
                "internal/utils/_lib/min_max_sse4_amd64.s",
                # "internal/utils/_lib/min_max.c",
                # "internal/utils/_lib/transpose_ints_avx2_amd64.s",
                "internal/utils/_lib/transpose_ints_sse4_amd64.s",
                # "internal/utils/_lib/transpose_ints.c",
                # "internal/utils/buf_reader.go",
                # "internal/utils/endians_default.go",
                # "internal/utils/endians_s390x.go",
                # "internal/utils/Makefile",
                # "internal/utils/math.go",
                "internal/utils/min_max_amd64.go",
                "internal/utils/min_max_arm64.go",
                # "internal/utils/min_max_avx2_amd64.go",
                # "internal/utils/min_max_avx2_amd64.s",
                "internal/utils/min_max_neon_arm64.go",
                "internal/utils/min_max_neon_arm64.s",
                # "internal/utils/min_max_noasm.go",
                "internal/utils/min_max_ppc64le.go",
                "internal/utils/min_max_s390x.go",
                "internal/utils/min_max_sse4_amd64.go",
                "internal/utils/min_max_sse4_amd64.s",
                # "internal/utils/min_max.go",
                # "internal/utils/transpose_ints_amd64.go.tmpl",
                # "internal/utils/transpose_ints_amd64.go",
                # "internal/utils/transpose_ints_arm64.go",
                # "internal/utils/transpose_ints_avx2_amd64.go",
                # "internal/utils/transpose_ints_avx2_amd64.s",
                # "internal/utils/transpose_ints_def.go",
                # "internal/utils/transpose_ints_noasm.go.tmpl",
                # "internal/utils/transpose_ints_noasm.go",
                # "internal/utils/transpose_ints_ppc64le.go",
                # "internal/utils/transpose_ints_s390x.go.tmpl",
                # "internal/utils/transpose_ints_s390x.go",
                # "internal/utils/transpose_ints_simd.go.tmpl",
                # "internal/utils/transpose_ints_sse4_amd64.go",
                # "internal/utils/transpose_ints_sse4_amd64.s",
                # "internal/utils/transpose_ints_test.go",
                # "internal/utils/transpose_ints.go.tmpl",
                # "internal/utils/transpose_ints.go",
                # "internal/utils/transpose_ints.tmpldata",
            ],
            version = ver,
        )

    if mod == "github.com/google/flatbuffers":
        download = go_mod_download(
            name = tag(name, "dl"),
            module = mod,
            strip = [
                "tests",
            ],
            version = ver,
        )

    if mod == "google.golang.org/grpc":
        install = [
            ".",
            # "compiler/protogen",
            # "internal/version",
            # "cmd/protoc-gen-go/internal_gengo",
        ]

    go_repo(
        name = name,
        download = download,
        install = install,
        module = mod,
        version = ver,
        visibility = ["PUBLIC"],
    )

go_module(
    name = "protoc-gen-go",
    binary = True,
    download = ":_google.golang.org_protobuf#dl",
    install = ["cmd/protoc-gen-go"],
    module = "google.golang.org/protobuf",
    visibility = ["PUBLIC"],
    deps = [":google.golang.org_protobuf"],
)

go_module(
    name = "protoc-gen-go-grpc",
    binary = True,
    module = "google.golang.org/grpc/cmd/protoc-gen-go-grpc",
    version = "v1.3.0",
    visibility = ["PUBLIC"],
    deps = [
        ":google.golang.org_protobuf",
    ],
)

go_module(
    name = "mvdan.cc_gofumpt",
    binary = True,
    module = "mvdan.cc/gofumpt",
    version = "v0.5.0",
    visibility = ["PUBLIC"],
    deps = [
        "///third_party/go/github.com_google_go-cmp//cmp",
        "///third_party/go/golang.org_x_mod//module",
        "///third_party/go/golang.org_x_sync//semaphore",
        "///third_party/go/golang.org_x_tools//go/ast/astutil",
    ],
)
