#!/bin/bash

PROGRAM_NAME="generic_rmOldCameraData"
TMP_DIR="/tmp/${PROGRAM_NAME}_test"
DATA_DIR="${TMP_DIR}/data"
BUILD_PATH="${TMP_DIR}/build/${PROGRAM_NAME}"
CHECK_PATH="${TMP_DIR}/build/check"
GEN_PATH="${TMP_DIR}/build/generate"

# Exit in case of error
set -e

function test_generic_rmOldCameraData() {
  # Clean previous test
  rm -Rf "$TMP_DIR"

  # Build executables
  go build -o "$BUILD_PATH" "${GOPATH}/src/github.com/Miguel-Dorta/surveillance-cameras/cmd/${PROGRAM_NAME}"
  go build -o "$CHECK_PATH" "${GOPATH}/src/github.com/Miguel-Dorta/surveillance-cameras/test/${PROGRAM_NAME}/check.go"
  go build -o "$GEN_PATH" "${GOPATH}/src/github.com/Miguel-Dorta/surveillance-cameras/test/${PROGRAM_NAME}/generateData.go"

  # Create testdata
  $GEN_PATH "$DATA_DIR"

  # Execute it
  $BUILD_PATH -path "$DATA_DIR" -days 0

  # Execute test for checking it
  $CHECK_PATH "$DATA_DIR"
}

test_generic_rmOldCameraData || (echo ":: FAIL - ${PROGRAM_NAME} test" && exit 1)

echo ":: PASS - ${PROGRAM_NAME} test"

# Clean up
rm -Rf "$TMP_DIR"
