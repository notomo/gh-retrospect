GH_NAME:=retrospect
GH_FULL_NAME:=gh-${GH_NAME}

build:
	go build -o ${GH_FULL_NAME} main.go

install: build
	gh extension remove ${GH_NAME} || echo
	gh extension install .

start: install
	gh ${GH_NAME}

test:
	TZ=UTC go test -v ./...

log:
	gh ${GH_NAME} -log=/tmp/${GH_FULL_NAME}
