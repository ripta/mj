FROM golang:1.9.2-alpine3.6 AS build
RUN apk add --update --no-cache git

ARG MJ_BUILD_DATE
ARG MJ_VERSION
ENV MJ_BUILD_DATE=$MJ_BUILD_DATE MJ_VERSION=$MJ_VERSION

RUN go-wrapper download github.com/ripta/mj
RUN go-wrapper install -ldflags "-s -w -X main.BuildVersion=$MJ_VERSION -X main.BuildDate=$MJ_BUILD_DATE" github.com/ripta/mj

FROM scratch
COPY --from=build /go/bin/mj /mj
ENTRYPOINT ["/mj"]
