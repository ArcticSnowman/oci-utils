.PHONY: generate
generate:
	@echo "Running $@"
	@go generate ./...


.PHONY: compile
compile: generate
	@echo "Running $@"
	@CGO_ENABLED=0 GOBIN=${PWD}/bin go install -v ./oci-list-volumes
	