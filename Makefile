BINARY=ai-agent-demo
VERSION=1.0.0
DATE=`date +%FT%T%z`
.PHONY: build deploy

default:
	@echo ${BINARY}
	@echo ${VERSION}
	@echo ${DATE}

build:
	@go build -o build/${BINARY}
	@echo "[ok] build ${BINARY}"

deploy:
