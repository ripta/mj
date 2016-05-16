MJVERSION=0.1.0
MJBUILDTS=`date -u +'%FT%T%z'`

FILES=$(shell find . -type f -name '*.go')

all: mj

clean:
	rm -f mj

mj: $(FILES)
	go build -ldflags "-X main.version=${MJVERSION} -X main.builtAt=${MJBUILDTS}"

test: $(FILES)
	go test
