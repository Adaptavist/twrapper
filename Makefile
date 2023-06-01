SRC=./cmd/twrapper
DIST=./dist

dist: clean check
	GOOS=darwin  GOARCH=amd64 go build -o ${DIST}/twrapper-darwin-amd64  ${SRC}
	GOOS=darwin  GOARCH=arm64 go build -o ${DIST}/twrapper-darwin-arm64  ${SRC}
	GOOS=linux   GOARCH=386   go build -o ${DIST}/twrapper-linux-386     ${SRC}
	GOOS=linux   GOARCH=amd64 go build -o ${DIST}/twrapper-linux-amd64   ${SRC}
	GOOS=linux   GOARCH=arm   go build -o ${DIST}/twrapper-linux-arm     ${SRC}
	GOOS=linux   GOARCH=arm64 go build -o ${DIST}/twrapper-linux-arm64   ${SRC}
	GOOS=windows GOARCH=386   go build -o ${DIST}/twrapper-windows-386   ${SRC}
	GOOS=windows GOARCH=amd64 go build -o ${DIST}/twrapper-windows-amd64 ${SRC}
	GOOS=windows GOARCH=arm   go build -o ${DIST}/twrapper-windows-arm   ${SRC}

clean :
	rm -rf ${DIST}

run : test
	go run ${SRC}

mods:
	go mod download

check : test
	staticcheck ./...

test : mods
	go test -v ./...

test_quiet : mods
	go test ./...

integration_tests:
	go test --tags=integration_test ./...

vet:
	go vet ./...