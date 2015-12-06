GOPATH	= ${CURDIR}
PKGS	= $(subst _${GOPATH},.,$(shell go list ./...))

all: main test coverage

main:
	go build main.go

clean: cleancoverage
	rm main

test:
	go test ./...

vtest:
	go test -v ./...

fmt:
	go fmt ./...

coverage:
	echo /dev/null > coverage.log
	for pkg in ${PKGS} ; do \
		go test -cover -coverprofile=$$pkg/coverage.out $$pkg; \
		if [ -e $$pkg/coverage.out ] ; then echo $$pkg/coverage.out >> coverage.log; fi \
	done
	grep -h '^mode:' `cat coverage.log` | sort | uniq	 > all-coverage.out
	grep -h -v '^mode:' `cat coverage.log`			>> all-coverage.out
	go tool cover -html=all-coverage.out

cleancoverage:
	rm -f coverage.log
	find . -name coverage.out | xargs rm -f
	rm -f all-coverage.out
