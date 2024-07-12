help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

functions := $(shell cd .. && find functions -name \*main.go | awk -F'/' '{print $$2}')

mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(patsubst %/,%,$(dir $(mkfile_path)))

build: ## Build golang binaries
	rm -rf ${current_dir}/bin
	mkdir -p ${current_dir}/bin
	echo "Build go artifacts in ${current_dir}..."
	@for function in $(functions) ; do \
		echo "Build $$function" ; \
		cd ${current_dir}/../functions/$$function ; \
		env GOARCH=amd64 GOOS=linux go build -o ${current_dir}/bin/$$function/bootstrap *.go ; \
		cd ${current_dir}/bin/$$function; \
		zip ../$$function.zip bootstrap; \
	done
