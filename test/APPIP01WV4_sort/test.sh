#!/bin/bash

PROGRAM_NAME="APPIP01WV4_sort"
TMP_DIR="/tmp/${PROGRAM_NAME}_test"
DATA_DIR="${TMP_DIR}/data"
SORT_DIR="${TMP_DIR}/sort"
BUILD_PATH="${TMP_DIR}/build/${PROGRAM_NAME}"
CHECK_PATH="${TMP_DIR}/build/check"
GEN_PATH="${TMP_DIR}/build/generate"
FILES_PATH="${TMP_DIR}/files.json"
PID_PATH="${TMP_DIR}/pid.pid"

# Exit in case of error
set -e

function test_APPIP01WV4_sort() {
  # Clean previous test
  rm -Rf "$TMP_DIR"

  # Build executables
  go build -o "$BUILD_PATH" "${GOPATH}/src/github.com/Miguel-Dorta/surveillance-cameras/cmd/${PROGRAM_NAME}"
  go build -o "$CHECK_PATH" "${GOPATH}/src/github.com/Miguel-Dorta/surveillance-cameras/test/${PROGRAM_NAME}/check.go"
  go build -o "$GEN_PATH" "${GOPATH}/src/github.com/Miguel-Dorta/surveillance-cameras/test/${PROGRAM_NAME}/generateData.go"

  # Create testdata
  $GEN_PATH "$DATA_DIR" "$FILES_PATH"

  # Execute it
  $BUILD_PATH -from "$DATA_DIR" -to "$SORT_DIR" -pid "$PID_PATH"

  # Execute test for checking it
  $CHECK_PATH "$SORT_DIR" "$FILES_PATH"
}

test_APPIP01WV4_sort || (echo ":: FAIL - ${PROGRAM_NAME} test" && exit 1)

echo ":: PASS - ${PROGRAM_NAME} test"

# Clean up
rm -Rf "$TMP_DIR"
