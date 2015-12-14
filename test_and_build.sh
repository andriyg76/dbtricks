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
    && go test -v dbtricks \
    && go test -v dbtricks/params \
    && go test -v dbtricks/orders \
    && go test -v dbtricks/writer \
    && go test -v mergesort \
    && go test -v pg/dumpsplit \
    && go test -v pg/datasplit \
    && go test -v pgdumpsplit \
    && go install -v pgdumpsplit #\
