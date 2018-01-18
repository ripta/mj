MJVERSION=0.1.0
MJBUILDTS=`date -u +'%FT%T%z'`

FILES=$(shell find . -type f -name '*.go')

all: mj

build-docker:
	docker build --build-arg MJ_VERSION=$(MJVERSION) --build-arg MJ_BUILD_DATE=$(MJBUILDTS) .

clean:
	rm -f mj

mj: $(FILES)
	go build -ldflags "-X main.BuildVersion=${MJVERSION} -X main.BuildDate=${MJBUILDTS}"

test: $(FILES)
	go test
