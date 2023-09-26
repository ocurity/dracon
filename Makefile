.PHONY: build clean clean-protos lint fmt fmt-go fmt-proto

proto_defs=$(shell find . -name "*.proto" -not -path "./vendor/*")
go_protos=$(proto_defs:.proto=.pb.go)

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
