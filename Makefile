BINARY_NAME=ratelimiter

.PHONY: build
build:
 GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin main.go

.PHONY: run
run:
 ./${BINARY_NAME}-darwin

.PHONY: build_and_run
build_and_run: build run

.PHONY: clean
clean:
 go clean
 rm ${BINARY_NAME}-darwin

.PHONY: test
test:
	go test .

.PHONY: test_coverage
test_coverage:
 go test . -coverprofile=coverage.out

.PHONY: install
install:
 go mod download

.PHONY: vet
vet:
 go vet
