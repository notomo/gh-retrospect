GH_NAME:=retrospect

build:
	go build -o gh-${GH_NAME} main.go

install: build
	gh extension remove ${GH_NAME} || echo
	gh extension install .

start: install
	gh ${GH_NAME}

test:
	TZ=UTC go test -v ./...

log:
	gh ${GH_NAME} -log=/dev/stdout
