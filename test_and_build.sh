#!/usr/bin/env bash

if [[ $( uname -o ) == 'Cygwin' ]] ;then
	GOPATH=$( cygpath -w $( pwd ) )
	GOBIN=$GOPATH\\bin
else
	GOPATH=$( pwd )
	GOBIN=$GOPATH/bin
fi

export GOROOT GOBIN
set | grep GO

rm -Rvf bin pkg

go get -v github.com/stretchr/testify/assert && \
	go test -v dbtricks && \
	# go test -v pgsplit && \
	go install -v pgsplit
