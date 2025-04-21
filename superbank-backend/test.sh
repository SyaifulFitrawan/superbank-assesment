#!/bin/bash

mkdir -p test

go test ./module/... -coverprofile test/coverage-full.out -covermode atomic -coverpkg ./...

head -n 1 test/coverage-full.out > test/coverage.out

grep 'module/' test/coverage-full.out | \
  grep -vE '\.dto\.go|\.interface\.go|\.container\.go' >> test/coverage.out

COVERAGE=$(go tool cover -func=test/coverage.out | grep total: | awk '{print substr($3, 1, length($3)-1)}')
THRESHOLD=95.0

go tool cover -func=test/coverage.out
go tool cover -html=test/coverage.out -o test/index.html

if (( $(echo "$COVERAGE < $THRESHOLD" | bc -l) )); then
    echo -e "\033[0;31mCoverage below threshold: $COVERAGE% < $THRESHOLD%\033[0m"
    exit 1
else
    echo -e "\033[0;32mCoverage passed: $COVERAGE%\033[0m"
fi