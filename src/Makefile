PROG_NAME := "filelist-mirror-downloader"
VERSION = 0.1.$(shell date +%Y%m%d.%H%M)
FLAGS := "-s -w -X main.version=${VERSION}"


build:
	GOOS=linux CGO_ENABLED=0 go build -ldflags=${FLAGS} -o ../${PROG_NAME} .
