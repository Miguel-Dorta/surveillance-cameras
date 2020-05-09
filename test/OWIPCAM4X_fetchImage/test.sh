#!/bin/bash

PROGRAM_NAME="OWIPCAM4X_fetchImage"
TMP_DIR="/tmp/${PROGRAM_NAME}_test"
DATA_DIR="${TMP_DIR}/data"
BUILD_PATH="${TMP_DIR}/build/${PROGRAM_NAME}"
CHECK_PATH="${TMP_DIR}/build/check"
SERVER_PATH="${TMP_DIR}/build/server"
PID_PATH="${TMP_DIR}"
URL_PATH="${TMP_DIR}/url.txt"

# Exit in case of error
set -e

function test_OWIPCAM4X_fetchImage() {
  # Clean previous test
  rm -Rf "$TMP_DIR"

  # Build executables
  go build -o "$BUILD_PATH" "${GOPATH}/src/github.com/Miguel-Dorta/surveillance-cameras/cmd/${PROGRAM_NAME}"
  go build -o "$CHECK_PATH" "${GOPATH}/src/github.com/Miguel-Dorta/surveillance-cameras/test/${PROGRAM_NAME}/check.go"
  go build -o "$SERVER_PATH" "${GOPATH}/src/github.com/Miguel-Dorta/surveillance-cameras/test/${PROGRAM_NAME}/server.go"

  # Create testdata
  $SERVER_PATH > "$URL_PATH" &
  SERVER_PID="$!"
  sleep 1

  # Execute it
  $BUILD_PATH -camera-name "CamName" -user "USER" -password "PASS" -url "$(cat $URL_PATH)/auto.jpg" -path "$DATA_DIR" -pid-directory "$PID_PATH" &
  sleep 10
  kill -SIGINT "$(cat ${PID_PATH}/OWIPCAM4X_fetchImage_CamName.pid)"
  kill -SIGINT "$SERVER_PID"

  # Execute test for checking it
  $CHECK_PATH "$DATA_DIR"
}

test_OWIPCAM4X_fetchImage || (echo ":: FAIL - ${PROGRAM_NAME} test" && exit 1)

echo ":: PASS - ${PROGRAM_NAME} test"

# Clean up
rm -Rf "$TMP_DIR"
