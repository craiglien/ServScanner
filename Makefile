include build-data-file.txt

OUTTESTS=tests/outside.txt
PROGTESTS=tests/test-prog.txt

TESTFILES=${OUTTESTS} ${PROGTESTS}

WINBIN=ServScanner-windows.exe
LNXBIN=ServScanner-linux

GOARGS=-ldflags "-w -s -X main.Version=${VERSION} -X main.buildTime=${BUILD}"

run: ServScanner.go data/input.go build-data-file.txt
	@echo -----
	go run ${GOARGS} ServScanner.go 

.PHONY: force-build
force-build: clean build

clean:
	rm -f ServScanner.exe
	rm -f ServScanner
	@true	

publish-to-web: build
	cp $(LNXBIN) /var/www/mapped-directory/

.PHONY: build for-linux for-windows
build:for-linux for-windows .git/COMMIT_EDITMSG

for-linux:$(LNXBIN)

for-windows:$(WINBIN)

$(LNXBIN):ServScanner.go data/input.go build-data-file.txt
	go build ${GOARGS} -o $@ ServScanner.go  

$(WINBIN):ServScanner.go data/input.go build-data-file.txt
	GOOS=windows GOARCH=386 go build ${GOARGS} -o $@ ServScanner.go  

buildTests: buildTests.go
	go build buildTests.go

data/input.go: buildTests ${TESTFILES} Makefile
	./buildTests ${TESTFILES}  > $@

.PHONY: test
test:
	go test -v

build-data-file.txt: Makefile .git/COMMIT_EDITMSG
	echo -n BUILD= > $@
	date +%FT%T%z >> $@
	echo -n GITHASH= >> $@
	git tag >> $@
	echo -n VERSION= >> $@
	git describe --tags >> $@

.PHONY: show-opts
show-opts:
	echo ${BUILD}
	sleep 1
	echo ${BUILD}
	echo ${GITHASH}
	echo ${VERSION}
