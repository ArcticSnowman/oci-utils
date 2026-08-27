.PHONY: generate
generate:
	@echo "Running $@"
	@go generate ./...


.PHONY: compile
compile: generate
	@echo "Running $@"
	@CGO_ENABLED=0 GOBIN=${PWD}/bin go install -v ./oci-compute-capacity-report
	
.PHONY: compress
compress: compile
	@echo "Running $@"
	@upx bin/oci-compute-capacity-report 