os 			:=	$(shell uname -s)
user		:=	$(shell whoami)
rev			:= 	$(shell git rev-parse --short HEAD)
ref			:=	$(shell git rev-parse --abbrev-ref HEAD)
revDate		:=	$(shell git log -1 --date=format:"%Y-%m-%d %T" --format="%ad")
revInfo		:=	$(shell echo "${ref}(${rev}: ${revDate})")

buildDate	:=	$(shell date +"%Y-%m-%d %T")
buildInfo	:=	$(shell echo ${buildDate})

basePath 	:=	$(shell pwd)
project		:= 	$(shell sed -n 1p go.mod | cut -d ' ' -f2)
output 		:=	$(shell echo "${project}_"`echo ${os} | tr '[A-Z]' '[a-z]'`)
outputDir	:=	"./bin"

ldflags		:=	$(shell echo "-s -w -X 'main.version=${revInfo}' -X 'main.buildInfo=${buildInfo}'")

# docker compose yaml dir
dockerComposeDir :=	""
# remote host
remoteHost	:= 	""
# remote work dir
remoteDir	:=	$(shell echo "${project}")

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
	@echo "use: build, build_office, office, office_build, office_login, sync_data, download_xlsx, test, test_coverage, tidy"

build:
	@echo "output: ${output}"
	go build -ldflags "${ldflags} -X main.env=develop" -o "${outputDir}/${output}"

build_office:
	@echo "output: ${output}"
	go build -ldflags "${ldflags} -X main.env=office" -o "${outputDir}/${output}"

office:
	@echo "office"
	@make office_build
	@echo "\nStart Run ..."
	@ssh -t ${remoteHost} "cd ${remoteDir}; ./${project}_linux; bash"

office_build:
	@echo "office build"
	@cd ${dockerComposeDir};docker-compose run --rm -w ${basePath} go make build_office
	@scp "${outputDir}/${project}_linux" "${remoteHost}:${remoteDir}"

office_login:
	@echo "office login"
	@ssh -t ${remoteHost} "cd ${remoteDir}; bash"

sync_data:
	@echo "sync data"
	@scp ./data/* ${remoteHost}:"${remoteDir}/data"
	@ssh -t ${remoteHost} "cd ${remoteDir}/data; ls -la"

download_xlsx:
	@echo "download xlsx"
	@scp -r "${remoteHost}:${remoteDir}/xlsx/*" "./xlsx/office"

test:
	echo "go test $(sed -n 1p go.mod | cut -d ' ' -f2)"
	go test -race ./...

test_coverage:
	echo "go test $(sed -n 1p go.mod | cut -d ' ' -f2)"
	go test -race -coverprofile=cover.pprof -covermode=atomic ./...

tidy:
	echo "go mod tidy $(sed -n 1p go.mod | cut -d ' ' -f2)"
	go mod tidy
