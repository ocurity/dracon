.PHONY: build clean

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
