FROM golang:alpine as builder

RUN apk --no-cache add git gcc g++

WORKDIR /go/src/github.com/seawolflin/gitlab-exporter/
COPY . .

RUN apk --no-cache add git gcc \
    && GOPROXY=https://goproxy.cn go mod download -x \
    && GOPATH=/go go build -o gitlab-exporter github.com/seawolflin/gitlab-exporter/cmd/gitlab-exporter


FROM golang:1.18-alpine as prod


WORKDIR /root/

COPY --from=0 /go/src/github.com/seawolflin/gitlab-exporter/gitlab-exporter .

ENV URL ""
ENV TOKEN ""

CMD /root/gitlab-exporter --url $URL --token $TOKEN
