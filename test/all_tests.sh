#!/bin/bash

TESTS_PATH="${GOPATH}/src/github.com/Miguel-Dorta/surveillance-cameras/test"

"${TESTS_PATH}/generic_listLargeDirs/test.sh"
"${TESTS_PATH}/generic_rmOldCameraData/test.sh"
