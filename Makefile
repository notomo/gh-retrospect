GH_NAME:=retrospect

build:
	go build -o gh-${GH_NAME} main.go

install: build
	gh extension remove ${GH_NAME}
	gh extension install .

start: install
	gh ${GH_NAME}
