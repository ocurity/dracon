.PHONY: build clean clean-protos lint fmt fmt-go fmt-proto

proto_defs=$(shell find . -name "*.proto" -not -path "./vendor/*")
go_protos=$(proto_defs:.proto=.pb.go)

TEKTON_VERSION=0.44.0

PROTOC=protoc

build: $(go_protos)
	@echo bla

$(go_protos): %.pb.go: %.proto
	$(PROTOC) --go_out=. --go_opt=paths=source_relative $<

clean: clean-protos

clean-protos:
	rm -rf $(go_protos)

fmt-proto:
	@echo "Tidying up Proto files"
	@buf format -w "./api/proto"

fmt-go:
	@echo "Tidying up Go files"
	@gofmt -l -w $$(find . -name *.go -not -path "./vendor/*" | xargs -n 1 dirname | uniq)

fmt: fmt-go fmt-proto

lint:
	@if [ "${CI}" = "true" ]; then\
		reviewdog -fail-on-error -reporter=github-pr-review;\
	else\
		reviewdog -fail-on-error -diff="git diff origin/main";\
	fi
	@golangci-lint

third_party/k8s/tektoncd/pipeline/swagger-v$(TEKTON_VERSION).json:
	wget "https://raw.githubusercontent.com/tektoncd/pipeline/v$(TEKTON_VERSION)/pkg/apis/pipeline/v1beta1/swagger.json" -O third_party/k8s/tektoncd/pipeline/swagger-v$(TEKTON_VERSION).json

third_party/k8s/tektoncd/pipeline/release-v$(TEKTON_VERSION).yaml:
	wget "https://storage.googleapis.com/tekton-releases/pipeline/previous/v$(TEKTON_VERSION)/release.yaml" -O third_party/k8s/tektoncd/pipeline/release-v$(TEKTON_VERSION).yaml

api/openapi/tekton/openapi_schema.json: third_party/k8s/tektoncd/pipeline/swagger-v$(TEKTON_VERSION).json
	./scripts/generate_openapi_schema.sh $< $@

mirror_images:
	./scripts/mirror_images.sh

third_party/k8s/tektoncd/pipeline/pipeline.yaml: third_party/k8s/tektoncd/pipeline/release-v$(TEKTON_VERSION).yaml
	cp third_party/k8s/tektoncd/pipeline/release-v$(TEKTON_VERSION).yaml third_party/k8s/tektoncd/pipeline/release.yaml
	kustomize build $(shell dirname $@) > $@
	rm third_party/k8s/tektoncd/pipeline/release.yaml
