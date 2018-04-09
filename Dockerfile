FROM golang:1.10-alpine3.7 AS build
RUN apk add --update --no-cache git
ENV GOPROJ=github.com/ripta/mj

ARG MJ_BUILD_DATE
ARG MJ_VERSION
ENV MJ_BUILD_DATE=$MJ_BUILD_DATE MJ_VERSION=$MJ_VERSION

WORKDIR $GOPATH/src/$GOPROJ
COPY . .
RUN go get -d ./...
RUN go install -ldflags "-s -w -X main.BuildCommit=$(git rev-parse HEAD) -X main.BuildVersion=$MJ_VERSION -X main.BuildDate=$MJ_BUILD_DATE" ./...
RUN mv $GOPATH/bin/mj /bin/mj

FROM scratch
COPY --from=build /bin/mj /mj
ENTRYPOINT ["/mj"]
