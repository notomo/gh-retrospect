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

LOG_DIR:=/tmp/${GH_FULL_NAME}
log:
	rm -rf ${LOG_DIR}
	gh ${GH_NAME} -log ${LOG_DIR}
