FROM golang:alpine as build

ADD . /go/src/github.com/terranodo/tegola
RUN cd /go/src/github.com/terranodo/tegola/cmd/tegola; go build -o tegola


