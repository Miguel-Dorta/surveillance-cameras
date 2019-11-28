#!/bin/bash

PROGRAM_NAME="generic_rmOldCameraData"
TMP_DIR="/tmp/${PROGRAM_NAME}"
DATA_DIR="${TMP_DIR}/data"
BUILD_PATH="${TMP_DIR}/build/${PROGRAM_NAME}"
CHECK_PATH="${TMP_DIR}/build/check"

# Clean previous test
rm -Rf "$TMP_DIR"

# Build executables
go build -o "$BUILD_PATH" "${GOPATH}/src/github.com/Miguel-Dorta/surveillance-cameras/cmd/${PROGRAM_NAME}"
go build -o "$CHECK_PATH" "${GOPATH}/src/github.com/Miguel-Dorta/surveillance-cameras/test/${PROGRAM_NAME}_check.go"

# Create testdata
for (( c = 0; c <= 5; c++ )); do
    cPath="${DATA_DIR}/C$(printf %03d $c)"

    for (( y = 1950; y <= 2100; y++ )); do
        yPath="${cPath}/${y}"

        for (( m = 1; m <= 12; m++ )); do
            mPath="${yPath}/$(printf %02d $m)"

            for (( d = 1; d <= 31; d++ )); do
                dPath="${mPath}/$(printf %02d $d)"
                mkdir -p "$dPath"
                touch "${dPath}/testfile"
            done
        done
    done
done

# Execute it
$BUILD_PATH -path "$DATA_DIR" -days 0

# Execute test for checking it
$CHECK_PATH "$DATA_DIR"
