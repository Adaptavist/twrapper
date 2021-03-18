SRC=./cmd/twapper
DIST=./dist

test : mods
	go test -v ./...

test_quiet : mods
	go test ./...

dist: clean test
	GOOS=darwin  GOARCH=amd64 go build -o ${DIST}/bpr-darwin-amd64  ${SRC}
	GOOS=darwin  GOARCH=arm64 go build -o ${DIST}/bpr-darwin-arm64  ${SRC}
	GOOS=linux   GOARCH=386   go build -o ${DIST}/bpr-linux-386     ${SRC}
	GOOS=linux   GOARCH=amd64 go build -o ${DIST}/bpr-linux-amd64   ${SRC}
	GOOS=linux   GOARCH=arm   go build -o ${DIST}/bpr-linux-arm     ${SRC}
	GOOS=linux   GOARCH=arm64 go build -o ${DIST}/bpr-linux-arm64   ${SRC}
	GOOS=windows GOARCH=386   go build -o ${DIST}/bpr-windows-386   ${SRC}
	GOOS=windows GOARCH=amd64 go build -o ${DIST}/bpr-windows-amd64 ${SRC}
	GOOS=windows GOARCH=arm   go build -o ${DIST}/bpr-windows-arm   ${SRC}

clean :
	rm -rf ${DIST}

run : test
	go run ${SRC}

mods:
	go mod download

check : test
	staticcheck ./cmd/twapper/
	staticcheck ./pkg/aws/
	staticcheck ./pkg/terraform/
