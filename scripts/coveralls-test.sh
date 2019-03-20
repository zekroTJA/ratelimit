#!/bin/bash

go get golang.org/x/tools/cmd/cover
go get github.com/mattn/goveralls

go test -v -covermode=count -coverprofile=coverage.out

$HOME/gopath/bin/goveralls -coverprofile=coverage.out \
    -service=travis-ci -repotoken $COVERALLS_TOKEN