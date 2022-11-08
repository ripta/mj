ARG MJ_BUILD_DATE
ARG MJ_VERSION

FROM golang:1.19-bullseye AS build
ENV MJ_BUILD_DATE=$MJ_BUILD_DATE MJ_VERSION=$MJ_VERSION
WORKDIR $GOPATH/src/github.com/ripta/mj
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go test -v ./...
RUN go install -ldflags "-s -w -X main.BuildCommit=$(git rev-parse HEAD) -X main.BuildVersion=$MJ_VERSION -X main.BuildDate=$MJ_BUILD_DATE" ./...
RUN mv $GOPATH/bin/mj /bin/mj

FROM debian:bullseye
COPY --from=build /bin/mj /mj
ENTRYPOINT ["/mj"]
