#!/bin/bash
GOVERSION="1.16.3"
REPOSITORY="your_repo_name"
VERSION="1.0.0"

cat > dockerfile <<EOS
FROM golang:$GOVERSION-buster
RUN apt-get -y update && \ 
apt-get install python3 -y && \
apt-get install openjdk-11-jdk-headless -y && \
go get github.com/bazelbuild/bazelisk && \
export PATH=\$PATH:\$(go env GOPATH)/bin
COPY go.mod go.sum bazelwebapi/*.go /app/web-cd-service/
WORKDIR /app/web-cd-service
RUN go build && rm -r /go/pkg/*
ENTRYPOINT ["/app/web-cd-service/web-cd-service"]
EXPOSE 8088
EOS

docker build -t ${REPOSITORY}/web-cd-service:${VERSION} .
docker login
docker push ${REPOSITORY}/web-cd-service:${VERSION}
docker logout