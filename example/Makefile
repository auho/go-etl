user	:=	$(shell whoami)
rev		:= 	$(shell git rev-parse --short HEAD)
os 		:=	$(shell uname -s)
project	:= 	$(shell sed -n 1p go.mod | cut -d ' ' -f2)
output 	:=	$(shell echo "${project}_"`echo ${os} | tr '[A-Z]' '[a-z]'`)

# Mac OS X
ifeq ($(shell uname),Darwin)
GOBIN		:=	$(shell echo ${GOBIN} | cut -d':' -f1)
GOPATH		:=	$(shell echo ${GOPATH} | cut -d':' -f1)
endif

# Linux
ifeq ($(os),Linux)
GOBIN		:=	$(shell echo ${GOBIN} | cut -d':' -f1)
GOPATH		:=	$(shell echo ${GOPATH} | cut -d':' -f1)
endif

# Windows
ifeq ($(os),MINGW)
GOBIN	:=	$(subst \,/,$(GOBIN))
GOPATH	:=	$(subst \,/,$(GOPATH))
GOBIN 	:=	/$(shell echo "$(GOBIN)" | cut -d';' -f1 | sed 's/://g')
GOPATH 	:=	/$(shell echo "$(GOPATH)" | cut -d';' -f1 | sed 's/://g')
endif

help:
	@echo "go bin: ${GOBIN}"
	@echo "go path: ${GOPATH}"
	@echo "os: ${os}"
	@echo "cgo enabled: ${CGO_ENABLED}"
	@echo ""
	@echo "use: build, build_office, test, test_coverage, tidy"

build:
	@echo "output: ${output}"
	go build -ldflags "-s -w -X main.env=develop -X main.version=${rev}" -o "${output}"

build_office:
	@echo "output: ${output}"
	go build -ldflags "-s -w -X main.env=office -X main.version=${rev}" -o "${output}"

test:
	echo "go test $(sed -n 1p go.mod | cut -d ' ' -f2)"
	go test -race ./...

test_coverage:
	echo "go test $(sed -n 1p go.mod | cut -d ' ' -f2)"
	go test -race -coverprofile=cover.pprof -covermode=atomic ./...

tidy:
	echo "go mod tidy $(sed -n 1p go.mod | cut -d ' ' -f2)"
	go mod tidy
