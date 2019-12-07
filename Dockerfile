FROM golang:latest

RUN # mkdir go/src

WORKDIR /go/src

ADD . /go/src

ENV export GO111MODULE=on