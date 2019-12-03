#!/bin/bash

PROGRAM_NAME="generic_listLargeDirs"
TMP_DIR="/tmp/${PROGRAM_NAME}_test"
DATA_DIR="${TMP_DIR}/data"
BUILD_PATH="${TMP_DIR}/build/${PROGRAM_NAME}"
CHECK_PATH="${TMP_DIR}/build/check"
GEN_PATH="${TMP_DIR}/build/generate"
FILES_PATH="${TMP_DIR}/files.txt"
STDOUT_PATH="${TMP_DIR}/stdout"

# Exit in case of error
set -e

function test_generic_listLargeDirs() {
  # Clean previous test
  rm -Rf "$TMP_DIR"

  # Build executables
  go build -o "$BUILD_PATH" "${GOPATH}/src/github.com/Miguel-Dorta/surveillance-cameras/cmd/${PROGRAM_NAME}"
  go build -o "$CHECK_PATH" "${GOPATH}/src/github.com/Miguel-Dorta/surveillance-cameras/test/${PROGRAM_NAME}/check.go"
  go build -o "$GEN_PATH" "${GOPATH}/src/github.com/Miguel-Dorta/surveillance-cameras/test/${PROGRAM_NAME}/generateData.go"

  # Create testdata
  $GEN_PATH "$DATA_DIR" "$FILES_PATH"

  # Execute it
  $BUILD_PATH "$DATA_DIR" > $STDOUT_PATH

  # Execute test for checking it
  $CHECK_PATH "$STDOUT_PATH" "$FILES_PATH"
}

test_generic_listLargeDirs || (echo ":: FAIL - generic_listLargeDirs test" && exit 1)

echo ":: PASS - generic_listLargeDirs test"

# Clean up
rm -Rf "$TMP_DIR"
