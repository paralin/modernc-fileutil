# Copyright (c) 2014 The fileutil authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

.PHONY: all clean edit editor build_all_targets todo

all: editor

build_all_targets:
	GOOS=darwin GOARCH=amd64 go build -v ./...
	GOOS=darwin GOARCH=arm64 go build -v ./...
	GOOS=freebsd GOARCH=386 go build -v ./...
	GOOS=freebsd GOARCH=amd64 go build -v ./...
	GOOS=linux GOARCH=386 go build -v ./...
	GOOS=linux GOARCH=amd64 go build -v ./...
	GOOS=linux GOARCH=arm go build -v ./...
	GOOS=linux GOARCH=arm64 go build -v ./...
	GOOS=linux GOARCH=ppc64le go build -v ./...
	GOOS=linux GOARCH=riscv64 go build -v ./...
	GOOS=linux GOARCH=s390x go build -v ./...
	GOOS=netbsd GOARCH=amd64 go build -v ./...
	GOOS=openbsd GOARCH=amd64 go build -v ./...
	GOOS=openbsd GOARCH=amd64 go build -v ./...
	GOOS=windows GOARCH=386 go build -v ./...
	GOOS=windows GOARCH=amd64 go build -v ./...

clean:
	rm -f log-* cpu.test mem.test *.out
	go clean

edit:
	@touch log
	@if [ -f "Session.vim" ]; then gvim -S & else gvim -p Makefile *.go & fi

editor:
	gofmt -l -s -w .
	go test
	go build
