#!/usr/bin/env bash

if [[ $( uname -o ) == 'Cygwin' ]] ;then
	GOPATH=$( cygpath -w $( pwd ) )
	GOBIN=$GOPATH\\bin
else
	GOPATH=$( pwd )
	GOBIN=$GOPATH/bin
fi

export GOPATH GOBIN
set | grep GO

rm -Rvf bin pkg

go get -v github.com/stretchr/testify/assert \
    && go test -v orders \
    && go test -v pgdumpsplit/args \
    && go install -v pgdumpsplit #\
#    && go test -v mydumpsplit \
#    && go install -v mydumpsplit/mydumpsplit

