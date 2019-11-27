#!/bin/bash

TMP_DIR="/tmp/generic_rmOldCameraData"
DATA_DIR="${TMP_DIR}/data"
BUILD_PATH="${TMP_DIR}/build/generic_rmOldCameraData"

# Clean previous test
rm -Rf "$TMP_DIR"

# Build executable
go build -o "$BUILD_PATH" "${GOPATH}/src/github.com/Miguel-Dorta/surveillance-cameras/cmd/generic_rmOldCameraData"

# Create testdata
for (( c = 0; c <= 5; c++ )); do
    cPath="${DATA_DIR}/C$(printf %03d $c)"

    for (( y = 1950; y <= 2100; y++ )); do
        yPath="${cPath}/${y}"

        for (( m = 1; m <= 12; m++ )); do
            mPath="${yPath}/$(printf %02d $m)"

            for (( d = 1; d < 31; d++ )); do
                dPath="${mPath}/$(printf %02d $d)"
                mkdir -p "$dPath"
                touch "${dPath}/testfile"
            done
        done
    done
done

# Execute it
$BUILD_PATH -path "$DATA_DIR" -days 0

# Output
echo "Please, check $DATA_DIR"
