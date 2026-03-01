#!/bin/bash

echo 'Running acceptance tests...'

echo "Building Go binary..."
go build -o hw3-logs ./cmd/app

if [ ! -f "hw3-logs" ]; then
    echo "ERROR: Binary hw3-logs not found!"
    exit 1
fi

testNumber=0
failedTests=0

function assertExitCode {
  RED='\033[0;31m'
  GREEN='\033[0;32m'
  NC='\033[0m'

  echo "Exit code: expected=$1, actual=$2"
  if [ $1 -ne $2 ]; then
     echo -e "${RED}Test №${testNumber} failed${NC}"
     failedTests=$((failedTests+1))
  else
     echo -e "${GREEN}Test №${testNumber} passed${NC}"
  fi
}

function assertJsonEquals {
  echo "Comparing JSON's..."
  if command -v jq >/dev/null 2>&1; then
    diff <(jq --sort-keys . "$1") <(jq --sort-keys . "$2")
    exit_code=$?
  else
    diff "$1" "$2"
    exit_code=$?
  fi
  assertExitCode 0 $exit_code
}

function verifyAllTestsPassed {
  RED='\033[0;31m'
  GREEN='\033[0;32m'
  NC='\033[0m'

  echo "Total failed tests: ${failedTests}"
  if [ $failedTests -ne 0 ]; then
     echo -e "${RED}Some tests have failed!${NC}"
     exit 1
  else
     echo -e "${GREEN}All tests passed${NC}"
     exit 0
  fi
}

function runTest {
 echo "Test [№${testNumber}][$1]: $2; expected exit code: $3; args: ${@:4};"
 

 local args=("${@:4}")
 for i in "${!args[@]}"; do
    args[i]="${args[i]//\/tmp\/data/$(pwd)/scripts/data}"
    args[i]="${args[i]//\/tmp\/input/$(pwd)/scripts/data/input}"
    args[i]="${args[i]//\/tmp\/output/$(pwd)/scripts/data/output}"
 done
 

 ./hw3-logs "${args[@]}"
 
 exit_code=$?
 assertExitCode $3 $exit_code

 testNumber=$((testNumber+1))
}

mkdir -p scripts/data/output
mkdir -p scripts/data/input/logs

echo "test log data" > scripts/data/input/file2.txt
echo '{"existing": "data"}' > scripts/data/output/existing.json

echo "Running negative tests..."

runTest "negative" "input file does not exist" 2 \
  -p /tmp/data/input/nonexistent.txt -f json -o /tmp/data/output/output1.json

runTest "negative" "input file has unsupported extension" 2 \
  -p /tmp/data/input/file1.html -f json -o /tmp/data/output/output2.json

runTest "negative" "output file already exists" 2 \
  -p /tmp/data/input/file2.txt -f json -o /tmp/data/output/existing.json

runTest "negative" "output file has unsupported extension (JSON)" 2 \
  -p /tmp/data/input/file2.txt -f json -o /tmp/data/output/output4.txt

runTest "negative" "output file has unsupported extension (MD)" 2 \
  -p /tmp/data/input/file2.txt -f markdown -o /tmp/data/output/output5.txt

runTest "negative" "output file has unsupported extension (AD)" 2 \
  -p /tmp/data/input/file2.txt -f adoc -o /tmp/data/output/output6.txt

runTest "negative" "unsupported output format" 2 \
  -p /tmp/data/input/file2.txt -f txt -o /tmp/data/output/output7.txt

runTest "negative" "invalid date format (--from)" 2 \
  -p /tmp/data/input/file2.txt -f json -o /tmp/data/output/output8.json --from="2025.01.02"

runTest "negative" "invalid date format (--to)" 2 \
  -p /tmp/data/input/file2.txt -f json -o /tmp/data/output/output9.json --to="2025.01.02"

runTest "negative" "--from > --to" 2 \
  -p /tmp/data/input/file2.txt -f json -o /tmp/data/output/output10.json --from="2025-01-02" --to="2025-01-01"

runTest "negative" "required parameter -p is missing" 2 \
  -f json -o /tmp/data/output/output11.json

runTest "negative" "required parameter -f is missing" 2 \
  -p /tmp/data/input/nonexistent.txt -o /tmp/data/output/output12.json

runTest "negative" "required parameter -o is missing" 2 \
  -p /tmp/data/input/nonexistent.txt -f json

runTest "negative" "unsupported parameter is present" 2 \
  -p /tmp/data/input/nonexistent.txt -f json -o /tmp/data/output/output14.json --custom=argument

runTest "positive" "properly calculate statistics from multiple local files" 0 \
  -p "/tmp/data/input/logs/part*.txt" -f json -o /tmp/data/output/stats.json

assertJsonEquals ./scripts/data/output/expected.json ./scripts/data/output/stats.json

verifyAllTestsPassed